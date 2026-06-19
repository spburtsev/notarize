.PHONY: generate gen-api gen-db migrate-new migrate-up migrate-down up down run

DB_URL ?= postgres://notarize:notarize@localhost:5432/notarize?sslmode=disable

generate: gen-api gen-db

gen-api:
	cd backend && go tool ogen --target internal/oas --package oas --clean ../openapi.yaml

gen-db:
	cd backend && go tool sqlc generate

migrate-new:
	cd backend && go tool goose -dir internal/db/migrations create $(name) sql

migrate-up:
	cd backend && go tool goose -dir internal/db/migrations postgres "$(DB_URL)" up

migrate-down:
	cd backend && go tool goose -dir internal/db/migrations postgres "$(DB_URL)" down
