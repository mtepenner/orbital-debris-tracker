from __future__ import annotations

import os

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.api.routes import catalog, conjunctions
from app.core.redis_client import OrbitCache


app = FastAPI(title="Orbital Debris Tracker API")
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

app.state.cache = OrbitCache(redis_url=os.getenv("REDIS_URL"))
app.state.compute_url = os.getenv("COMPUTE_ENGINE_URL", "http://127.0.0.1:7001")

app.include_router(catalog.router)
app.include_router(conjunctions.router)


@app.get("/health")
async def health() -> dict[str, str]:
    return {"status": "ok"}
