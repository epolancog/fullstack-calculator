.PHONY: test-backend run-backend coverage-backend install-frontend test-frontend dev-frontend

test-backend:
	cd backend && go test ./... -v -count=1

run-backend:
	cd backend && go run ./cmd/server/

coverage-backend:
	cd backend && go test ./... -coverprofile=coverage.out -count=1
	cd backend && go tool cover -func=coverage.out

install-frontend:
	cd frontend && npm install

test-frontend:
	cd frontend && npm test

dev-frontend:
	cd frontend && npm run dev
