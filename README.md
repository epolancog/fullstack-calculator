# Fullstack Calculator

A full-stack calculator application with a **React (TypeScript) frontend** and a **Go backend** microservice. The frontend is a thin UI layer with a frosted glassmorphism design that delegates all computation to the backend REST API. The backend owns all arithmetic logic including operator precedence evaluation via the shunting-yard algorithm.

## Architecture

```
┌─────────────────────────────────────────────┐
│                  Frontend                    │
│              (React + TypeScript)            │
│                                             │
│  ┌─────────┐  ┌───────────┐  ┌───────────┐ │
│  │ Display  │  │ButtonGrid │  │  Error    │ │
│  │Component │  │ Component │  │ Message   │ │
│  └────┬─────┘  └─────┬─────┘  └───────────┘ │
│       │              │                       │
│  ┌────┴──────────────┴─────┐                │
│  │   useCalculator Hook     │                │
│  │   (state machine)        │                │
│  └────────────┬─────────────┘                │
│               │                              │
│  ┌────────────┴─────────────┐                │
│  │   API Client (interface) │                │
│  └────────────┬─────────────┘                │
└───────────────┼──────────────────────────────┘
                │ HTTP (JSON)
                │ POST /api/calculate
                │ POST /api/calculate/expression
                │ GET  /api/operations
┌───────────────┼──────────────────────────────┐
│               │         Backend              │
│  ┌────────────┴─────────────┐                │
│  │   HTTP Handlers          │                │
│  │   (routing + response)   │                │
│  └────────────┬─────────────┘                │
│               │                              │
│  ┌────────────┴─────────────┐                │
│  │   Validator              │                │
│  │   (input validation)     │                │
│  └────────────┬─────────────┘                │
│               │                              │
│  ┌────────────┴─────────────┐                │
│  │   Calculator Engine      │                │
│  │   (interface)            │                │
│  │   ┌───────────────────┐  │                │
│  │   │ Operator Registry  │  │                │
│  │   │ (strategy pattern) │  │                │
│  │   └───────────────────┘  │                │
│  │   ┌───────────────────┐  │                │
│  │   │ Expression Parser  │  │                │
│  │   │ (shunting-yard)   │  │                │
│  │   └───────────────────┘  │                │
│  └──────────────────────────┘                │
└──────────────────────────────────────────────┘
```

## Tech Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Backend | Go 1.26+, stdlib `net/http` | Zero external dependencies for routing; Go 1.22+ enhanced routing patterns |
| Frontend | Vite 6 + React 19 + TypeScript | Fast HMR, native TypeScript/ESM support |
| Styling | Tailwind CSS v4 + CVA | Utility-first CSS with type-safe component variants |
| Backend Testing | Go `testing` package | Table-driven tests, idiomatic Go |
| Frontend Testing | Vitest + React Testing Library | Native Vite integration, same API as Jest |
| Expression Eval | Shunting-yard algorithm | Compact, respects precedence, demonstrates CS fundamentals |
| State Management | `useReducer` hook | Calculator state machine fits reducer pattern perfectly |
| Containerization | Docker + docker-compose | Multi-stage builds for small images |

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Monorepo | Single repo with `backend/` and `frontend/` | Simplifies review, CI, and local dev |
| Expression evaluation | Shunting-yard algorithm | Well-understood, compact, respects precedence |
| Percentage behavior | Context-dependent | `50 + 10%` → `55` (percent of left operand with +/-), `200 * 10%` → `20` (simple divide-by-100 otherwise) |
| `sqrt` handling | Unary prefix operator | Pushed onto operator stack like a function, pops one value |
| Implicit multiplication | `5 √ 9` → `5 * sqrt 9` | Natural UX — typing a number then sqrt implies multiplication |
| Negative numbers | Unary minus in tokenizer | `-` at start or after operator/`(` treated as negation |
| API client abstraction | Interface + concrete class | Enables testing with mocks (DIP) |
| CORS | Backend middleware | Production-correct pattern; Vite dev proxy for convenience |

## SOLID Principles

| Principle | Where Applied | How |
|-----------|--------------|-----|
| **SRP** | Each operator is one function. Handlers only route. Engine only calculates. | Genuinely distinct concerns with different reasons to change |
| **OCP** | Operator registry (`map[string]Operator`) | Adding a new operator = register one new struct, zero changes to existing code |
| **DIP** | Handlers depend on `Calculator` interface. Frontend components depend on `CalculatorApi` interface | Enables unit testing with mocks |

## Setup

### Prerequisites

- Go 1.26+
- Node.js 22+
- npm

### Local Development

Start the backend:

```bash
make run-backend
```

In a separate terminal, start the frontend:

```bash
make install-frontend  # first time only
make dev-frontend
```

Open http://localhost:5173 — the Vite dev proxy forwards `/api` requests to the backend on port 8080.

### Docker

Requires [Docker Desktop](https://www.docker.com/products/docker-desktop/) running.

```bash
docker-compose up --build
```

- Frontend: http://localhost:3001
- Backend API: http://localhost:8080

## API Documentation

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/calculate` | Single arithmetic operation |
| POST | `/api/calculate/expression` | Full expression with precedence |
| GET | `/api/operations` | List supported operations |

### Examples

**Single operation:**

```bash
curl -s -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operand_a": 10, "operator": "+", "operand_b": 5}'
# {"result": 15}
```

**Expression with precedence:**

```bash
curl -s -X POST http://localhost:8080/api/calculate/expression \
  -H "Content-Type: application/json" \
  -d '{"expression": "5 + 3 * 2"}'
# {"result": 11, "expression": "5 + 3 * 2"}
```

**Percentage:**

```bash
curl -s -X POST http://localhost:8080/api/calculate/expression \
  -H "Content-Type: application/json" \
  -d '{"expression": "50 + 10%"}'
# {"result": 55, "expression": "50 + 10%"}
```

**List operations:**

```bash
curl -s http://localhost:8080/api/operations
# {"operations": ["%", "*", "+", "-", "/", "^", "sqrt"]}
```

### Error Response Format

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "division by zero is not allowed"
  }
}
```

Error codes: `DIVISION_BY_ZERO`, `INVALID_OPERATOR`, `INVALID_OPERAND`, `INVALID_EXPRESSION`, `MALFORMED_JSON`, `SQRT_NEGATIVE`.

## Testing

Run all backend tests:

```bash
make test-backend
```

Run all frontend tests:

```bash
make test-frontend
```

Generate coverage reports:

```bash
make coverage-backend
cd frontend && npx vitest run --coverage
```

## Project Structure

```
fullstack-calculator/
├── backend/
│   ├── cmd/server/
│   │   └── main.go                # Entry point, dependency injection
│   ├── internal/
│   │   ├── calculator/
│   │   │   ├── calculator.go       # Calculator interface + implementation
│   │   │   ├── operator.go         # Operator interface + concrete operators
│   │   │   └── expression.go       # Shunting-yard expression evaluator
│   │   ├── handler/
│   │   │   ├── handler.go          # HTTP handlers + routes
│   │   │   ├── request.go          # Request DTOs
│   │   │   └── response.go         # Response DTOs + error formatting
│   │   ├── middleware/             # CORS, logging, recovery, content-type
│   │   └── validator/             # Input validation
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── api/calculator.ts       # API client interface + implementation
│   │   ├── components/
│   │   │   ├── Button/             # Button with CVA variants
│   │   │   ├── ButtonGrid/         # Calculator button layout
│   │   │   ├── Calculator/         # Main container + keyboard listener
│   │   │   ├── Display/            # Expression + result display
│   │   │   └── ErrorMessage/       # Error display with shake animation
│   │   ├── hooks/useCalculator.ts  # Calculator state machine (useReducer)
│   │   ├── types/index.ts          # Shared TypeScript types
│   │   └── lib/utils.ts            # cn() utility for class merging
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
├── .github/workflows/ci.yml       # GitHub Actions CI pipeline
├── docker-compose.yml
├── Makefile
└── PLAN.md
```

## Future Enhancements

- Calculation history with persistence
- Dark/light theme toggle
- E2E tests with Playwright
- Memory functions (M+, M-, MR, MC)
- Scientific calculator mode
