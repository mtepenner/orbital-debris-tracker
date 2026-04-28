from __future__ import annotations

from fastapi import APIRouter, HTTPException, Request

from app.api.routes.conjunctions import fetch_prediction_snapshot


router = APIRouter(prefix="/catalog", tags=["catalog"])


@router.get("")
async def list_catalog(request: Request) -> list[dict[str, object]]:
    cache = request.app.state.cache
    objects = await cache.list_positions()
    if not objects:
        snapshot = await fetch_prediction_snapshot(request.app.state.compute_url)
        await cache.seed_positions(snapshot["objects"])
        objects = await cache.list_positions()
    return objects


@router.get("/{object_id}")
async def get_catalog_object(object_id: str, request: Request) -> dict[str, object]:
    cache = request.app.state.cache
    item = await cache.get_position(object_id)
    if item is None:
        snapshot = await fetch_prediction_snapshot(request.app.state.compute_url)
        await cache.seed_positions(snapshot["objects"])
        item = await cache.get_position(object_id)
    if item is None:
        raise HTTPException(status_code=404, detail="catalog object not found")
    return item
