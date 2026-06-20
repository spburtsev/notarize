.PHONY: generate gen-api

generate: gen-api

gen-api:
	cd backend && go tool ogen --target internal/oas --package oas --clean ../openapi.yaml
