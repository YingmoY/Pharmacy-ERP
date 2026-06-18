from datetime import date, datetime
from typing import Any, Literal
from pydantic import BaseModel, Field

from app.models.supplier import SupplierCandidate


class QualityWarning(BaseModel):
    level: Literal["LOW", "MEDIUM", "HIGH"]
    code: str | None = None
    field: str | None = None
    message: str
    suggestion: str | None = None


class DrugCandidateInvoice(BaseModel):
    drug_id: int
    drug_code: str
    common_name: str
    trade_name: str | None = None
    specification: str
    dosage_form: str | None = None
    manufacturer: str
    approval_number: str | None = None
    barcode: str | None = None
    unit: str | None = None
    retail_price: str | None = None
    purchase_price: str | None = None
    is_prescription: bool = False
    is_medicare: bool = False
    inventory: dict | None = None
    score: float = 0.0
    match_reason: str | None = None
    highlights: list[str] = []
    confidence: float = 0.0
    candidate_type: Literal["EXACT", "FUZZY", "SEMANTIC", "ALIAS", "MANUAL"] = "FUZZY"


class InvoiceRecognizeItem(BaseModel):
    row_index: int = Field(..., ge=1)
    drug_name: str | None = None
    specification: str | None = None
    manufacturer: str | None = None
    approval_number: str | None = None
    batch_number: str | None = None
    expire_date: date | None = None
    quantity: str | None = None
    unit_price: str | None = None
    amount: str | None = None
    confidence: float | None = None
    matched_drug_id: int | None = None
    drug_candidates: list[DrugCandidateInvoice] = []
    warnings: list[QualityWarning] = []


class InvoiceRecognizeResult(BaseModel):
    recognized_supplier_name: str | None = None
    supplier_candidates: list[SupplierCandidate] = []
    matched_supplier_id: int | None = None
    invoice_no: str | None = None
    invoice_date: date | None = None
    total_amount: str | None = None
    confidence: float | None = None
    items: list[InvoiceRecognizeItem] = []
    warnings: list[QualityWarning] = []


class InvoiceRecognizeResponse(BaseModel):
    request_id: str
    status: Literal["PENDING", "PROCESSING", "COMPLETED", "FAILED"]
    result: InvoiceRecognizeResult | None = None
    raw_response: Any | None = None
    error_code: str | None = None
    error_message: str | None = None
    started_at: datetime | None = None
    finished_at: datetime | None = None
    duration_ms: int | None = None


class InvoiceQualityCheckRequest(BaseModel):
    result: InvoiceRecognizeResult
    strict: bool = False
