from __future__ import annotations

import asyncio
from datetime import datetime, timezone
from typing import Any

import httpx
from fastapi import APIRouter, WebSocket, WebSocketDisconnect


router = APIRouter(tags=["conjunctions"])


@router.get("/conjunctions")
async def list_conjunctions() -> list[dict[str, Any]]:
    snapshot = await fetch_prediction_snapshot()
    return snapshot["conjunctions"]


@router.websocket("/conjunctions/stream")
async def conjunction_stream(websocket: WebSocket) -> None:
    await websocket.accept()
    try:
        while True:
            snapshot = await fetch_prediction_snapshot()
            await websocket.send_json(snapshot["conjunctions"])
            await asyncio.sleep(2.0)
    except WebSocketDisconnect:
        return


async def fetch_prediction_snapshot(compute_url: str | None = None) -> dict[str, Any]:
    target = compute_url or "http://127.0.0.1:7001"
    try:
        async with httpx.AsyncClient(timeout=3.0) as client:
            response = await client.post(f"{target}/predict", json={"horizon_minutes": 45})
            response.raise_for_status()
            return response.json()
    except Exception:
        return sample_snapshot()


def sample_snapshot() -> dict[str, Any]:
    generated_at = datetime.now(timezone.utc).isoformat()
    objects = [
        {"object_id": "25544", "name": "ISS (ZARYA)", "x_km": 6780.0, "y_km": 1200.0, "z_km": 820.0, "speed_km_s": 7.67},
        {"object_id": "29716", "name": "FENGYUN 1C DEB", "x_km": -3220.0, "y_km": 5300.0, "z_km": 4210.0, "speed_km_s": 7.13},
        {"object_id": "54237", "name": "COSMOS 1408 DEB", "x_km": -3600.0, "y_km": 5180.0, "z_km": 4180.0, "speed_km_s": 7.38},
    ]
    conjunctions = [
        {
            "primary_id": "29716",
            "secondary_id": "54237",
            "miss_distance_km": 28.4,
            "probability_collision": 0.00057,
            "tca": generated_at,
        }
    ]
    return {"generated_at": generated_at, "objects": objects, "conjunctions": conjunctions}
