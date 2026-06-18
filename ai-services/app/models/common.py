from typing import Any
from pydantic import BaseModel


class ApiResponseBase(BaseModel):
    code: int
    message: str
    request_id: str


class ErrorDetail(BaseModel):
    error_code: str
    details: list[str] = []


class ErrorResponse(ApiResponseBase):
    data: Any = None
    error: ErrorDetail | None = None


class EmptyApiResponse(ApiResponseBase):
    data: dict | None = {}


def ok(data: Any, request_id: str, message: str = "success") -> dict:
    return {"code": 200, "message": message, "request_id": request_id, "data": data}


def err(code: int, message: str, request_id: str, error_code: str = "INTERNAL_ERROR", details: list[str] | None = None) -> dict:
    return {
        "code": code,
        "message": message,
        "request_id": request_id,
        "data": None,
        "error": {"error_code": error_code, "details": details or []},
    }
