.PHONY: run dev test build lint clean

# Run Go backend
run:
	go run ./cmd/api

# Run Go backend + frontend dev server
dev:
	@echo "Starting backend and frontend..."
	@(go run ./cmd/api &) && cd frontend && npm run dev

# Run all tests
test:
	go test ./... -count=1
	cd frontend && npx vitest run

# Build everything
build:
	go build -o bin/focustrack ./cmd/api
	cd frontend && npm run build

# Lint
lint:
	go vet ./...
	cd frontend && npx vue-tsc --noEmit

# Clean build artifacts
clean:
	rm -rf bin/ frontend/dist/ data/
