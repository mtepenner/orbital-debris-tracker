from __future__ import annotations

import json
from dataclasses import dataclass, field
from typing import Any

try:
    from redis.asyncio import Redis
except Exception:  # pragma: no cover - fallback path for limited environments
    Redis = None


@dataclass(slots=True)
class OrbitCache:
    redis_url: str | None = None
    _memory: dict[str, dict[str, Any]] = field(init=False, repr=False)
    _redis: Any = field(init=False, repr=False, default=None)

    def __post_init__(self) -> None:
        self._memory = {}
        self._redis = Redis.from_url(self.redis_url, decode_responses=True) if self.redis_url and Redis else None

    async def seed_positions(self, objects: list[dict[str, Any]]) -> None:
        for item in objects:
            self._memory[item["object_id"]] = item
        if self._redis is not None:
            payload = {item["object_id"]: json.dumps(item) for item in objects}
            if payload:
                await self._redis.hset("orbital:positions", mapping=payload)

    async def list_positions(self) -> list[dict[str, Any]]:
        if self._redis is not None:
            values = await self._redis.hvals("orbital:positions")
            if values:
                return [json.loads(value) for value in values]
        return sorted(self._memory.values(), key=lambda item: item["object_id"])

    async def get_position(self, object_id: str) -> dict[str, Any] | None:
        if self._redis is not None:
            value = await self._redis.hget("orbital:positions", object_id)
            if value is not None:
                return json.loads(value)
        return self._memory.get(object_id)
