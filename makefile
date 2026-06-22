.PHONY: generate gen-api gen-client

generate: gen-api gen-client

gen-api:
	cd backend && go tool ogen --target internal/oas --package oas --clean ../openapi.yaml

gen-client:
	cd frontend && pnpm gen:api
