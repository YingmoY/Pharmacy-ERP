"""Invoice recognition pipeline: OCR → LLM correction → DB match → persist."""
import hashlib
import io
import json
import time
from datetime import datetime, timezone, date as date_t

import asyncpg

from app.database import get_pool
from app.models.invoice import (
    InvoiceRecognizeItem,
    InvoiceRecognizeResult,
    InvoiceRecognizeResponse,
    QualityWarning,
    DrugCandidateInvoice,
)
from app.models.drug import DrugMatchInputItem
from app.models.supplier import SupplierCandidate
from app.request_id import new_request_id
from app.services import llm_client, matching_service


def _parse_date(val: str | None) -> date_t | None:
    if not val:
        return None
    for fmt in ("%Y-%m-%d", "%Y/%m/%d", "%Y年%m月%d日", "%Y.%m.%d"):
        try:
            return datetime.strptime(val, fmt).date()
        except ValueError:
            continue
    return None


def _parse_decimal_str(val) -> str | None:
    if val is None:
        return None
    s = str(val).strip()
    return s if s else None


async def _load_image_bytes_from_pdf(content: bytes) -> list[tuple[bytes, str]]:
    """Convert PDF pages to JPEG images. Returns list of (bytes, mime_type)."""
    try:
        import pdf2image  # type: ignore
        images = pdf2image.convert_from_bytes(content, dpi=200, fmt="jpeg")
        results = []
        for img in images:
            buf = io.BytesIO()
            img.save(buf, format="JPEG", quality=90)
            results.append((buf.getvalue(), "image/jpeg"))
        return results
    except Exception:
        # Fall back: treat raw bytes as single image
        return [(content, "application/pdf")]


def _images_from_content(content: bytes, content_type: str) -> list[tuple[bytes, str]]:
    ct = (content_type or "").lower()
    if "pdf" in ct:
        import asyncio
        loop = asyncio.get_event_loop()
        # pdf2image is sync; run in executor
        return []  # handled in async caller
    if ct in ("image/png",):
        return [(content, "image/png")]
    return [(content, "image/jpeg")]


async def recognize_invoice(
    file_content: bytes,
    file_name: str,
    content_type: str,
    erp_request_id: str | None,
    erp_file_id: str | None,
    match_master_data: bool,
    return_raw_response: bool,
) -> InvoiceRecognizeResponse:
    started_at = datetime.now(timezone.utc)
    request_id = new_request_id("INVAI")
    file_hash = hashlib.sha256(file_content).hexdigest()

    pool = get_pool()

    # Persist job as PROCESSING
    job_id: int = await pool.fetchval(
        """
        INSERT INTO ai.invoice_recognition_job
          (request_id, erp_request_id, erp_file_id, file_name,
           content_type, file_size, file_hash, status, started_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,'PROCESSING',$8)
        RETURNING id
        """,
        request_id, erp_request_id, erp_file_id, file_name,
        content_type, len(file_content), file_hash, started_at,
    )

    raw_response: dict = {}
    result: InvoiceRecognizeResult | None = None
    error_code: str | None = None
    error_message: str | None = None

    try:
        # ── Step 1: extract raw text ──────────────────────────────────────────
        ct = (content_type or "").lower()
        if "pdf" in ct:
            # For digital PDFs, extract text directly (fast, no image conversion)
            try:
                combined_ocr_text = llm_client.ocr_invoice_pdf(file_content)
                raw_response["ocr_method"] = "pypdf"
            except Exception:
                # Fall back to image-based OCR if pypdf fails (scanned PDF)
                image_pages = await _load_image_bytes_from_pdf(file_content)
                ocr_texts: list[str] = []
                for img_bytes, mime in image_pages:
                    text = await llm_client.ocr_invoice_image(img_bytes, mime)
                    ocr_texts.append(text)
                combined_ocr_text = "\n\n--- 下一页 ---\n\n".join(ocr_texts)
                raw_response["ocr_method"] = "llm_vision_fallback"
        elif ct == "image/png":
            combined_ocr_text = await llm_client.ocr_invoice_image(file_content, "image/png")
            raw_response["ocr_method"] = "llm_vision"
        else:
            combined_ocr_text = await llm_client.ocr_invoice_image(file_content, "image/jpeg")
            raw_response["ocr_method"] = "llm_vision"

        raw_response["ocr_text"] = combined_ocr_text

        # ── Step 2: load supplier context from DB ─────────────────────────────
        supplier_hints: list[str] = []
        try:
            rows = await pool.fetch(
                "SELECT name FROM ai.v_active_supplier_source ORDER BY name LIMIT 100"
            )
            supplier_hints = [r["name"] for r in rows]
        except Exception:
            pass  # context is optional; proceed without it

        # ── Step 3: LLM correction / structured extraction ────────────────────
        structured = await llm_client.llm_correct_invoice(
            combined_ocr_text, supplier_hints=supplier_hints or None
        )
        raw_response["llm_structured"] = structured

        # ── Step 4: parse structured output ───────────────────────────────────
        invoice_date = _parse_date(structured.get("invoice_date"))
        supplier_name: str | None = structured.get("supplier_name")
        invoice_no: str | None = structured.get("invoice_no")
        total_amount: str | None = _parse_decimal_str(structured.get("total_amount"))

        raw_items: list[dict] = structured.get("items") or []
        recognized_items: list[InvoiceRecognizeItem] = []

        for idx, raw in enumerate(raw_items, start=1):
            row_index = int(raw.get("row_index") or idx)
            item = InvoiceRecognizeItem(
                row_index=row_index,
                drug_name=raw.get("drug_name") or None,
                specification=raw.get("specification") or None,
                manufacturer=raw.get("manufacturer") or None,
                approval_number=raw.get("approval_number") or None,
                batch_number=raw.get("batch_number") or None,
                expire_date=_parse_date(raw.get("expire_date")),
                quantity=_parse_decimal_str(raw.get("quantity")),
                unit_price=_parse_decimal_str(raw.get("unit_price")),
                amount=_parse_decimal_str(raw.get("amount")),
                confidence=0.75,
            )
            recognized_items.append(item)

        # ── Step 5: match master data ─────────────────────────────────────────
        supplier_candidates: list[SupplierCandidate] = []
        matched_supplier_id: int | None = None

        if match_master_data and supplier_name:
            supplier_candidates = await matching_service.match_supplier(supplier_name, limit=5)
            if supplier_candidates and supplier_candidates[0].confidence >= 0.7:
                matched_supplier_id = supplier_candidates[0].supplier_id

        if match_master_data:
            for item in recognized_items:
                if item.drug_name:
                    match_input = DrugMatchInputItem(
                        row_index=item.row_index,
                        drug_name=item.drug_name,
                        specification=item.specification,
                        manufacturer=item.manufacturer,
                        approval_number=item.approval_number,
                    )
                    candidates = await matching_service.match_drugs(match_input, limit=5)
                    item.drug_candidates = candidates
                    if candidates and candidates[0].confidence >= 0.8:
                        item.matched_drug_id = candidates[0].drug_id

        # ── Step 6: quality warnings ──────────────────────────────────────────
        warnings: list[QualityWarning] = []
        _add_quality_warnings(recognized_items, total_amount, warnings)

        overall_confidence = _compute_confidence(recognized_items, supplier_candidates)

        result = InvoiceRecognizeResult(
            recognized_supplier_name=supplier_name,
            supplier_candidates=supplier_candidates,
            matched_supplier_id=matched_supplier_id,
            invoice_no=invoice_no,
            invoice_date=invoice_date,
            total_amount=total_amount,
            confidence=overall_confidence,
            items=recognized_items,
            warnings=warnings,
        )

        finished_at = datetime.now(timezone.utc)
        duration_ms = int((finished_at - started_at).total_seconds() * 1000)

        # Persist completed job
        await pool.execute(
            """
            UPDATE ai.invoice_recognition_job SET
              status='COMPLETED',
              recognized_supplier_name=$2,
              matched_supplier_id=$3,
              invoice_no=$4,
              invoice_date=$5,
              confidence=$6,
              result_json=$7,
              raw_response_json=$8,
              warnings_json=$9,
              finished_at=$10,
              duration_ms=$11
            WHERE id=$1
            """,
            job_id,
            supplier_name,
            matched_supplier_id,
            invoice_no,
            invoice_date,
            overall_confidence,
            json.dumps(result.model_dump(mode="json"), ensure_ascii=False),
            json.dumps(raw_response, ensure_ascii=False) if return_raw_response else None,
            json.dumps([w.model_dump() for w in warnings], ensure_ascii=False),
            finished_at,
            duration_ms,
        )

        # Cache items
        for item in recognized_items:
            try:
                await pool.execute(
                    """
                    INSERT INTO ai.invoice_recognition_item_cache
                      (job_id, row_no, drug_name, specification, manufacturer,
                       batch_number, expire_date, quantity, unit_price, amount,
                       matched_drug_id, confidence, candidates_json, warnings_json)
                    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
                    ON CONFLICT (job_id, row_no) DO NOTHING
                    """,
                    job_id,
                    item.row_index,
                    item.drug_name,
                    item.specification,
                    item.manufacturer,
                    item.batch_number,
                    item.expire_date,
                    float(item.quantity) if item.quantity else None,
                    float(item.unit_price) if item.unit_price else None,
                    float(item.amount) if item.amount else None,
                    item.matched_drug_id,
                    item.confidence,
                    json.dumps([c.model_dump(mode="json") for c in item.drug_candidates], ensure_ascii=False),
                    json.dumps([w.model_dump() for w in item.warnings], ensure_ascii=False),
                )
            except Exception:
                pass

        return InvoiceRecognizeResponse(
            request_id=request_id,
            status="COMPLETED",
            result=result,
            raw_response=raw_response if return_raw_response else None,
            started_at=started_at,
            finished_at=finished_at,
            duration_ms=duration_ms,
        )

    except Exception as exc:
        finished_at = datetime.now(timezone.utc)
        duration_ms = int((finished_at - started_at).total_seconds() * 1000)
        error_code = "LLM_ERROR"
        error_message = str(exc)

        await pool.execute(
            """
            UPDATE ai.invoice_recognition_job SET
              status='FAILED', error_code=$2, error_message=$3,
              finished_at=$4, duration_ms=$5
            WHERE id=$1
            """,
            job_id, error_code, error_message, finished_at, duration_ms,
        )

        return InvoiceRecognizeResponse(
            request_id=request_id,
            status="FAILED",
            error_code=error_code,
            error_message=error_message,
            started_at=started_at,
            finished_at=finished_at,
            duration_ms=duration_ms,
        )


async def get_job(request_id: str) -> InvoiceRecognizeResponse | None:
    pool = get_pool()
    row = await pool.fetchrow(
        "SELECT * FROM ai.invoice_recognition_job WHERE request_id=$1 AND deleted_at IS NULL",
        request_id,
    )
    if row is None:
        return None

    result: InvoiceRecognizeResult | None = None
    if row["result_json"]:
        try:
            result = InvoiceRecognizeResult.model_validate(
                json.loads(row["result_json"])
            )
        except Exception:
            pass

    raw = None
    if row["raw_response_json"]:
        try:
            raw = json.loads(row["raw_response_json"])
        except Exception:
            pass

    duration_ms = row["duration_ms"]

    return InvoiceRecognizeResponse(
        request_id=row["request_id"],
        status=row["status"],
        result=result,
        raw_response=raw,
        error_code=row["error_code"],
        error_message=row["error_message"],
        started_at=row["started_at"],
        finished_at=row["finished_at"],
        duration_ms=int(duration_ms) if duration_ms is not None else None,
    )


def _compute_confidence(
    items: list[InvoiceRecognizeItem],
    supplier_candidates: list[SupplierCandidate],
) -> float:
    scores = [item.confidence or 0.0 for item in items if item.confidence is not None]
    if supplier_candidates:
        scores.append(supplier_candidates[0].confidence)
    return round(sum(scores) / len(scores), 4) if scores else 0.5


def _add_quality_warnings(
    items: list[InvoiceRecognizeItem],
    total_amount: str | None,
    warnings: list[QualityWarning],
) -> None:
    # Check amount consistency: quantity * unit_price ≈ amount
    computed_total = 0.0
    for item in items:
        try:
            qty = float(item.quantity or "0")
            price = float(item.unit_price or "0")
            amt = float(item.amount or "0")
            expected = round(qty * price, 2)
            if qty > 0 and price > 0 and abs(expected - amt) > 0.02:
                item.warnings.append(
                    QualityWarning(
                        level="MEDIUM",
                        code="AMOUNT_MISMATCH",
                        field=f"items[{item.row_index}].amount",
                        message=f"明细金额 {amt} 与数量×单价 {expected} 不一致，请人工确认",
                        suggestion="请核对数量、单价和金额",
                    )
                )
            computed_total += amt if amt else expected
        except (ValueError, TypeError):
            pass

    # Check total amount
    if total_amount:
        try:
            declared = float(total_amount)
            if abs(declared - computed_total) > 0.1:
                warnings.append(
                    QualityWarning(
                        level="HIGH",
                        code="TOTAL_AMOUNT_MISMATCH",
                        field="total_amount",
                        message=f"发票合计金额 {declared} 与明细合计 {round(computed_total,2)} 不符",
                        suggestion="请核对发票总金额",
                    )
                )
        except (ValueError, TypeError):
            pass

    # Check missing fields
    for item in items:
        if not item.drug_name:
            item.warnings.append(
                QualityWarning(
                    level="HIGH",
                    code="MISSING_DRUG_NAME",
                    field=f"items[{item.row_index}].drug_name",
                    message="药品名称未识别",
                )
            )
        if not item.expire_date:
            item.warnings.append(
                QualityWarning(
                    level="LOW",
                    code="MISSING_EXPIRE_DATE",
                    field=f"items[{item.row_index}].expire_date",
                    message="有效期未识别，请人工填写",
                )
            )


def check_quality(
    result: InvoiceRecognizeResult, strict: bool
) -> tuple[bool, list[QualityWarning]]:
    warnings: list[QualityWarning] = list(result.warnings)

    for item in result.items:
        warnings.extend(item.warnings)
        if strict and item.matched_drug_id is None and item.drug_name:
            warnings.append(
                QualityWarning(
                    level="HIGH",
                    code="NO_DRUG_MATCH",
                    field=f"items[{item.row_index}].matched_drug_id",
                    message=f"药品[{item.drug_name}]未匹配到系统药品，请人工确认",
                )
            )

    if strict and not result.matched_supplier_id:
        warnings.append(
            QualityWarning(
                level="HIGH",
                code="NO_SUPPLIER_MATCH",
                field="matched_supplier_id",
                message="供应商未匹配，请人工确认",
            )
        )

    high_count = sum(1 for w in warnings if w.level == "HIGH")
    passed = high_count == 0
    return passed, warnings
