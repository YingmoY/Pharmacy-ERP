from fastapi import APIRouter, File, Form, UploadFile, HTTPException

from app.models.common import ok, err
from app.models.invoice import InvoiceQualityCheckRequest
from app.request_id import new_request_id
from app.services import invoice_service

router = APIRouter()

ALLOWED_CONTENT_TYPES = {
    "image/jpeg", "image/jpg", "image/png", "image/webp",
    "application/pdf",
}
MAX_FILE_SIZE = 20 * 1024 * 1024  # 20 MB


@router.post("/invoices/recognize")
async def recognize_invoice(
    file: UploadFile = File(...),
    erp_request_id: str | None = Form(None),
    erp_file_id: str | None = Form(None),
    async_mode: bool = Form(False, alias="async"),
    match_master_data: bool = Form(True),
    return_raw_response: bool = Form(False),
):
    request_id = new_request_id()

    ct = (file.content_type or "").lower()
    if ct not in ALLOWED_CONTENT_TYPES:
        return err(415, f"不支持的文件类型: {ct}", request_id, "UNSUPPORTED_MEDIA_TYPE")

    content = await file.read()
    if len(content) > MAX_FILE_SIZE:
        return err(413, "文件过大，最大支持 20MB", request_id, "PAYLOAD_TOO_LARGE")

    response = await invoice_service.recognize_invoice(
        file_content=content,
        file_name=file.filename or "invoice",
        content_type=ct,
        erp_request_id=erp_request_id,
        erp_file_id=erp_file_id,
        match_master_data=match_master_data,
        return_raw_response=return_raw_response,
    )
    return ok(response.model_dump(mode="json"), request_id=response.request_id)


@router.get("/invoices/jobs/{request_id}")
async def get_invoice_job(request_id: str):
    api_request_id = new_request_id()
    response = await invoice_service.get_job(request_id)
    if response is None:
        return err(404, "识别任务不存在", api_request_id, "NOT_FOUND")
    return ok(response.model_dump(mode="json"), request_id=api_request_id)


@router.post("/invoices/quality-check")
async def quality_check(body: InvoiceQualityCheckRequest):
    request_id = new_request_id()
    passed, warnings = invoice_service.check_quality(body.result, body.strict)
    return ok(
        {
            "passed": passed,
            "warnings": [w.model_dump() for w in warnings],
        },
        request_id=request_id,
    )
