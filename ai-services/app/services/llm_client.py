"""OpenAI-compatible client pointing at the local LM Studio instance."""
import base64
import json
import re
from pathlib import Path

from openai import AsyncOpenAI

from app.config import settings

_client: AsyncOpenAI | None = None


def get_client() -> AsyncOpenAI:
    global _client
    if _client is None:
        _client = AsyncOpenAI(
            base_url=settings.llm_base_url,
            api_key=settings.llm_api_key,
        )
    return _client


def _encode_image(image_bytes: bytes, mime_type: str = "image/jpeg") -> str:
    return f"data:{mime_type};base64,{base64.b64encode(image_bytes).decode()}"


def ocr_invoice_pdf(pdf_bytes: bytes) -> str:
    """Extract text from a PDF using pypdf (no image conversion needed)."""
    try:
        import pypdf  # type: ignore
        reader = pypdf.PdfReader(__import__("io").BytesIO(pdf_bytes))
        pages: list[str] = []
        for page in reader.pages:
            text = page.extract_text() or ""
            if text.strip():
                pages.append(text)
        return "\n\n--- 下一页 ---\n\n".join(pages)
    except Exception as exc:
        raise RuntimeError(f"pypdf 文本提取失败: {exc}") from exc


async def ocr_invoice_image(image_bytes: bytes, mime_type: str = "image/jpeg") -> str:
    """Pass 1 — ask the multimodal model to extract raw text from the invoice image."""
    client = get_client()
    data_url = _encode_image(image_bytes, mime_type)
    resp = await client.chat.completions.create(
        model=settings.llm_model,
        messages=[
            {
                "role": "user",
                "content": [
                    {
                        "type": "image_url",
                        "image_url": {"url": data_url},
                    },
                    {
                        "type": "text",
                        "text": (
                            "你是一个 OCR 引擎。请将这张药品采购发票图片中的所有文字原样提取出来，"
                            "保持原始布局，不要遗漏任何数字、汉字、英文或标点符号。"
                            "直接输出提取到的文本，不要添加任何解释或格式化。"
                        ),
                    },
                ],
            }
        ],
        temperature=0.0,
        max_tokens=4096,
    )
    return resp.choices[0].message.content or ""


async def llm_correct_invoice(raw_ocr_text: str, supplier_hints: list[str] | None = None) -> dict:
    """Pass 2 — ask the LLM to parse and correct the raw OCR text into structured JSON.

    supplier_hints: list of known active supplier names to help the LLM match correctly.
    """
    client = get_client()

    supplier_context = ""
    if supplier_hints:
        names = "、".join(supplier_hints[:50])  # limit to 50 to avoid blowing context
        supplier_context = f"\n\n系统中已有的活跃供应商（仅供参考，用于纠正 OCR 错别字）：\n{names}\n"

    prompt = f"""你是一个药品采购发票解析专家。下面是从发票中提取的原始文本，可能存在 OCR 识别错误或排版混乱。
{supplier_context}
请仔细分析文本，提取并修正以下信息，以 JSON 格式返回，不要输出任何其他内容：

{{
  "invoice_no": "发票号码（字符串，无则 null）",
  "invoice_date": "开票日期，格式 YYYY-MM-DD（无则 null）",
  "supplier_name": "供应商/销售方名称（字符串，无则 null）",
  "total_amount": "价税合计金额，纯数字字符串保留两位小数（无则 null）",
  "items": [
    {{
      "row_index": 1,
      "drug_name": "药品名称",
      "specification": "规格",
      "manufacturer": "生产厂家",
      "approval_number": "批准文号",
      "batch_number": "批号",
      "expire_date": "有效期，格式 YYYY-MM-DD（无则 null）",
      "quantity": "数量，纯数字字符串（无则 null）",
      "unit_price": "单价，纯数字字符串保留两位小数（无则 null）",
      "amount": "金额，纯数字字符串保留两位小数（无则 null）"
    }}
  ]
}}

OCR 原始文本：
{raw_ocr_text}

JSON 结果："""

    resp = await client.chat.completions.create(
        model=settings.llm_model,
        messages=[{"role": "user", "content": prompt}],
        temperature=0.0,
        max_tokens=4096,
    )
    content = resp.choices[0].message.content or "{}"
    return _parse_json_from_llm(content)


async def llm_rank_drugs_for_query(user_query: str, candidates: list[dict]) -> dict:
    """Given a pre-fetched list of candidate drugs from inventory, ask the LLM which ones
    best match the user's query (symptoms, disease, approximate drug name, etc.).

    candidates: list of dicts with keys: drug_id, common_name, specification, dosage_form
    Returns: {
        "selected": [{"drug_id": int, "reason": str}, ...],  — ranked best→worst
        "explanation": str
    }
    Always returns at least one selection even if nothing is a strong match.
    """
    client = get_client()

    lines = []
    for i, c in enumerate(candidates, start=1):
        spec = c.get("specification") or ""
        form = c.get("dosage_form") or ""
        lines.append(f"{i}. ID={c['drug_id']} {c['common_name']} {spec} {form}".strip())
    drug_list = "\n".join(lines)

    max_select = min(5, len(candidates))

    prompt = (
        f"你是专业药剂师助手。药店工作人员输入了以下描述，请从库存药品中选出最合适的药品。\n\n"
        f"描述：{user_query}\n\n"
        f"库存药品列表（格式：序号. ID=编号 药品名称 规格 剂型）：\n{drug_list}\n\n"
        f"要求：\n"
        f"1. 只能从上面列表中选择，不能编造列表外的药品\n"
        f"2. drug_id 必须使用列表中的 ID= 后面的数字\n"
        f"3. 选1到{max_select}个，按适合程度从高到低排序\n"
        f"4. 即使没有完全匹配，也必须选出最有可能对症的药品\n\n"
        f"只输出以下 JSON，不要其他内容：\n"
        f'{{"selected":[{{"drug_id":数字,"reason":"推荐理由10字内"}}],"explanation":"1句话总结"}}'
    )

    resp = await client.chat.completions.create(
        model=settings.llm_model,
        messages=[
            {"role": "system", "content": "你是专业药剂师助手，直接输出 JSON，不要解释。"},
            {"role": "user", "content": prompt},
        ],
        temperature=0.1,
        max_tokens=4096,
    )
    raw = (resp.choices[0].message.content or "").strip()
    parsed = _parse_json_from_llm(raw)

    valid_ids = {c["drug_id"] for c in candidates}
    # Build name → drug_id map for fuzzy fallback
    name_to_id: dict[str, int] = {c["common_name"]: c["drug_id"] for c in candidates}

    selected = []
    for s in parsed.get("selected") or []:
        did = s.get("drug_id")
        reason = s.get("reason") or "AI推荐"

        if isinstance(did, (int, float)) and int(did) in valid_ids:
            # LLM returned a valid ID from our list
            selected.append({"drug_id": int(did), "reason": reason})
        else:
            # LLM may have invented an ID — try to match by the drug_name field if present
            llm_name = s.get("drug_name") or s.get("name") or ""
            if llm_name:
                # Exact match first
                matched_id = name_to_id.get(llm_name)
                if matched_id is None:
                    # Substring match: find any candidate whose name contains or is contained in llm_name
                    for cname, cid in name_to_id.items():
                        if llm_name in cname or cname in llm_name:
                            matched_id = cid
                            break
                if matched_id is not None:
                    selected.append({"drug_id": matched_id, "reason": reason})

    return {
        "selected": selected,
        "explanation": parsed.get("explanation") or "",
    }


def _parse_json_from_llm(text: str) -> dict:
    text = text.strip()
    # Strip markdown code fences if present
    text = re.sub(r"^```(?:json)?\s*", "", text)
    text = re.sub(r"\s*```$", "", text)
    text = text.strip()
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        # Try to extract first JSON object
        match = re.search(r"\{.*\}", text, re.DOTALL)
        if match:
            try:
                return json.loads(match.group())
            except json.JSONDecodeError:
                pass
    return {}
