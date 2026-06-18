from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.config import settings
from app.database import init_pool, close_pool
from app.models.common import err
from app.request_id import new_request_id
from app.routers import health, drugs, invoices, matching


@asynccontextmanager
async def lifespan(app: FastAPI):
    await init_pool()
    yield
    await close_pool()


app = FastAPI(
    title="PharmacyERP AI Service",
    version=settings.service_version,
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

PREFIX = "/ai/api/v1"

app.include_router(health.router, prefix=PREFIX)
app.include_router(drugs.router, prefix=PREFIX)
app.include_router(invoices.router, prefix=PREFIX)
app.include_router(matching.router, prefix=PREFIX)


@app.exception_handler(Exception)
async def global_exception_handler(request: Request, exc: Exception):
    request_id = new_request_id()
    return JSONResponse(
        status_code=500,
        content=err(500, "服务内部错误", request_id, "INTERNAL_ERROR", [str(exc)]),
    )
