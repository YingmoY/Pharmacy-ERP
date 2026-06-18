from datetime import datetime, timezone

from fastapi import APIRouter

from app.config import settings
from app.database import get_pool
from app.models.common import ok
from app.request_id import new_request_id

router = APIRouter()


@router.get("/health")
async def health_check():
    request_id = new_request_id()
    db_status = "READY"
    try:
        pool = get_pool()
        await pool.fetchval("SELECT 1")
    except Exception:
        db_status = "UNAVAILABLE"

    llm_status = "READY"
    try:
        from openai import AsyncOpenAI
        client = AsyncOpenAI(base_url=settings.llm_base_url, api_key=settings.llm_api_key)
        await client.models.list()
    except Exception:
        llm_status = "DEGRADED"

    overall = "UP" if db_status == "READY" else "DOWN"
    if llm_status != "READY" and overall == "UP":
        overall = "DEGRADED"

    return ok(
        {
            "status": overall,
            "service_time": datetime.now(timezone.utc).isoformat(),
            "version": settings.service_version,
            "models": [
                {
                    "module": "DRUG_SEARCH",
                    "provider": "LOCAL_SEARCH",
                    "model_name": "hybrid-keyword-fuzzy",
                    "status": "READY" if db_status == "READY" else "UNAVAILABLE",
                },
                {
                    "module": "INVOICE_RECOGNITION",
                    "provider": "LOCAL_LLM",
                    "model_name": settings.llm_model,
                    "status": llm_status,
                },
            ],
        },
        request_id=request_id,
    )
