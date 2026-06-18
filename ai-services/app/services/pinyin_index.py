"""In-memory pinyin index for drug search.

Loads drug names from DB at startup (and refreshes every REFRESH_INTERVAL_S),
then provides fast pinyin abbreviation + full-pinyin matching.

Design:
  - ASCII-only queries are treated as possible pinyin input.
  - Matches abbreviation (首字母, e.g. "abo" → 阿波罗) and full pinyin prefix
    (e.g. "abo" matches "abokeluo").
  - Returns a set of drug_ids so the caller can add a SQL `WHERE drug_id = ANY(...)`.
"""

import asyncio
import logging
import re
import time
from typing import Optional

from pypinyin import lazy_pinyin, Style

logger = logging.getLogger(__name__)

REFRESH_INTERVAL_S = 300  # rebuild index every 5 minutes

_index: dict[int, tuple[str, str]] = {}  # drug_id -> (abbrev, full_pinyin)
_last_built: float = 0.0
_lock = asyncio.Lock()


def _to_pinyin(text: str) -> tuple[str, str]:
    """Return (abbreviation, full_pinyin) for a Chinese string."""
    chars = lazy_pinyin(text, style=Style.NORMAL, errors="ignore")
    full = "".join(chars)
    abbrev = "".join(c[0] for c in chars if c)
    return abbrev, full


def _build_entry(common_name: str, trade_name: str) -> tuple[str, str]:
    combined = f"{common_name} {trade_name}".strip()
    abbrevs: list[str] = []
    fulls: list[str] = []
    for part in (common_name, trade_name):
        if part:
            a, f = _to_pinyin(part)
            abbrevs.append(a)
            fulls.append(f)
    return " ".join(filter(None, abbrevs)), " ".join(filter(None, fulls))


async def _rebuild(pool) -> None:
    global _index, _last_built
    try:
        async with pool.acquire() as conn:
            rows = await conn.fetch(
                "SELECT drug_id, common_name, COALESCE(trade_name, '') AS trade_name "
                "FROM ai.v_drug_search_source"
            )
        new_index: dict[int, tuple[str, str]] = {}
        for row in rows:
            drug_id = row["drug_id"]
            abbrev, full = _build_entry(row["common_name"], row["trade_name"])
            new_index[drug_id] = (abbrev, full)
        _index = new_index
        _last_built = time.monotonic()
        logger.info("Pinyin index built: %d drugs", len(_index))
    except Exception:
        logger.exception("Failed to build pinyin index")


async def ensure_index(pool) -> None:
    """Refresh index if stale; called before each pinyin search."""
    if time.monotonic() - _last_built < REFRESH_INTERVAL_S:
        return
    async with _lock:
        if time.monotonic() - _last_built < REFRESH_INTERVAL_S:
            return
        await _rebuild(pool)


_PINYIN_RE = re.compile(r"^[a-zA-Z]+$")


def is_pinyin_query(query: str) -> bool:
    """Return True if the query looks like a pinyin input (ASCII letters only)."""
    return bool(_PINYIN_RE.match(query))


def search_pinyin(query: str) -> set[int]:
    """Return drug_ids whose pinyin abbreviation or full pinyin starts with query."""
    q = query.lower()
    matched: set[int] = set()
    for drug_id, (abbrev, full) in _index.items():
        if abbrev.startswith(q) or full.startswith(q):
            matched.add(drug_id)
    return matched
