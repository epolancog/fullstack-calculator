package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/epolancog/fullstack-calculator/backend/internal/calculator"
	"github.com/epolancog/fullstack-calculator/backend/internal/validator"
)

func newTestHandler() *Handler {
	calc := calculator.NewCalculator()
	v := validator.NewValidator(calc.SupportedOperations())
	return NewHandler(calc, v)
}

func newTestMux() *http.ServeMux {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

// --- POST /api/calculate ---

func TestCalculateValidAddition(t *testing.T) {
	mux := newTestMux()
	body := `{"operand_a": 5, "operator": "+", "operand_b": 3}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp CalculateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Result != 8 {
		t.Errorf("expected result 8, got %v", resp.Result)
	}
}

func TestCalculateValidDivision(t *testing.T) {
	mux := newTestMux()
	body := `{"operand_a": 15, "operator": "/", "operand_b": 3}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp CalculateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Result != 5 {
		t.Errorf("expected result 5, got %v", resp.Result)
	}
}

func TestCalculateDivisionByZero(t *testing.T) {
	mux := newTestMux()
	body := `{"operand_a": 10, "operator": "/", "operand_b": 0}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "DIVISION_BY_ZERO")
}

func TestCalculateInvalidOperator(t *testing.T) {
	mux := newTestMux()
	body := `{"operand_a": 5, "operator": "invalid", "operand_b": 3}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "INVALID_OPERATOR")
}

func TestCalculateMalformedJSON(t *testing.T) {
	mux := newTestMux()
	body := `not json`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "MALFORMED_JSON")
}

func TestCalculateEmptyBody(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

// --- POST /api/calculate/expression ---

func TestExpressionValid(t *testing.T) {
	mux := newTestMux()
	body := `{"expression": "5 + 3 * 2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate/expression", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp ExpressionResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Result != 11 {
		t.Errorf("expected result 11, got %v", resp.Result)
	}
	if resp.Expression != "5 + 3 * 2" {
		t.Errorf("expected expression echoed, got %q", resp.Expression)
	}
}

func TestExpressionPrecedence(t *testing.T) {
	mux := newTestMux()
	body := `{"expression": "2 + 3 * 4 - 6 / 2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate/expression", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp ExpressionResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Result != 11 {
		t.Errorf("expected result 11, got %v", resp.Result)
	}
}

func TestExpressionEmpty(t *testing.T) {
	mux := newTestMux()
	body := `{"expression": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate/expression", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "INVALID_EXPRESSION")
}

func TestExpressionInvalid(t *testing.T) {
	mux := newTestMux()
	body := `{"expression": "5 & 3"}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate/expression", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "INVALID_EXPRESSION")
}

// --- GET /api/operations ---

func TestOperations(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/operations", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp OperationsResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Operations) == 0 {
		t.Error("expected non-empty operations list")
	}
}

// --- Additional error code coverage ---

func TestCalculateSqrtNegative(t *testing.T) {
	mux := newTestMux()
	body := `{"operand_a": -4, "operator": "sqrt", "operand_b": 0}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "SQRT_NEGATIVE")
}

func TestExpressionDivisionByZero(t *testing.T) {
	mux := newTestMux()
	body := `{"expression": "10 / 0"}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate/expression", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "DIVISION_BY_ZERO")
}

func TestExpressionSqrtNegative(t *testing.T) {
	mux := newTestMux()
	body := `{"expression": "sqrt -4"}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate/expression", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "SQRT_NEGATIVE")
}

func TestExpressionMalformedJSON(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodPost, "/api/calculate/expression", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	assertErrorCode(t, w, "MALFORMED_JSON")
}

// --- Method not allowed ---

func TestCalculateMethodNotAllowed(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/calculate", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d: %s", w.Code, w.Body.String())
	}
}

// helper

func assertErrorCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode string) {
	t.Helper()
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp.Error.Code != expectedCode {
		t.Errorf("expected error code %q, got %q (message: %s)", expectedCode, resp.Error.Code, resp.Error.Message)
	}
}
