from fastapi import APIRouter

from app.models.common import ok
from app.models.drug import DrugMatchRequest
from app.models.supplier import SupplierMatchRequest
from app.request_id import new_request_id
from app.services import matching_service

router = APIRouter()


@router.post("/drugs/match")
async def match_drugs(body: DrugMatchRequest):
    request_id = new_request_id()
    result_items = []
    for input_item in body.items:
        candidates = await matching_service.match_drugs(input_item, body.limit_per_item)
        best = candidates[0] if candidates else None
        result_items.append(
            {
                "row_index": input_item.row_index,
                "candidates": [c.model_dump(mode="json") for c in candidates],
                "best_candidate": best.model_dump(mode="json") if best else None,
                "warnings": [],
            }
        )
    return ok({"items": result_items}, request_id=request_id)


@router.post("/suppliers/match")
async def match_supplier(body: SupplierMatchRequest):
    request_id = new_request_id()
    candidates = await matching_service.match_supplier(body.supplier_name, body.limit)
    best = candidates[0] if candidates else None
    return ok(
        {
            "recognized_supplier_name": body.supplier_name,
            "best_candidate": best.model_dump(mode="json") if best else None,
            "candidates": [c.model_dump(mode="json") for c in candidates],
        },
        request_id=request_id,
    )
