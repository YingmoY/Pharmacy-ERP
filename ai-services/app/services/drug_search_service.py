"""Drug search service — queries ai.v_drug_search_source using pg_trgm + ILIKE + pinyin."""
import time
from datetime import date

import asyncpg

from app.database import get_pool
from app.models.drug import (
    DrugSearchFilters,
    DrugSearchItem,
    DrugInventorySummary,
    DrugSearchResponse,
    DrugRecommendResponse,
)
from app.services import pinyin_index, llm_client


def _row_to_item(row: asyncpg.Record, score: float, reason: str) -> DrugSearchItem:
    inv = DrugInventorySummary(
        available_qty=row["available_qty"],
        in_stock_qty=row["in_stock_qty"],
        reserved_qty=row["reserved_qty"],
        pending_qty=row["pending_qty"],
        abnormal_qty=row["abnormal_qty"],
        near_expire_available_qty=row["near_expire_available_qty"],
        nearest_expire_date=row["nearest_expire_date"],
    )
    return DrugSearchItem(
        drug_id=row["drug_id"],
        drug_code=row["drug_code"],
        common_name=row["common_name"],
        trade_name=row["trade_name"],
        specification=row["specification"],
        dosage_form=row["dosage_form"],
        manufacturer=row["manufacturer"],
        approval_number=row["approval_number"],
        barcode=row["barcode"],
        unit=row["unit"],
        retail_price=str(row["retail_price"]) if row["retail_price"] is not None else None,
        purchase_price=str(row["purchase_price"]) if row["purchase_price"] is not None else None,
        is_prescription=row["is_prescription"],
        is_medicare=row["is_medicare"],
        inventory=inv,
        score=score,
        match_reason=reason,
        highlights=[],
    )


def _build_filter_clauses(filters: DrugSearchFilters) -> tuple[list[str], list]:
    clauses: list[str] = []
    params: list = []

    if filters.only_available:
        clauses.append("available_qty > 0")
    if filters.is_prescription is not None:
        params.append(filters.is_prescription)
        clauses.append(f"is_prescription = ${len(params)}")
    if filters.is_medicare is not None:
        params.append(filters.is_medicare)
        clauses.append(f"is_medicare = ${len(params)}")
    if filters.manufacturer:
        params.append(f"%{filters.manufacturer}%")
        clauses.append(f"manufacturer ILIKE ${len(params)}")
    if filters.dosage_form:
        params.append(f"%{filters.dosage_form}%")
        clauses.append(f"dosage_form ILIKE ${len(params)}")
    if filters.storage_condition:
        params.append(f"%{filters.storage_condition}%")
        clauses.append(f"storage_condition ILIKE ${len(params)}")
    if filters.near_expire_only:
        clauses.append("near_expire_available_qty > 0")

    return clauses, params


async def search_drugs(
    query: str,
    search_mode: str,
    limit: int,
    offset: int,
    filters: DrugSearchFilters,
    request_id: str | None = None,
    client_ip: str | None = None,
) -> DrugSearchResponse:
    pool = get_pool()
    start = time.monotonic()

    filter_clauses, filter_params = _build_filter_clauses(filters)
    where_base = " AND ".join(filter_clauses) if filter_clauses else "TRUE"

    rows: list[asyncpg.Record] = []
    score_map: dict[int, tuple[float, str]] = {}

    # Detect pinyin input and resolve to drug_ids via in-memory index
    pinyin_ids: list[int] = []
    if pinyin_index.is_pinyin_query(query):
        await pinyin_index.ensure_index(pool)
        pinyin_ids = list(pinyin_index.search_pinyin(query))

    async with pool.acquire() as conn:
        # --- Pinyin search (always run when query looks like pinyin) ---
        if pinyin_ids:
            base_idx = len(filter_params) + 1
            py_params = filter_params + [pinyin_ids, limit + offset]
            py_sql = f"""
                SELECT * FROM ai.v_drug_search_source
                WHERE ({where_base})
                  AND drug_id = ANY(${base_idx}::int[])
                ORDER BY available_qty DESC
                LIMIT ${base_idx + 1}
            """
            py_rows = await conn.fetch(py_sql, *py_params)
            for r in py_rows:
                did = r["drug_id"]
                if did not in score_map:
                    score_map[did] = (0.9, "拼音匹配")
                rows.append(r)

        if search_mode in ("KEYWORD", "HYBRID"):
            # Param index starts after filter_params
            base_idx = len(filter_params) + 1
            kw_params = filter_params + [f"%{query}%", limit + offset]
            kw_sql = f"""
                SELECT * FROM ai.v_drug_search_source
                WHERE ({where_base})
                  AND search_text ILIKE ${base_idx}
                ORDER BY
                  CASE WHEN common_name ILIKE ${base_idx} THEN 0 ELSE 1 END,
                  available_qty DESC
                LIMIT ${base_idx + 1}
            """
            kw_rows = await conn.fetch(kw_sql, *kw_params)
            for r in kw_rows:
                did = r["drug_id"]
                if did not in score_map:
                    score_map[did] = (0.8, "关键词匹配")
                rows.append(r)

        if search_mode in ("FUZZY", "HYBRID", "SEMANTIC"):
            base_idx = len(filter_params) + 1
            fz_params = filter_params + [query, limit + offset]
            fz_sql = f"""
                SELECT *, similarity(search_text, ${base_idx}) AS sim
                FROM ai.v_drug_search_source
                WHERE ({where_base})
                  AND similarity(search_text, ${base_idx}) > 0.1
                ORDER BY sim DESC, available_qty DESC
                LIMIT ${base_idx + 1}
            """
            fz_rows = await conn.fetch(fz_sql, *fz_params)
            for r in fz_rows:
                did = r["drug_id"]
                sim = float(r["sim"])
                if did not in score_map or score_map[did][0] < sim:
                    score_map[did] = (sim, "模糊匹配")
                if not any(x["drug_id"] == did for x in rows):
                    rows.append(r)

    # Deduplicate, sort by score, apply offset/limit
    seen: set[int] = set()
    unique_rows: list[asyncpg.Record] = []
    for r in rows:
        did = r["drug_id"]
        if did not in seen:
            seen.add(did)
            unique_rows.append(r)

    unique_rows.sort(key=lambda r: score_map.get(r["drug_id"], (0.0, ""))[0], reverse=True)
    total = len(unique_rows)
    paged = unique_rows[offset: offset + limit]

    items = [
        _row_to_item(r, *score_map.get(r["drug_id"], (0.5, "匹配")))
        for r in paged
    ]

    latency_ms = int((time.monotonic() - start) * 1000)

    # Log async (fire-and-forget)
    pool2 = get_pool()
    try:
        async with pool2.acquire() as conn2:
            await conn2.execute(
                """
                INSERT INTO ai.drug_search_log
                  (request_id, query_text, search_mode, result_count,
                   top_result_drug_id, latency_ms, client_ip)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                """,
                request_id,
                query,
                search_mode,
                total,
                items[0].drug_id if items else None,
                latency_ms,
                client_ip,
            )
    except Exception:
        pass

    return DrugSearchResponse(
        query=query,
        normalized_query=query,
        search_mode=search_mode,
        total=total,
        items=items,
        suggestions=[],
    )


async def get_drug_detail(drug_id: int) -> DrugSearchItem | None:
    pool = get_pool()
    async with pool.acquire() as conn:
        row = await conn.fetchrow(
            "SELECT * FROM ai.v_drug_search_source WHERE drug_id = $1",
            drug_id,
        )
    if row is None:
        return None
    return _row_to_item(row, 1.0, "精确 ID 查询")


async def save_feedback(
    request_id: str | None,
    query: str,
    selected_drug_id: int | None,
    feedback_type: str,
    feedback_note: str | None,
    client_ip: str | None,
) -> None:
    pool = get_pool()
    async with pool.acquire() as conn:
        await conn.execute(
            """
            INSERT INTO ai.drug_search_feedback
              (request_id, query_text, selected_drug_id, feedback_type, feedback_note, client_ip)
            VALUES ($1, $2, $3, $4, $5, $6)
            """,
            request_id,
            query,
            selected_drug_id,
            feedback_type,
            feedback_note,
            client_ip,
        )


_CANDIDATE_LIMIT = 80  # how many drugs to pull from DB as LLM input


async def ai_recommend_drugs(
    query: str,
    filters: DrugSearchFilters,
    limit: int,
) -> DrugRecommendResponse:
    """Fetch candidate drugs from inventory, then let the LLM rank which ones fit the query.

    Flow:
    1. Broad DB search (keyword + fuzzy) for candidates from live inventory.
    2. If too few candidates, supplement with top-stocked drugs so LLM always has choices.
    3. LLM receives the candidate list and picks the most appropriate ones for the user's query.
    4. Return LLM-ordered results, all guaranteed to be in stock.
    """
    pool = get_pool()
    filter_clauses, filter_params = _build_filter_clauses(filters)
    where_base = " AND ".join(filter_clauses) if filter_clauses else "TRUE"

    candidate_map: dict[int, asyncpg.Record] = {}

    async with pool.acquire() as conn:
        base_idx = len(filter_params) + 1

        # Keyword search
        kw_rows = await conn.fetch(
            f"""SELECT * FROM ai.v_drug_search_source
                WHERE ({where_base}) AND search_text ILIKE ${base_idx}
                ORDER BY available_qty DESC LIMIT ${base_idx + 1}""",
            *filter_params, f"%{query}%", _CANDIDATE_LIMIT,
        )
        for r in kw_rows:
            candidate_map[r["drug_id"]] = r

        # Fuzzy / pg_trgm search (low threshold to cast a wide net)
        fz_rows = await conn.fetch(
            f"""SELECT * FROM ai.v_drug_search_source
                WHERE ({where_base}) AND similarity(search_text, ${base_idx}) > 0.05
                ORDER BY similarity(search_text, ${base_idx}) DESC, available_qty DESC
                LIMIT ${base_idx + 1}""",
            *filter_params, query, _CANDIDATE_LIMIT,
        )
        for r in fz_rows:
            candidate_map[r["drug_id"]] = r

        # If still not enough candidates, supplement with top-stocked drugs
        if len(candidate_map) < 15:
            top_rows = await conn.fetch(
                f"""SELECT * FROM ai.v_drug_search_source
                    WHERE ({where_base})
                    ORDER BY available_qty DESC LIMIT ${base_idx}""",
                *filter_params, _CANDIDATE_LIMIT,
            )
            for r in top_rows:
                if r["drug_id"] not in candidate_map:
                    candidate_map[r["drug_id"]] = r
                if len(candidate_map) >= _CANDIDATE_LIMIT:
                    break

    candidates = list(candidate_map.values())

    # Build a compact list for the LLM (id + name + spec + form)
    llm_candidates = [
        {
            "drug_id": r["drug_id"],
            "common_name": r["common_name"],
            "specification": r["specification"] or "",
            "dosage_form": r["dosage_form"] or "",
        }
        for r in candidates
    ]

    # Ask LLM to rank/select
    ranked = await llm_client.llm_rank_drugs_for_query(query, llm_candidates)
    selected: list[dict] = ranked.get("selected") or []
    explanation: str = ranked.get("explanation") or ""

    # Map selected drug_ids back to full records, preserving LLM order
    selected_ids: list[int] = [s["drug_id"] for s in selected if isinstance(s.get("drug_id"), int)]
    reason_map: dict[int, str] = {
        s["drug_id"]: s.get("reason") or "AI推荐"
        for s in selected
        if isinstance(s.get("drug_id"), int)
    }

    # If LLM returned nothing or returned unknown IDs, fall back to top candidates
    valid_ids = [did for did in selected_ids if did in candidate_map]
    if not valid_ids:
        valid_ids = [r["drug_id"] for r in candidates[:limit]]
        reason_map = {did: "库存推荐" for did in valid_ids}

    result_ids = valid_ids[:limit]
    items = [
        _row_to_item(candidate_map[did], 1.0, reason_map.get(did, "AI推荐"))
        for did in result_ids
    ]

    return DrugRecommendResponse(
        query=query,
        explanation=explanation,
        terms=[],
        total=len(result_ids),
        items=items,
    )
