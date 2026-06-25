.PHONY: generate gen-api gen-client lint-worker

generate: gen-api gen-client

gen-api:
	cd backend && go tool ogen --target internal/oas --package oas --clean ../openapi.yaml

gen-client:
	cd frontend && pnpm gen:api

lint-worker:
	cd worker && uv run --only-group dev ruff format . && uv run --only-group dev ruff check --fix .
