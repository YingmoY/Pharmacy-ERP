from pydantic import BaseModel, Field


class SupplierMatchRequest(BaseModel):
    supplier_name: str = Field(..., min_length=1, max_length=255)
    limit: int = Field(5, ge=1, le=10)


class SupplierCandidate(BaseModel):
    supplier_id: int
    supplier_code: str
    name: str
    license_no: str | None = None
    contact_name: str | None = None
    contact_phone: str | None = None
    confidence: float
    match_reason: str
