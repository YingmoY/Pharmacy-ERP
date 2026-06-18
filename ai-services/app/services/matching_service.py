"""Drug and supplier candidate matching using pg_trgm similarity."""
import asyncpg

from app.database import get_pool
from app.models.drug import DrugMatchInputItem, DrugSearchItem, DrugInventorySummary
from app.models.invoice import DrugCandidateInvoice
from app.models.supplier import SupplierCandidate


async def match_drugs(
    item: DrugMatchInputItem, limit: int
) -> list[DrugCandidateInvoice]:
    pool = get_pool()

    # Build a combined search string from all provided fields
    parts = [item.drug_name]
    if item.specification:
        parts.append(item.specification)
    if item.manufacturer:
        parts.append(item.manufacturer)
    search_str = " ".join(parts)

    async with pool.acquire() as conn:
        # Try exact barcode/approval_number match first
        exact_rows: list[asyncpg.Record] = []
        if item.barcode:
            r = await conn.fetchrow(
                "SELECT * FROM ai.v_invoice_drug_match_source WHERE barcode = $1 LIMIT 1",
                item.barcode,
            )
            if r:
                exact_rows.append(r)
        if not exact_rows and item.approval_number:
            r = await conn.fetchrow(
                "SELECT * FROM ai.v_invoice_drug_match_source WHERE approval_number = $1 LIMIT 1",
                item.approval_number,
            )
            if r:
                exact_rows.append(r)

        # Fuzzy match by similarity
        fuzzy_rows = await conn.fetch(
            """
            SELECT *, similarity(search_text, $1) AS sim
            FROM ai.v_invoice_drug_match_source
            WHERE similarity(search_text, $1) > 0.15
            ORDER BY sim DESC
            LIMIT $2
            """,
            search_str,
            limit,
        )

    candidates: list[DrugCandidateInvoice] = []
    seen: set[int] = set()

    for r in exact_rows:
        did = r["drug_id"]
        seen.add(did)
        candidates.append(
            DrugCandidateInvoice(
                drug_id=did,
                drug_code=r["drug_code"],
                common_name=r["common_name"],
                trade_name=r["trade_name"],
                specification=r["specification"],
                dosage_form=r["dosage_form"],
                manufacturer=r["manufacturer"],
                approval_number=r["approval_number"],
                barcode=r["barcode"],
                unit=r["unit"],
                retail_price=str(r["retail_price"]) if r["retail_price"] else None,
                purchase_price=str(r["purchase_price"]) if r["purchase_price"] else None,
                is_prescription=r["is_prescription"],
                is_medicare=r["is_medicare"],
                score=1.0,
                confidence=1.0,
                candidate_type="EXACT",
                match_reason="条形码/批准文号精确匹配",
            )
        )

    for r in fuzzy_rows:
        did = r["drug_id"]
        if did in seen:
            continue
        seen.add(did)
        sim = float(r["sim"])
        candidates.append(
            DrugCandidateInvoice(
                drug_id=did,
                drug_code=r["drug_code"],
                common_name=r["common_name"],
                trade_name=r["trade_name"],
                specification=r["specification"],
                dosage_form=r["dosage_form"],
                manufacturer=r["manufacturer"],
                approval_number=r["approval_number"],
                barcode=r["barcode"],
                unit=r["unit"],
                retail_price=str(r["retail_price"]) if r["retail_price"] else None,
                purchase_price=str(r["purchase_price"]) if r["purchase_price"] else None,
                is_prescription=r["is_prescription"],
                is_medicare=r["is_medicare"],
                score=sim,
                confidence=sim,
                candidate_type="FUZZY",
                match_reason=f"模糊匹配相似度 {sim:.2f}",
            )
        )
        if len(candidates) >= limit:
            break

    return candidates[:limit]


async def match_supplier(name: str, limit: int) -> list[SupplierCandidate]:
    pool = get_pool()
    async with pool.acquire() as conn:
        rows = await conn.fetch(
            """
            SELECT *,
              GREATEST(
                similarity(name, $1),
                similarity(COALESCE(license_no, ''), $1)
              ) AS sim
            FROM ai.v_active_supplier_source
            WHERE similarity(name, $1) > 0.1
               OR name ILIKE $2
            ORDER BY sim DESC
            LIMIT $3
            """,
            name,
            f"%{name}%",
            limit,
        )

    candidates: list[SupplierCandidate] = []
    for r in rows:
        sim = float(r["sim"])
        reason = "名称模糊匹配"
        if r["name"] == name:
            sim = 1.0
            reason = "名称精确匹配"
        elif name in r["name"] or r["name"] in name:
            sim = max(sim, 0.85)
            reason = "名称包含匹配"
        candidates.append(
            SupplierCandidate(
                supplier_id=r["supplier_id"],
                supplier_code=r["supplier_code"],
                name=r["name"],
                license_no=r["license_no"],
                contact_name=r["contact_name"],
                contact_phone=r["contact_phone"],
                confidence=round(sim, 4),
                match_reason=reason,
            )
        )

    return candidates
