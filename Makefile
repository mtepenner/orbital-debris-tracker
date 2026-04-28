.PHONY: proto validate up down

proto:
	@echo "Proto definitions are staged in proto/orbit_engine.proto for optional generation workflows."

validate:
	cd tle_ingestion && go test ./...
	cd compute_engine && go test ./...
	python -m compileall api_gateway/app
	cd frontend && npm install
	cd frontend && npm run build

up:
	docker compose up -d --build

down:
	docker compose down
