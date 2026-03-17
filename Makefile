.PHONY: test-backend run-backend coverage-backend

test-backend:
	cd backend && go test ./... -v -count=1

run-backend:
	cd backend && go run ./cmd/server/

coverage-backend:
	cd backend && go test ./... -coverprofile=coverage.out -count=1
	cd backend && go tool cover -func=coverage.out
