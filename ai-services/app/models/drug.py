from datetime import date
from decimal import Decimal
from typing import Literal
from pydantic import BaseModel, Field


class DrugSearchFilters(BaseModel):
    only_available: bool = False
    is_prescription: bool | None = None
    is_medicare: bool | None = None
    manufacturer: str | None = Field(None, max_length=100)
    dosage_form: str | None = Field(None, max_length=50)
    storage_condition: str | None = Field(None, max_length=100)
    near_expire_only: bool = False


class DrugSearchRequest(BaseModel):
    query: str = Field(..., min_length=1, max_length=200)
    search_mode: Literal["KEYWORD", "FUZZY", "SEMANTIC", "HYBRID"] = "HYBRID"
    limit: int = Field(10, ge=1, le=50)
    offset: int = Field(0, ge=0)
    filters: DrugSearchFilters = Field(default_factory=DrugSearchFilters)
    context: dict | None = None


class DrugInventorySummary(BaseModel):
    available_qty: int = 0
    in_stock_qty: int = 0
    reserved_qty: int = 0
    pending_qty: int = 0
    abnormal_qty: int = 0
    near_expire_available_qty: int = 0
    nearest_expire_date: date | None = None


class DrugSearchItem(BaseModel):
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
    is_prescription: bool
    is_medicare: bool
    inventory: DrugInventorySummary | None = None
    score: float
    match_reason: str | None = None
    highlights: list[str] = []


class DrugSearchResponse(BaseModel):
    query: str
    normalized_query: str | None = None
    search_mode: str
    total: int
    items: list[DrugSearchItem]
    suggestions: list[str] = []


class DrugSearchFeedbackRequest(BaseModel):
    request_id: str | None = None
    query: str = Field(..., max_length=200)
    selected_drug_id: int | None = None
    feedback_type: Literal["CLICK", "SELECT", "NO_RESULT", "BAD_RESULT", "MANUAL_CORRECTION"]
    feedback_note: str | None = Field(None, max_length=500)


class DrugMatchInputItem(BaseModel):
    row_index: int = Field(..., ge=1)
    drug_name: str = Field(..., max_length=200)
    specification: str | None = Field(None, max_length=100)
    manufacturer: str | None = Field(None, max_length=200)
    approval_number: str | None = Field(None, max_length=100)
    barcode: str | None = Field(None, max_length=100)


class DrugMatchRequest(BaseModel):
    items: list[DrugMatchInputItem] = Field(..., min_length=1, max_length=200)
    limit_per_item: int = Field(5, ge=1, le=10)


class DrugRecommendRequest(BaseModel):
    query: str = Field(..., min_length=1, max_length=500)
    limit: int = Field(10, ge=1, le=30)
    filters: DrugSearchFilters = Field(default_factory=DrugSearchFilters)


class DrugRecommendResponse(BaseModel):
    query: str
    explanation: str
    terms: list[str]
    total: int
    items: list[DrugSearchItem]
