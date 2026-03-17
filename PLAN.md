# Full-Stack Calculator — Implementation Plan

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Technology Decisions](#technology-decisions)
- [SOLID Principles Application](#solid-principles-application)
- [Monorepo Structure](#monorepo-structure)
- [Session 1: Go Backend — Calculation Engine](#session-1-go-backend--calculation-engine)
- [Session 2: Go Backend — HTTP API Layer](#session-2-go-backend--http-api-layer)
- [Session 3: Frontend — Scaffolding & Component Architecture](#session-3-frontend--scaffolding--component-architecture)
- [Session 4: Frontend — Calculator Logic & UI Integration](#session-4-frontend--calculator-logic--ui-integration)
- [Session 5: Styling, Responsiveness & Polish](#session-5-styling-responsiveness--polish)
- [Session 6: Docker, CI, README & Final QA](#session-6-docker-ci-readme--final-qa)
- [Final Manual QA Checklist](#final-manual-qa-checklist)
- [Progress Tracker](#progress-tracker)

---

## Overview

A full-stack calculator application with a **React (TypeScript) frontend** and a **Go backend** microservice. The frontend is a thin UI layer that delegates all computation to the backend REST API. The backend owns all arithmetic logic including operator precedence evaluation.

**Key characteristics:**

- Monorepo with `backend/` and `frontend/` directories
- Backend: Go 1.26+, stdlib `net/http` router, zero external dependencies for core logic
- Frontend: Vite + React + TypeScript, Tailwind CSS v4, CVA (class-variance-authority)
- Frosted glassmorphism UI theme (single theme, no dark/light toggle)
- Mathematical precedence respected (backend evaluates using shunting-yard algorithm)
- Stateless: no database, no calculation history
- Target test coverage: ~95% backend, ~90% frontend

---

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

### API Endpoints

| Method | Path                        | Description                          | Request Body Example                                                        | Response Example                          |
|--------|-----------------------------|--------------------------------------|-----------------------------------------------------------------------------|------------------------------------------|
| POST   | `/api/calculate`            | Single arithmetic operation          | `{"operand_a": 10, "operator": "+", "operand_b": 5}`                       | `{"result": 15}`                         |
| POST   | `/api/calculate/expression` | Full expression with precedence      | `{"expression": "5 + 3 * 2"}`                                              | `{"result": 11, "expression": "5 + 3 * 2"}` |
| GET    | `/api/operations`           | List supported operations            | —                                                                           | `{"operations": ["+", "-", "*", "/", "^", "sqrt", "%"]}` |

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

---

## Technology Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Monorepo | Single repo | Simplifies review, CI, and local dev for evaluators |
| Go HTTP router | `net/http` (Go 1.22+ enhanced routing) | Zero deps for routing; stdlib patterns like `POST /api/{path}` cover our needs |
| Expression evaluation | Shunting-yard algorithm | Well-understood, compact (~50-80 lines), respects precedence, demonstrates CS fundamentals |
| Frontend build tool | Vite | Fast HMR, native TypeScript/ESM support, industry standard |
| Styling | Tailwind CSS v4 + CVA | Utility-first CSS for rapid styling; CVA for type-safe component variants without CSS-in-JS runtime |
| Testing (Go) | Standard `testing` package | No unnecessary deps; table-driven tests are idiomatic Go |
| Testing (Frontend) | Vitest + React Testing Library | Vitest is native to Vite (shared config, same transform), same API as Jest, faster execution |
| E2E testing | Not included | Scope too small to justify Playwright/Cypress overhead; documented as future enhancement |
| CORS | Backend middleware | Production-correct pattern; Vite dev proxy also configured for convenience |
| State management | `useReducer` hook | Calculator state machine fits reducer pattern perfectly; no external lib needed |

---

## SOLID Principles Application

Applied where they naturally fit — not forced for the sake of completeness.

| Principle | Where Applied | How | Why It Fits |
|---|---|---|---|
| **SRP** — Single Responsibility | Each operator is one function. Handlers only route. Engine only calculates. Validator only validates. | No class/struct does more than one job | Natural separation — these are genuinely distinct concerns with different reasons to change |
| **OCP** — Open/Closed | Operator registry (map of string → Operator). Adding a new operator = register one new struct, zero changes to existing code | New behavior via extension, not modification | Textbook use case — a calculator with an extensible set of operations |
| **DIP** — Dependency Inversion | HTTP handlers depend on `Calculator` interface, not concrete struct. Frontend components depend on API client abstraction (injectable for testing) | High-level modules depend on abstractions | Genuinely useful — enables unit testing handlers and components with mocks |

**Not applied (would be forced):**
- **ISP**: The interfaces in this project are already small (3 methods max). There's no bloated interface to segregate.
- **LSP**: Naturally satisfied as a side effect of OCP (all operators implement the same interface), but not a driving design decision for this scope.

---

## Monorepo Structure

```
fullstack-calculator/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # Entry point, dependency injection, server start
│   ├── internal/
│   │   ├── calculator/
│   │   │   ├── calculator.go         # Calculator interface + concrete implementation
│   │   │   ├── calculator_test.go    # Engine unit tests
│   │   │   ├── operator.go           # Operator interface + concrete operators
│   │   │   ├── operator_test.go      # Operator unit tests
│   │   │   ├── expression.go         # Shunting-yard expression parser/evaluator
│   │   │   └── expression_test.go    # Expression evaluation tests
│   │   ├── handler/
│   │   │   ├── handler.go            # HTTP handler struct + routes
│   │   │   ├── handler_test.go       # Handler unit tests (httptest)
│   │   │   ├── request.go            # Request DTOs
│   │   │   └── response.go           # Response DTOs + error formatting
│   │   ├── middleware/
│   │   │   ├── cors.go               # CORS middleware
│   │   │   ├── logging.go            # Request logging middleware
│   │   │   ├── recovery.go           # Panic recovery middleware
│   │   │   └── content_type.go       # JSON content-type enforcement
│   │   └── validator/
│   │       ├── validator.go          # Input validation logic
│   │       └── validator_test.go     # Validator tests
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── api/
│   │   │   ├── calculator.ts         # API client interface + implementation
│   │   │   └── calculator.test.ts    # API client tests
│   │   ├── components/
│   │   │   ├── Calculator/
│   │   │   │   ├── Calculator.tsx     # Main calculator container
│   │   │   │   └── Calculator.test.tsx
│   │   │   ├── Display/
│   │   │   │   ├── Display.tsx        # Expression + result display
│   │   │   │   └── Display.test.tsx
│   │   │   ├── ButtonGrid/
│   │   │   │   ├── ButtonGrid.tsx     # Grid layout of buttons
│   │   │   │   └── ButtonGrid.test.tsx
│   │   │   ├── Button/
│   │   │   │   ├── Button.tsx         # Individual button with CVA variants
│   │   │   │   └── Button.test.tsx
│   │   │   └── ErrorMessage/
│   │   │       ├── ErrorMessage.tsx   # Error display component
│   │   │       └── ErrorMessage.test.tsx
│   │   ├── hooks/
│   │   │   ├── useCalculator.ts       # Calculator state machine (useReducer)
│   │   │   └── useCalculator.test.ts  # Hook tests
│   │   ├── types/
│   │   │   └── index.ts              # Shared TypeScript types
│   │   ├── App.tsx
│   │   ├── App.test.tsx
│   │   ├── main.tsx
│   │   └── index.css                 # Tailwind directives + glassmorphism utilities
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   └── package.json
├── .github/
│   └── workflows/
│       └── ci.yml                    # GitHub Actions: lint, test, build both layers
├── docker-compose.yml
├── Makefile                          # Orchestration: run, test, build, docker
├── .gitignore
├── PLAN.md                           # This file
└── README.md                         # Setup, API docs, design decisions
```

---

## Session 1: Go Backend — Calculation Engine

**Goal**: Implement and fully test the core calculation engine — operators, calculator, and expression evaluator. No HTTP layer yet.

**Commit message**: `feat(backend): add calculation engine with operator registry and expression evaluator`

### Steps

- [ ] **1.1** Initialize Go module
  - `cd backend && go mod init github.com/epolancog/fullstack-calculator/backend`
  - Create directory structure: `cmd/server/`, `internal/calculator/`, `internal/handler/`, `internal/middleware/`, `internal/validator/`

- [ ] **1.2** Implement `Operator` interface and concrete operators (`internal/calculator/operator.go`)
  - `Operator` interface: `Execute(a, b float64) (float64, error)`
  - Concrete operators: `Add`, `Subtract`, `Multiply`, `Divide`, `Power`, `SquareRoot`, `Percentage`
  - `Divide` returns error on division by zero
  - `SquareRoot` returns error on negative input (uses only operand_a, ignores operand_b)
  - `Percentage` is unary — uses only operand_a, returns `a / 100` (e.g., `50` → `0.5`)
  - Each operator is a struct satisfying the interface (strategy pattern)

- [ ] **1.3** Write operator unit tests (`internal/calculator/operator_test.go`)
  - Table-driven tests for each operator
  - Test cases:
    - Addition: positive, negative, zero, decimals
    - Subtraction: positive, negative, zero, result-negative
    - Multiplication: positive, negative, zero, decimals
    - Division: normal, by zero (error), decimals
    - Power: positive exponent, zero exponent, negative exponent
    - SquareRoot: positive, zero, negative (error)
    - Percentage: `Percentage(50, _)` → `0.5`, `Percentage(100, _)` → `1.0`, `Percentage(0, _)` → `0`
  - Verify error messages are descriptive

- [ ] **1.4** Implement `Calculator` interface and concrete implementation (`internal/calculator/calculator.go`)
  - `Calculator` interface:
    ```go
    type Calculator interface {
        Calculate(operandA float64, operator string, operandB float64) (float64, error)
        EvaluateExpression(expression string) (float64, error)
        SupportedOperations() []string
    }
    ```
  - Concrete `Calc` struct with operator registry: `map[string]Operator`
  - `NewCalculator()` constructor that registers all operators
  - `Calculate()` method: lookup operator, execute
  - `SupportedOperations()` method: return sorted list of operator keys

- [ ] **1.5** Implement expression parser/evaluator (`internal/calculator/expression.go`)
  - Shunting-yard algorithm for parsing infix expressions
  - Supported tokens: numbers (int, float, negative), operators (`+`, `-`, `*`, `/`, `^`, `sqrt`, `%`), parentheses (`(`, `)`)
  - Operator precedence:
    - `+`, `-`: precedence 1
    - `*`, `/`: precedence 2
    - `^`: precedence 3, right-associative
    - `sqrt`: precedence 4 (unary, prefix)
    - `%`: precedence 4 (unary, postfix — binds tightly to immediate left number)
  - **Percentage behavior (`%`)**: context-dependent based on the preceding binary operator:
    - With `*` or `/`: simple divide-by-100. E.g., `200 * 10%` → `200 * 0.1` → `20`
    - With `+` or `-`: percent of the left operand. E.g., `50 + 10%` → `50 + (10% of 50)` → `55`
    - Standalone: just divide by 100. E.g., `50%` → `0.5`
    - No left operand: implicit `0`. E.g., `% 50` → `0`
    - Chained: each `%` applies to its immediate left value. E.g., `50%%` → `0.005`
    - With `^`: `%` binds to the number, not the result. E.g., `2 ^ 3%` → `2 ^ 0.03` (Google-style)
    - Implementation: during RPN evaluation, when processing `%`, peek at the pending binary operator to determine behavior (~20-25 extra lines in evaluator)
  - Parentheses: `(` and `)` handled natively by shunting-yard algorithm for grouping sub-expressions
  - Tokenizer: split expression string into number and operator tokens
  - **Unary minus (negative numbers)**: during tokenization, if `-` appears at the start of the expression, immediately after another operator, or immediately after `(`, treat it as part of the next number (unary negation), not as the binary subtraction operator. This is ~5-10 extra lines in the tokenizer.
  - **`sqrt` syntax**: space-required, e.g., `sqrt 16`. `sqrt16` (no space) is rejected as an invalid token. `sqrt` has precedence 4 so it binds only to the immediate next value; use parentheses for complex operands: `sqrt (16 + 9)`
  - Evaluator: convert to postfix (RPN) via shunting-yard, then evaluate the RPN stack
  - Error handling: malformed expressions, mismatched operators/operands, division by zero during evaluation
  - `EvaluateExpression()` wired into the `Calculator` implementation

- [ ] **1.6** Write calculator unit tests (`internal/calculator/calculator_test.go`)
  - `Calculate()` tests:
    - Valid operations for each operator
    - Unknown operator → error
  - `SupportedOperations()` test: returns expected list

- [ ] **1.7** Write expression evaluator tests (`internal/calculator/expression_test.go`)
  - Precedence tests:
    - `"5 + 3 * 2"` → `11` (not 16)
    - `"10 - 2 * 3 + 4"` → `8`
    - `"2 ^ 3 + 1"` → `9`
    - `"10 / 2 + 3 * 4"` → `17`
  - Basic operation tests:
    - `"1 + 1"` → `2`
    - `"10 - 3"` → `7`
    - `"4 * 5"` → `20`
    - `"20 / 4"` → `5`
  - Edge cases:
    - Single number: `"42"` → `42`
    - Decimal numbers: `"1.5 + 2.5"` → `4`
    - Division by zero within expression → error
    - Empty expression → error
    - Invalid characters → error
    - Consecutive operators → error (e.g., `"5 + + 3"`)
  - Negative number (unary minus) tests:
    - `"-5 + 3"` → `-2` (negation at start of expression)
    - `"5 + -3"` → `2` (negation after operator)
    - `"5 * -2 + 1"` → `-9` (negation with precedence)
    - `"-5 * -2"` → `10` (double negation)
    - `"10 / -2"` → `-5` (negation with division)
  - Percentage (unary postfix, context-dependent) tests:
    - `"50%"` → `0.5` (standalone, divide by 100)
    - `"200 * 10%"` → `20` (with `*`, simple divide by 100)
    - `"50 + 10%"` → `55` (with `+`, 10% of 50 added)
    - `"50 - 20%"` → `40` (with `-`, 20% of 50 subtracted)
    - `"100 + 10% + 20%"` → `132` (chained: 100+10=110, then 110+22=132)
    - `"50%%"` → `0.005` (chained: 50%=0.5, then 0.5%=0.005)
    - `"(50 + 10)%"` → `0.6` (after parens, divide by 100)
    - `"2 * 3 + 50%"` → `9` (mid-precedence: 6 + 50% of 6 = 9)
  - Parentheses tests:
    - `"(5 + 3) * 2"` → `16`
    - `"((2 + 3)) * 4"` → `20`
    - `"10 * (2 + 3)"` → `50`
    - Mismatched parentheses → error (e.g., `"(5 + 3"`, `"5 + 3)"`)
  - Complex expressions:
    - `"2 + 3 * 4 - 6 / 2"` → `11`
    - `"2 ^ 3 * 2 + 1"` → `17`
    - `"(2 + 3) * (4 - 1)"` → `15`

- [ ] **1.8** Run tests and verify coverage
  - `go test ./internal/calculator/... -v -cover`
  - Target: 95%+ coverage on the calculator package
  - Fix any failing tests

- [ ] **1.9** Create root `.gitignore`
  - Go: binaries, `vendor/`, `*.exe`, `*.out`
  - Node: `node_modules/`, `dist/`, `.env`
  - IDE: `.vscode/`, `.idea/`
  - OS: `.DS_Store`, `Thumbs.db`

- [ ] **1.10** Create root `Makefile` (backend targets only for now)
  - `make test-backend` — run Go tests with coverage
  - `make run-backend` — run the server (placeholder, will wire in Session 2)
  - `make coverage-backend` — generate coverage report

- [ ] **1.11** Commit all changes

### Manual Test Scenarios (Session 1)

Since there is no HTTP server yet, manual testing is done via `go test` output:

| # | Scenario | Command | Expected |
|---|----------|---------|----------|
| 1 | All operator tests pass | `cd backend && go test ./internal/calculator/... -v` | All PASS |
| 2 | Coverage meets target | `cd backend && go test ./internal/calculator/... -cover` | ≥ 95% coverage |
| 3 | Expression `5 + 3 * 2` evaluates to 11 | Verified in test output | `11` (not 16) |
| 4 | Division by zero returns error | Verified in test output | Descriptive error, no panic |
| 5 | Sqrt of negative returns error | Verified in test output | Descriptive error, no panic |
| 6 | Empty/invalid expression returns error | Verified in test output | Descriptive error, no panic |
| 7 | Expression `(5 + 3) * 2` evaluates to 16 | Verified in test output | `16` (not 16 regardless — parentheses override precedence) |
| 8 | Expression `200 * 10%` evaluates to 20 | Verified in test output | `20` (percentage with `*`) |
| 9 | Expression `50 + 10%` evaluates to 55 | Verified in test output | `55` (percentage with `+`, 10% of 50) |
| 10 | Mismatched parentheses returns error | Verified in test output | Descriptive error, no panic |

---

## Session 2: Go Backend — HTTP API Layer

**Goal**: Expose the calculation engine via REST endpoints with validation, middleware, and comprehensive HTTP tests.

**Commit message**: `feat(backend): add REST API with handlers, validation, and middleware`

### Steps

- [ ] **2.1** Implement request/response DTOs (`internal/handler/request.go`, `internal/handler/response.go`)
  - Request DTOs:
    ```go
    type CalculateRequest struct {
        OperandA float64 `json:"operand_a"`
        Operator string  `json:"operator"`
        OperandB float64 `json:"operand_b"`
    }
    type ExpressionRequest struct {
        Expression string `json:"expression"`
    }
    ```
  - Response DTOs:
    ```go
    type CalculateResponse struct {
        Result float64 `json:"result"`
    }
    type ExpressionResponse struct {
        Result     float64 `json:"result"`
        Expression string  `json:"expression"`
    }
    type OperationsResponse struct {
        Operations []string `json:"operations"`
    }
    type ErrorResponse struct {
        Error ErrorDetail `json:"error"`
    }
    type ErrorDetail struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    }
    ```

- [ ] **2.2** Implement input validator (`internal/validator/validator.go`)
  - `ValidateCalculateRequest(req)` → error with code
  - `ValidateExpressionRequest(req)` → error with code
  - Validations:
    - Operator is in the supported list
    - Expression is not empty
    - Expression contains only valid characters (digits, operators, parentheses, spaces, decimal points)
    - Operands are finite numbers (not NaN, not Inf)
  - Return structured errors with error codes (`INVALID_OPERATOR`, `INVALID_OPERAND`, `INVALID_EXPRESSION`)

- [ ] **2.3** Write validator tests (`internal/validator/validator_test.go`)
  - Valid requests pass
  - Empty operator → `INVALID_OPERATOR`
  - Unsupported operator → `INVALID_OPERATOR`
  - NaN/Inf operand → `INVALID_OPERAND`
  - Empty expression → `INVALID_EXPRESSION`

- [ ] **2.4** Implement HTTP handlers (`internal/handler/handler.go`)
  - `Handler` struct with `Calculator` interface dependency (DIP)
  - `NewHandler(calc Calculator) *Handler`
  - `RegisterRoutes(mux *http.ServeMux)` method
  - Route registration using Go 1.22+ patterns:
    - `POST /api/calculate`
    - `POST /api/calculate/expression`
    - `GET /api/operations`
  - Each handler method:
    1. Decode JSON body
    2. Validate input
    3. Call calculator
    4. Encode JSON response
  - Proper HTTP status codes: 200 (success), 400 (validation/calc error), 405 (method not allowed), 500 (unexpected)

- [ ] **2.5** Implement middleware (`internal/middleware/`)
  - **CORS** (`cors.go`):
    - Allow origins: `*` (configurable)
    - Allow methods: `GET, POST, OPTIONS`
    - Allow headers: `Content-Type`
    - Handle preflight `OPTIONS` requests
  - **Logging** (`logging.go`):
    - Log method, path, status code, duration
    - Use `slog` (Go stdlib structured logging)
  - **Recovery** (`recovery.go`):
    - Catch panics, return 500 with generic error
    - Log the panic + stack trace
  - **Content-Type** (`content_type.go`):
    - For POST requests, enforce `Content-Type: application/json`
    - Return 415 (Unsupported Media Type) otherwise

- [ ] **2.6** Wire up `cmd/server/main.go`
  - Create calculator with all operators registered
  - Create handler with calculator injected
  - Create `http.ServeMux`, register routes
  - Wrap mux with middleware chain: Recovery → Logging → CORS → ContentType → mux
  - Read port from `PORT` env var, default to `8080`
  - Graceful shutdown on SIGINT/SIGTERM
  - Log server start/stop

- [ ] **2.7** Write handler unit tests (`internal/handler/handler_test.go`)
  - Use `httptest.NewRecorder()` for unit tests
  - Tests for `POST /api/calculate`:
    - Valid addition → 200, correct result
    - Valid division → 200, correct result
    - Division by zero → 400, `DIVISION_BY_ZERO` error code
    - Invalid operator → 400, `INVALID_OPERATOR` error code
    - Malformed JSON → 400, `MALFORMED_JSON` error code
    - Empty body → 400
    - Wrong content-type → 415
  - Tests for `POST /api/calculate/expression`:
    - Valid expression → 200, correct result + echo expression
    - Expression with precedence → 200, correct result
    - Empty expression → 400, `INVALID_EXPRESSION`
    - Invalid expression → 400
  - Tests for `GET /api/operations`:
    - Returns 200, list of operations
  - Test for unsupported method:
    - `GET /api/calculate` → 405

- [ ] **2.8** Write integration test (full HTTP round-trip)
  - Use `httptest.NewServer()` to spin up the full server with middleware
  - Test a real HTTP request through the full middleware chain
  - Verify CORS headers are present
  - Verify logging doesn't break the response

- [ ] **2.9** Update `Makefile`
  - `make run-backend` — now runs the actual server
  - `make test-backend` — runs all backend tests
  - `make coverage-backend` — generates coverage report with `-coverprofile`

- [ ] **2.10** Run full test suite and verify coverage
  - `go test ./... -v -cover`
  - Target: 95%+ across all packages

- [ ] **2.11** Commit all changes

### Manual Test Scenarios (Session 2)

Start the server with `make run-backend` (or `cd backend && go run ./cmd/server/`) and test with `curl`:

| # | Scenario | Command | Expected |
|---|----------|---------|----------|
| 1 | Server starts | `make run-backend` | Logs "server started on :8080" |
| 2 | Single addition | `curl -s -X POST http://localhost:8080/api/calculate -H "Content-Type: application/json" -d '{"operand_a": 5, "operator": "+", "operand_b": 3}'` | `{"result": 8}` |
| 3 | Division by zero | `curl -s -X POST http://localhost:8080/api/calculate -H "Content-Type: application/json" -d '{"operand_a": 10, "operator": "/", "operand_b": 0}'` | 400 with `DIVISION_BY_ZERO` |
| 4 | Expression with precedence | `curl -s -X POST http://localhost:8080/api/calculate/expression -H "Content-Type: application/json" -d '{"expression": "5 + 3 * 2"}'` | `{"result": 11, "expression": "5 + 3 * 2"}` |
| 5 | List operations | `curl -s http://localhost:8080/api/operations` | JSON array of supported operators |
| 6 | Invalid JSON | `curl -s -X POST http://localhost:8080/api/calculate -H "Content-Type: application/json" -d 'not json'` | 400 with `MALFORMED_JSON` |
| 7 | Wrong content-type | `curl -s -X POST http://localhost:8080/api/calculate -H "Content-Type: text/plain" -d '{}'` | 415 Unsupported Media Type |
| 8 | CORS preflight | `curl -s -X OPTIONS http://localhost:8080/api/calculate -H "Origin: http://localhost:3000" -H "Access-Control-Request-Method: POST" -v 2>&1 \| grep -i "access-control"` | CORS headers present |
| 9 | Invalid operator | `curl -s -X POST http://localhost:8080/api/calculate -H "Content-Type: application/json" -d '{"operand_a": 5, "operator": "invalid", "operand_b": 3}'` | 400 with `INVALID_OPERATOR` |
| 10 | Complex expression | `curl -s -X POST http://localhost:8080/api/calculate/expression -H "Content-Type: application/json" -d '{"expression": "2 + 3 * 4 - 6 / 2"}'` | `{"result": 11, ...}` |
| 11 | Percentage expression | `curl -s -X POST http://localhost:8080/api/calculate/expression -H "Content-Type: application/json" -d '{"expression": "200 * 10%"}'` | `{"result": 20, ...}` |
| 12 | Parentheses expression | `curl -s -X POST http://localhost:8080/api/calculate/expression -H "Content-Type: application/json" -d '{"expression": "(5 + 3) * 2"}'` | `{"result": 16, ...}` |

---

## Session 3: Frontend — Scaffolding & Component Architecture

**Goal**: Scaffold the Vite project, configure Tailwind + CVA, implement the API client and base components (Button, ErrorMessage), with tests.

**Commit message**: `feat(frontend): scaffold Vite project with Tailwind, CVA, API client, and base components`

### Steps

- [ ] **3.1** Scaffold Vite React TypeScript project
  - `npm create vite@latest frontend -- --template react-ts`
  - Verify the project runs: `cd frontend && npm install && npm run dev`

- [ ] **3.2** Install and configure Tailwind CSS v4
  - Install Tailwind and its Vite plugin
  - Configure `index.css` with Tailwind directives
  - Verify Tailwind classes render correctly

- [ ] **3.3** Install and configure CVA
  - `npm install class-variance-authority`
  - Install `clsx` and `tailwind-merge` for utility class merging
  - Create a `cn()` utility function (`src/lib/utils.ts`): `cn(...inputs) => twMerge(clsx(...inputs))`

- [ ] **3.4** Install and configure testing tools
  - `npm install -D vitest @testing-library/react @testing-library/jest-dom @testing-library/user-event jsdom`
  - Configure `vitest` in `vite.config.ts` (jsdom environment, setup file)
  - Create test setup file (`src/test/setup.ts`) importing `@testing-library/jest-dom`

- [ ] **3.5** Configure Vite dev proxy
  - In `vite.config.ts`, proxy `/api` to `http://localhost:8080`
  - This avoids CORS issues during development

- [ ] **3.6** Define shared TypeScript types (`src/types/index.ts`)
  - ```typescript
    // API types
    interface CalculateRequest {
      operand_a: number;
      operator: string;
      operand_b: number;
    }
    interface ExpressionRequest {
      expression: string;
    }
    interface CalculateResponse {
      result: number;
    }
    interface ExpressionResponse {
      result: number;
      expression: string;
    }
    interface OperationsResponse {
      operations: string[];
    }
    interface ApiError {
      error: {
        code: string;
        message: string;
      };
    }

    // Calculator UI types
    type ButtonVariant = "number" | "operator" | "action" | "equals";
    type ButtonSize = "default" | "wide";
    interface CalculatorButton {
      label: string;
      value: string;
      variant: ButtonVariant;
      size?: ButtonSize;
    }
    ```

- [ ] **3.7** Implement API client (`src/api/calculator.ts`)
  - `CalculatorApi` interface:
    ```typescript
    interface CalculatorApi {
      calculate(operandA: number, operator: string, operandB: number): Promise<CalculateResponse>;
      evaluateExpression(expression: string): Promise<ExpressionResponse>;
      getOperations(): Promise<OperationsResponse>;
    }
    ```
  - `HttpCalculatorApi` class implementing the interface
  - Uses `fetch` — no axios dependency
  - Handles error responses: parse error JSON, throw typed error
  - Base URL configurable (default: `/api` for Vite proxy)

- [ ] **3.8** Write API client tests (`src/api/calculator.test.ts`)
  - Mock `fetch` using `vi.fn()` (Vitest)
  - Test successful calculate call → returns result
  - Test successful expression call → returns result + expression
  - Test successful getOperations call → returns list
  - Test API error response → throws with error code and message
  - Test network error → throws appropriate error

- [ ] **3.9** Implement Button component (`src/components/Button/Button.tsx`)
  - CVA variants:
    - `variant`: number (neutral), operator (accent), action (secondary), equals (primary/highlight)
    - `size`: default (1 column), wide (2 columns)
  - Props: `label`, `onClick`, `variant`, `size`, `disabled`, `ariaLabel`
  - Basic glassmorphism styling (transparent bg, border, hover/active states)
  - Press animation: scale down on `:active`

- [ ] **3.10** Write Button tests (`src/components/Button/Button.test.tsx`)
  - Renders label text
  - Applies correct variant classes
  - Applies wide class for size="wide"
  - Calls onClick when clicked
  - Renders with aria-label when provided
  - Disabled button doesn't fire onClick

- [ ] **3.11** Implement ErrorMessage component (`src/components/ErrorMessage/ErrorMessage.tsx`)
  - Props: `message: string | null`
  - Renders nothing when message is null
  - Displays error with icon and text
  - Subtle red/warning styling compatible with glass theme

- [ ] **3.12** Write ErrorMessage tests (`src/components/ErrorMessage/ErrorMessage.test.tsx`)
  - Renders nothing when message is null
  - Renders error message text
  - Has appropriate role/aria attributes for accessibility

- [ ] **3.13** Update `Makefile`
  - `make install-frontend` — `npm install`
  - `make test-frontend` — run Vitest
  - `make dev-frontend` — run Vite dev server

- [ ] **3.14** Run tests and verify all pass
  - `cd frontend && npm test`

- [ ] **3.15** Commit all changes

### Manual Test Scenarios (Session 3)

| # | Scenario | How to Test | Expected |
|---|----------|-------------|----------|
| 1 | Vite dev server starts | `cd frontend && npm run dev` | Opens on http://localhost:5173, no errors |
| 2 | Tailwind classes work | Add a `className="bg-red-500 text-white p-4"` to App.tsx, check browser | Red background, white text, padding visible |
| 3 | Button renders in browser | Temporarily render `<Button label="7" variant="number" onClick={() => {}} />` in App.tsx | Button visible with correct styling |
| 4 | Button variants look different | Render one of each variant side by side | Visual difference between number, operator, action, equals |
| 5 | All frontend tests pass | `cd frontend && npm test` | All PASS |
| 6 | Vite proxy works | Start backend (`make run-backend`) + frontend (`make dev-frontend`), open browser console, call `fetch('/api/operations').then(r => r.json()).then(console.log)` | Returns operations list from backend |

---

## Session 4: Frontend — Calculator Logic & UI Integration

**Goal**: Implement the calculator state machine, Display, ButtonGrid, and full Calculator component. Fully interactive calculator calling the backend API.

**Commit message**: `feat(frontend): implement calculator state machine, display, button grid, and full UI integration`

### Steps

- [ ] **4.1** Implement calculator state machine (`src/hooks/useCalculator.ts`)
  - Uses `useReducer` for predictable state transitions
  - State shape:
    ```typescript
    interface CalculatorState {
      expression: string;      // Full expression being built (e.g., "5 + 3 * ")
      currentInput: string;    // Current number being typed (e.g., "2")
      result: string | null;   // Result after evaluation
      error: string | null;    // Error message
      isLoading: boolean;      // Awaiting API response
    }
    ```
  - Actions/transitions:
    - `DIGIT` — append digit to currentInput
    - `DECIMAL` — append decimal point (prevent double decimal)
    - `OPERATOR` — append currentInput + operator to expression, clear currentInput
    - `EQUALS` — send full expression to API, display result
    - `CLEAR` — reset all state
    - `BACKSPACE` — remove last character from currentInput
    - `CLEAR_ENTRY` — clear only currentInput
  - Behavior rules:
    - After `EQUALS` + result displayed: new digit starts fresh, new operator uses result as start
    - Prevent leading zeros (except `0.`)
    - Prevent consecutive operators (replace last operator)
    - `sqrt` is handled as a unary prefix operator
    - `%` is handled as a unary postfix operator — pressing `%` appends `%` to currentInput immediately (no second operand needed). E.g., typing `50` then `%` produces `50%` in the expression. Multiple `%` presses allowed (e.g., `50%%` → `0.005`).
    - Parentheses: `(` and `)` buttons append to expression. Track open paren count to validate when `)` is allowed.
    - **Negative numbers**: pressing `-` when there is no currentInput and expression is empty (or ends with an operator or `(`) starts a negative number (appends `-` to currentInput). The backend tokenizer handles the actual unary minus parsing.
  - Receives `CalculatorApi` as dependency (DIP — injectable for testing)

- [ ] **4.2** Write useCalculator tests (`src/hooks/useCalculator.test.ts`)
  - Use `@testing-library/react` `renderHook`
  - Mock API client
  - Test state transitions:
    - Type digits → currentInput updates
    - Type decimal → added once only
    - Select operator → expression builds
    - Press equals → API called, result displayed
    - Clear → resets all
    - Backspace → removes last digit
  - Test expression building:
    - `5 + 3 * 2 =` → API called with `"5 + 3 * 2"`
    - Result displayed after API response
  - Test error handling:
    - API returns error → error state set
    - Network failure → error state set
  - Test edge cases:
    - Pressing equals with no expression → no API call
    - Pressing operator with no input → no-op (or uses 0)
    - After result, new digit starts fresh expression
    - After result, operator continues with result
    - Prevent double decimal
    - Prevent leading zeros
  - Test negative number input:
    - Press `-` as first input → currentInput is `"-"`
    - Type `-`, `5`, `+`, `3`, `=` → expression sent as `"-5 + 3"`
    - Press `5`, `+`, `-`, `3`, `=` → expression sent as `"5 + -3"`

- [ ] **4.3** Implement Display component (`src/components/Display/Display.tsx`)
  - Props: `expression: string`, `currentInput: string`, `result: string | null`, `isLoading: boolean`
  - Two lines:
    - Top: expression history (smaller, muted)
    - Bottom: current input or result (larger, primary)
  - Shows loading indicator when `isLoading`
  - Number formatting: add commas for thousands, limit decimal places for display
  - Text overflow: shrink font or scroll for long numbers/expressions
  - Glassmorphism styling for the display panel

- [ ] **4.4** Write Display tests (`src/components/Display/Display.test.tsx`)
  - Renders current input
  - Renders expression
  - Renders result when present
  - Shows loading indicator
  - Handles long numbers (no overflow/breakage)

- [ ] **4.5** Implement ButtonGrid component (`src/components/ButtonGrid/ButtonGrid.tsx`)
  - Defines button layout as data (array of `CalculatorButton` objects)
  - Layout (4 columns, classic calculator):
    ```
    [ C  ] [ ⌫  ] [  (  ] [  )  ]
    [ √  ] [ ^  ] [  %  ] [  ÷  ]
    [ 7  ] [ 8  ] [  9  ] [  ×  ]
    [ 4  ] [ 5  ] [  6  ] [  -  ]
    [ 1  ] [ 2  ] [  3  ] [  +  ]
    [ 0 (wide) ] [  .  ] [ = (equals) ]
    ```
  - Note: `%` is a unary postfix button (appends to current number, no second operand)
  - Renders Button components with correct variants
  - Passes click callbacks

- [ ] **4.6** Write ButtonGrid tests (`src/components/ButtonGrid/ButtonGrid.test.tsx`)
  - Renders all number buttons (0-9)
  - Renders all operator buttons
  - Renders action buttons (C, ⌫, =)
  - Renders advanced operation buttons (√, ^, %)
  - Renders parentheses buttons ((, ))
  - Click on button fires correct callback with correct value

- [ ] **4.7** Implement Calculator container (`src/components/Calculator/Calculator.tsx`)
  - Composes: Display + ButtonGrid + ErrorMessage
  - Uses `useCalculator` hook
  - Passes state and callbacks down to children
  - Keyboard event listener:
    - `0-9`, `.` → digit/decimal
    - `+`, `-`, `*`, `/`, `^` → binary operator
    - `%` → postfix unary operator (appends to current number)
    - `(`, `)` → parentheses
    - `Enter`, `=` → equals
    - `Escape` → clear
    - `Backspace` → backspace

- [ ] **4.8** Write Calculator integration tests (`src/components/Calculator/Calculator.test.tsx`)
  - Mock API client
  - Test: click digits → display updates
  - Test: click operator → expression builds
  - Test: click equals → API called, result shown
  - Test: click clear → resets
  - Test: keyboard input works (simulate keydown events)
  - Test: error from API → error message displayed
  - Test: full calculation flow: `5 + 3 * 2 =` → shows result from API

- [ ] **4.9** Update App.tsx
  - Render Calculator component
  - Pass real API client instance
  - Basic layout: centered on page

- [ ] **4.10** Write App test (`src/App.test.tsx`)
  - Renders without crashing
  - Calculator component is present

- [ ] **4.11** Run all frontend tests
  - `npm test`
  - Verify all pass, check coverage

- [ ] **4.12** Commit all changes

### Manual Test Scenarios (Session 4)

Start both backend (`make run-backend`) and frontend (`make dev-frontend`):

| # | Scenario | Action | Expected |
|---|----------|--------|----------|
| 1 | Basic addition | Click: 5, +, 3, = | Display shows expression `5 + 3 =` and result `8` |
| 2 | Precedence | Click: 5, +, 3, ×, 2, = | Result is `11` (not 16) |
| 3 | Division by zero | Click: 1, 0, ÷, 0, = | Error message displayed (not a crash) |
| 4 | Clear | Click: 5, +, 3, C | Display resets to empty/0 |
| 5 | Backspace | Click: 1, 2, 3, ⌫ | Display shows `12` |
| 6 | Decimal input | Click: 1, ., 5, +, 2, ., 5, = | Result is `4` |
| 7 | Keyboard input | Type: `5`, `+`, `3`, `Enter` | Same as clicking, result `8` |
| 8 | Keyboard Escape | Type some digits, press `Escape` | Calculator clears |
| 9 | Consecutive operators | Click: 5, +, -, 3, = | Operator replaced: `5 - 3 = 2` |
| 10 | After result, new digit | Click: 5, +, 3, =, 7 | New expression starts with `7` |
| 11 | After result, new operator | Click: 5, +, 3, =, +, 2, = | Expression `8 + 2 = 10` (uses previous result) |
| 12 | Square root | Click: √, 9, = (or 9, √) | Result is `3` |
| 13 | Negative number at start | Click: -, 5, +, 3, = | Result is `-2` |
| 14 | Negative number after operator | Click: 5, +, -, 3, ×, 2, = | Expression `5 + -3 * 2`, result is `-1` |
| 15 | All frontend tests pass | `cd frontend && npm test` | All PASS |

---

## Session 5: Styling, Responsiveness & Polish

**Goal**: Premium frosted glassmorphism design, responsive layout, micro-interactions, accessibility. No logic changes — purely visual.

**Commit message**: `style(frontend): add frosted glassmorphism design, responsive layout, and micro-interactions`

### Steps

- [ ] **5.1** Design the glassmorphism foundation
  - Background: animated gradient mesh or static gradient (dark-ish, to make glass pop)
  - Define Tailwind custom utilities/theme extensions for glass effects:
    - `glass-panel`: semi-transparent bg + `backdrop-filter: blur(16px)` + subtle border
    - `glass-button`: lighter glass effect for buttons
    - `glass-display`: darker glass for the display area

- [ ] **5.2** Style the Calculator container
  - Centered card on desktop, full-width on mobile
  - Max-width: ~400px
  - Rounded corners (xl or 2xl)
  - Outer glass panel with shadow
  - Padding and spacing

- [ ] **5.3** Style the Display component
  - Glass panel inset within the calculator
  - Expression line: smaller text, muted/secondary color, right-aligned
  - Current input / result line: large text (2xl-4xl), right-aligned, white/bright
  - Text truncation or font-size scaling for long inputs
  - Subtle inner shadow for depth

- [ ] **5.4** Style the ButtonGrid
  - CSS Grid: 4 columns, 6 rows, consistent gap
  - Wide buttons span 2 columns
  - Consistent button height

- [ ] **5.5** Style the Button variants (CVA)
  - **Number buttons**: glass background, white text, subtle border
  - **Operator buttons**: slightly tinted glass (accent color — blue/purple), bolder text
  - **Action buttons** (C, ⌫): slightly different tint (muted/gray)
  - **Equals button**: strong accent color (gradient or solid), stands out
  - All buttons:
    - Rounded corners
    - Hover: brighten/glow
    - Active/press: scale(0.95) + darken
    - Transition: smooth (150ms)
    - Min height for touch targets (48px+)

- [ ] **5.6** Style ErrorMessage
  - Subtle red/amber tint glass
  - Icon + message text
  - Shake animation on appear (CSS keyframes)

- [ ] **5.7** Implement responsive design
  - Mobile (< 640px): calculator fills width with padding, larger buttons for touch
  - Tablet (640px - 1024px): centered card, moderate size
  - Desktop (> 1024px): centered card, max-width 400px
  - Media queries via Tailwind breakpoints

- [ ] **5.8** Add micro-interactions
  - Button press: `transform: scale(0.95)` on `:active`
  - Result display: fade-in animation when result appears
  - Error: shake animation
  - Loading: subtle pulse or spinner on display
  - Transitions: all color/transform changes animated (150ms ease)

- [ ] **5.9** Accessibility audit
  - All buttons have `aria-label` (especially symbols like `×`, `÷`)
  - Focus indicators visible on keyboard navigation (ring/outline)
  - Color contrast: verify text is readable against glass backgrounds
  - Screen reader: display has `aria-live="polite"` for result announcements

- [ ] **5.10** Cross-viewport visual QA
  - Test at 375px (mobile), 768px (tablet), 1440px (desktop)
  - Verify no overflow, no broken layouts
  - Verify glass effects render (fallback for browsers without backdrop-filter if needed)

- [ ] **5.11** Run all tests (ensure styling changes didn't break anything)
  - `npm test`
  - All existing tests should still pass

- [ ] **5.12** Commit all changes

### Manual Test Scenarios (Session 5)

| # | Scenario | How to Test | Expected |
|---|----------|-------------|----------|
| 1 | Desktop appearance | Open http://localhost:5173 at 1440px width | Centered glass calculator card, gradient background, premium look |
| 2 | Mobile appearance | Open browser DevTools, toggle device toolbar, select iPhone 14 (390px) | Calculator fills width, buttons are touch-friendly, no horizontal scroll |
| 3 | Tablet appearance | DevTools, iPad view (768px) | Centered, appropriately sized |
| 4 | Button hover effect | Hover over a number button | Subtle glow/brighten |
| 5 | Button press animation | Click and hold a button | Button scales down slightly |
| 6 | Result animation | Complete a calculation (5 + 3 =) | Result fades in smoothly |
| 7 | Error animation | Trigger error (divide by zero) | Error message appears with shake |
| 8 | Keyboard focus | Tab through buttons | Visible focus ring/outline on each button |
| 9 | Glass effect visible | Look at calculator and buttons | Semi-transparent panels with blur, background gradient visible through |
| 10 | Long expression display | Type a long expression like `123456 + 789012 * 345678` | Text doesn't overflow; shrinks or scrolls gracefully |
| 11 | All tests still pass | `cd frontend && npm test` | All PASS |

---

## Session 6: Docker, CI, README & Final QA

**Goal**: Dockerize both services, CI pipeline, comprehensive README, coverage reports, final QA pass.

**Commit message**: `docs: add Docker setup, CI pipeline, and comprehensive README`

### Steps

- [ ] **6.1** Create backend Dockerfile (`backend/Dockerfile`)
  - Multi-stage build:
    - Stage 1 (`builder`): `golang:1.26-alpine`, copy source, `go build -o server ./cmd/server/`
    - Stage 2 (`runtime`): `alpine:latest` (or `scratch` if no CGO), copy binary, expose port 8080
  - Small final image size

- [ ] **6.2** Create frontend Dockerfile (`frontend/Dockerfile`)
  - Multi-stage build:
    - Stage 1 (`builder`): `node:22-alpine`, copy source, `npm ci`, `npm run build`
    - Stage 2 (`runtime`): `nginx:alpine`, copy built files to nginx html dir
  - Nginx config: serve SPA (fallback to index.html) + proxy `/api` to backend

- [ ] **6.3** Create nginx config for frontend (`frontend/nginx.conf`)
  - Serve static files from `/usr/share/nginx/html`
  - Proxy `/api/` to `http://backend:8080/api/`
  - SPA fallback: `try_files $uri $uri/ /index.html`

- [ ] **6.4** Create `docker-compose.yml`
  - Services:
    - `backend`: build from `./backend`, port `8080:8080`
    - `frontend`: build from `./frontend`, port `3000:80`, depends on `backend`
  - Network: shared bridge network
  - Health check for backend

- [ ] **6.5** Test Docker setup
  - `docker-compose up --build`
  - Verify frontend accessible at http://localhost:3000
  - Verify API calls work through nginx proxy
  - Verify calculation works end-to-end

- [ ] **6.6** Create GitHub Actions CI workflow (`.github/workflows/ci.yml`)
  - Triggers: push to `master`, pull requests
  - Jobs:
    - **backend**:
      - Setup Go 1.26
      - `go vet ./...`
      - `go test ./... -v -coverprofile=coverage.out`
      - Upload coverage artifact
    - **frontend**:
      - Setup Node 22
      - `npm ci`
      - `npm run lint` (if ESLint configured)
      - `npm test -- --coverage`
      - `npm run build`
      - Upload coverage artifact

- [ ] **6.7** Configure ESLint for frontend (if not already)
  - Basic TypeScript + React rules
  - Add `lint` script to `package.json`

- [ ] **6.8** Generate coverage reports
  - Backend: `go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out`
  - Frontend: `npm test -- --coverage`
  - Verify both meet targets (95% backend, 90% frontend)

- [ ] **6.9** Write comprehensive README.md
  - **Project overview**: what it is, screenshot/demo
  - **Architecture**: ASCII diagram (same as this plan), data flow
  - **Tech stack**: Go, React, TypeScript, Tailwind, Vite
  - **Design decisions table**: each decision + rationale (from this plan)
  - **SOLID principles**: where and how applied
  - **Setup instructions**:
    - Prerequisites (Go 1.22+, Node 18+, npm)
    - Local development (start backend, start frontend)
    - Docker (`docker-compose up --build`)
  - **API documentation**:
    - Each endpoint with method, path, description
    - Request/response examples with `curl`
    - Error response format and codes
  - **Testing**:
    - How to run tests (backend, frontend, all)
    - Coverage report generation
  - **Project structure**: directory tree with descriptions
  - **Future enhancements**:
    - Parentheses support in expressions
    - Calculation history with persistence
    - Dark/light theme toggle
    - E2E tests with Playwright
    - Memory functions (M+, M-, MR, MC)
    - Scientific calculator mode

- [ ] **6.10** Final pass: verify all Makefile targets work
  - `make test-backend`
  - `make test-frontend`
  - `make run-backend`
  - `make dev-frontend`
  - `make coverage-backend`
  - Add `make docker-up` / `make docker-down` targets

- [ ] **6.11** Commit all changes

### Manual Test Scenarios (Session 6)

| # | Scenario | How to Test | Expected |
|---|----------|-------------|----------|
| 1 | Docker build succeeds | `docker-compose build` | Both images build without errors |
| 2 | Docker-compose up | `docker-compose up` | Both services start, logs visible |
| 3 | Frontend via Docker | Open http://localhost:3000 | Calculator UI loads |
| 4 | API via Docker nginx proxy | Complete a calculation in the Docker-served frontend | Result returns correctly |
| 5 | CI workflow syntax valid | `act` (local) or push branch and check GitHub Actions | Workflow runs or parses without errors |
| 6 | Backend coverage target | `make coverage-backend` | ≥ 95% reported |
| 7 | Frontend coverage target | `cd frontend && npm test -- --coverage` | ≥ 90% reported |
| 8 | README renders correctly | View README.md on GitHub (or preview locally) | All sections, tables, code blocks render properly |
| 9 | curl examples from README work | Copy each curl example from README, run it against running backend | All return expected responses |
| 10 | Makefile targets | Run each `make` target | All succeed without errors |

---

## Final Manual QA Checklist

After all 6 sessions are complete, run through this comprehensive end-to-end checklist:

### Setup & Build

| # | Scenario | Steps | Expected |
|---|----------|-------|----------|
| 1 | Fresh clone setup | `git clone`, follow README instructions | Project runs successfully |
| 2 | Backend starts | `make run-backend` | Server logs "started on :8080" |
| 3 | Frontend starts | `make dev-frontend` | Vite dev server on :5173 |
| 4 | Docker full stack | `docker-compose up --build` | Both services healthy, frontend on :3000 |

### Calculator — Basic Operations

| # | Scenario | Input | Expected Result |
|---|----------|-------|-----------------|
| 5 | Addition | `7 + 3 =` | `10` |
| 6 | Subtraction | `10 - 4 =` | `6` |
| 7 | Multiplication | `6 × 8 =` | `48` |
| 8 | Division | `15 ÷ 3 =` | `5` |
| 9 | Exponentiation | `2 ^ 10 =` | `1024` |
| 10 | Percentage (multiply) | `200 × 10 % =` | `20` (10% of 200) |
| 11 | Percentage (add) | `50 + 10 % =` | `55` (10% of 50, added) |
| 12 | Square root | `√ 144 =` | `12` |
| 13 | Parentheses | `( 5 + 3 ) × 2 =` | `16` |

### Calculator — Precedence

| # | Scenario | Input | Expected Result |
|---|----------|-------|-----------------|
| 14 | Multiply before add | `5 + 3 × 2 =` | `11` |
| 15 | Divide before subtract | `10 - 6 ÷ 2 =` | `7` |
| 16 | Mixed precedence | `2 + 3 × 4 - 6 ÷ 2 =` | `11` |
| 17 | Power precedence | `2 ^ 3 + 1 =` | `9` |
| 18 | Parentheses override precedence | `( 2 + 3 ) × 2 =` | `10` (not `8`) |

### Calculator — Edge Cases

| # | Scenario | Input | Expected |
|---|----------|-------|----------|
| 19 | Division by zero | `10 ÷ 0 =` | Error message displayed |
| 20 | Sqrt of negative | `√ -4 =` | Error message displayed |
| 21 | Decimal input | `1.5 + 2.5 =` | `4` |
| 22 | Large numbers | `999999 × 999999 =` | Correct result, no overflow in display |
| 23 | Single number equals | `42 =` | `42` |
| 24 | Double decimal prevented | `1..5` | Only one decimal: `1.5` |
| 25 | Leading zero prevented | `007` | Displays `7` (or `0.07` if decimal) |
| 26 | Negative number at start | `-5 + 3 =` | `-2` |
| 27 | Negative after operator | `5 + -3 =` | `2` |
| 28 | Negative with precedence | `5 * -2 + 1 =` | `-9` |
| 29 | Double negation | `-5 * -2 =` | `10` |
| 30 | Mismatched parentheses | `( 5 + 3 =` | Error message displayed |

### Calculator — UX

| # | Scenario | Input | Expected |
|---|----------|-------|----------|
| 31 | Clear button | Type digits, press C | Resets to initial state |
| 32 | Backspace | Type `123`, press ⌫ | Shows `12` |
| 33 | Keyboard digits | Type `5`, `+`, `3`, `Enter` on keyboard | Result `8` |
| 34 | Keyboard Escape | Type digits, press Escape | Clears calculator |
| 35 | Keyboard Backspace | Type digits, press Backspace | Removes last digit |
| 36 | Consecutive operators | `5 + - 3 =` | Treats as `5 - 3 = 2` (replaced operator) |
| 37 | Continue after result | `5 + 3 =` (shows 8), then `+ 2 =` | `10` (continues from result) |
| 38 | New expression after result | `5 + 3 =` (shows 8), then type `7` | Starts fresh with `7` |

### Responsive Design

| # | Scenario | How to Test | Expected |
|---|----------|-------------|----------|
| 39 | Mobile layout | DevTools → iPhone 14 (390px) | Full-width, touch-friendly buttons |
| 40 | Tablet layout | DevTools → iPad (768px) | Centered, properly sized |
| 41 | Desktop layout | Full browser window (1440px) | Centered card, max-width ~400px |
| 42 | No horizontal scroll | All viewports | No overflow, no scrollbar |

### Accessibility

| # | Scenario | How to Test | Expected |
|---|----------|-------------|----------|
| 43 | Keyboard navigation | Tab through all buttons | Visible focus indicator on each |
| 44 | Screen reader | Inspect aria attributes | Buttons have aria-labels, result has aria-live |

### API (Direct)

| # | Scenario | curl command | Expected |
|---|----------|-------------|----------|
| 45 | Single operation | `curl -X POST .../api/calculate -d '{"operand_a":5,"operator":"+","operand_b":3}'` | `{"result":8}` |
| 46 | Expression | `curl -X POST .../api/calculate/expression -d '{"expression":"5 + 3 * 2"}'` | `{"result":11,...}` |
| 47 | List operations | `curl .../api/operations` | JSON array |
| 48 | Invalid input | `curl -X POST .../api/calculate -d '{"operand_a":5,"operator":"invalid","operand_b":3}'` | 400 error |

### Tests & Coverage

| # | Scenario | Command | Expected |
|---|----------|---------|----------|
| 49 | Backend tests pass | `make test-backend` | All PASS |
| 50 | Frontend tests pass | `make test-frontend` | All PASS |
| 51 | Backend coverage | `make coverage-backend` | ≥ 95% |
| 52 | Frontend coverage | Coverage report | ≥ 90% |

---

## Progress Tracker

Update this section at the start/end of each session to track overall progress.

| Session | Status | Date Started | Date Completed | Notes |
|---------|--------|-------------|----------------|-------|
| 1 — Go Calculation Engine | Not Started | — | — | — |
| 2 — Go HTTP API Layer | Not Started | — | — | — |
| 3 — Frontend Scaffolding | Not Started | — | — | — |
| 4 — Calculator Logic & UI | Not Started | — | — | — |
| 5 — Styling & Polish | Not Started | — | — | — |
| 6 — Docker, CI, README | Not Started | — | — | — |
| Final QA | Not Started | — | — | — |

### How to Use This Plan Across Sessions

1. At the start of a new Claude session, share this file and say: **"Continue with Session N of the plan in PLAN.md"**
2. The agent will read this plan, check the progress tracker, and proceed with the next uncompleted session
3. At the end of each session, the agent will update the progress tracker and commit
4. Each session produces exactly **one commit**
5. Run the manual test scenarios at the end of each session before committing
