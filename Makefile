.PHONY: run test

run:
	go run ./cmd/ingestion-gateway

test:
	go test ./...