import asyncio

from fastapi import APIRouter, Request, WebSocket, WebSocketDisconnect

from app.models.common import ok, err
from app.models.drug import DrugSearchRequest, DrugSearchFeedbackRequest, DrugRecommendRequest, DrugSearchFilters
from app.request_id import new_request_id
from app.services import drug_search_service

router = APIRouter()

_WS_PROGRESS_HINTS = [
    "正在搜索库存候选药品...",
    "AI正在分析描述...",
    "匹配症状对应药品...",
    "AI深度推理中...",
    "即将完成...",
    "AI深度推理中...",
]


def _client_ip(request: Request) -> str | None:
    forwarded = request.headers.get("x-forwarded-for")
    if forwarded:
        return forwarded.split(",")[0].strip()
    return request.client.host if request.client else None


@router.post("/drugs/search")
async def search_drugs(body: DrugSearchRequest, request: Request):
    request_id = new_request_id()
    result = await drug_search_service.search_drugs(
        query=body.query,
        search_mode=body.search_mode,
        limit=body.limit,
        offset=body.offset,
        filters=body.filters,
        request_id=request_id,
        client_ip=_client_ip(request),
    )
    return ok(result.model_dump(mode="json"), request_id=request_id)


@router.get("/drugs/{drug_id}")
async def get_drug_detail(drug_id: int):
    request_id = new_request_id()
    item = await drug_search_service.get_drug_detail(drug_id)
    if item is None:
        return err(404, "药品不存在", request_id=request_id, error_code="NOT_FOUND")
    return ok(item.model_dump(mode="json"), request_id=request_id)


@router.post("/drugs/recommend")
async def recommend_drugs(body: DrugRecommendRequest, request: Request):
    """HTTP fallback for non-CDN environments. Prefer the WebSocket endpoint."""
    request_id = new_request_id()
    body.filters.only_available = True
    result = await drug_search_service.ai_recommend_drugs(
        query=body.query,
        filters=body.filters,
        limit=body.limit,
    )
    return ok(result.model_dump(mode="json"), request_id=request_id)


@router.websocket("/drugs/recommend/ws")
async def recommend_drugs_ws(websocket: WebSocket):
    """WebSocket endpoint for AI drug recommendation.

    Protocol:
      client → server: {"query": str, "limit": int, "filters": {"only_available": bool}}
      server → client: {"type": "progress", "message": str}   (every ~3s while thinking)
      server → client: {"type": "result",   "data": DrugRecommendResponse}
      server → client: {"type": "error",    "message": str}   (on failure)
    """
    await websocket.accept()
    recommend_task: asyncio.Task | None = None
    try:
        # Wait up to 10s for the client to send the request payload
        data = await asyncio.wait_for(websocket.receive_json(), timeout=10.0)

        query: str = str(data.get("query") or "").strip()
        limit: int = int(data.get("limit") or 10)
        filters_raw: dict = data.get("filters") or {}
        if not query:
            await websocket.send_json({"type": "error", "message": "query 不能为空"})
            return

        filters = DrugSearchFilters(**{k: v for k, v in filters_raw.items() if k in DrugSearchFilters.model_fields})
        filters.only_available = True  # always restrict to in-stock for POS

        # Launch the heavy LLM work as a background task
        recommend_task = asyncio.create_task(
            drug_search_service.ai_recommend_drugs(query, filters, limit)
        )

        # Send progress hints every 3s while the task is running
        hint_idx = 0
        while not recommend_task.done():
            try:
                await asyncio.wait_for(asyncio.shield(recommend_task), timeout=3.0)
            except asyncio.TimeoutError:
                hint = _WS_PROGRESS_HINTS[hint_idx % len(_WS_PROGRESS_HINTS)]
                hint_idx += 1
                await websocket.send_json({"type": "progress", "message": hint})
            except Exception:
                break  # task raised; handled below

        result = recommend_task.result()
        await websocket.send_json({"type": "result", "data": result.model_dump(mode="json")})

    except WebSocketDisconnect:
        if recommend_task and not recommend_task.done():
            recommend_task.cancel()
    except Exception as exc:
        try:
            await websocket.send_json({"type": "error", "message": str(exc)})
        except Exception:
            pass
        if recommend_task and not recommend_task.done():
            recommend_task.cancel()


@router.post("/drugs/search-feedback")
async def submit_feedback(body: DrugSearchFeedbackRequest, request: Request):
    request_id = new_request_id()
    await drug_search_service.save_feedback(
        request_id=body.request_id,
        query=body.query,
        selected_drug_id=body.selected_drug_id,
        feedback_type=body.feedback_type,
        feedback_note=body.feedback_note,
        client_ip=_client_ip(request),
    )
    return ok({}, request_id=request_id)
