package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/epolancog/fullstack-calculator/backend/internal/calculator"
	"github.com/epolancog/fullstack-calculator/backend/internal/middleware"
	"github.com/epolancog/fullstack-calculator/backend/internal/validator"
)

func newIntegrationServer() *httptest.Server {
	calc := calculator.NewCalculator()
	v := validator.NewValidator(calc.SupportedOperations())
	h := NewHandler(calc, v)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Full middleware chain
	var chain http.Handler = mux
	chain = middleware.ContentType(chain)
	chain = middleware.CORS(chain)
	chain = middleware.Logging(chain)
	chain = middleware.Recovery(chain)

	return httptest.NewServer(chain)
}

func TestIntegrationCalculate(t *testing.T) {
	server := newIntegrationServer()
	defer server.Close()

	body := `{"operand_a": 10, "operator": "+", "operand_b": 5}`
	resp, err := http.Post(server.URL+"/api/calculate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result CalculateResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Result != 15 {
		t.Errorf("expected 15, got %v", result.Result)
	}
}

func TestIntegrationExpression(t *testing.T) {
	server := newIntegrationServer()
	defer server.Close()

	body := `{"expression": "(5 + 3) * 2"}`
	resp, err := http.Post(server.URL+"/api/calculate/expression", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result ExpressionResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Result != 16 {
		t.Errorf("expected 16, got %v", result.Result)
	}
}

func TestIntegrationCORSHeaders(t *testing.T) {
	server := newIntegrationServer()
	defer server.Close()

	req, _ := http.NewRequest(http.MethodOptions, server.URL+"/api/calculate", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("missing or incorrect Access-Control-Allow-Origin header: %q", resp.Header.Get("Access-Control-Allow-Origin"))
	}
	if resp.Header.Get("Access-Control-Allow-Methods") == "" {
		t.Error("missing Access-Control-Allow-Methods header")
	}
}

func TestIntegrationContentTypeEnforcement(t *testing.T) {
	server := newIntegrationServer()
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/calculate", "text/plain", strings.NewReader("{}"))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("expected 415, got %d", resp.StatusCode)
	}
}

func TestIntegrationOperations(t *testing.T) {
	server := newIntegrationServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/operations")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result OperationsResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if len(result.Operations) == 0 {
		t.Error("expected non-empty operations list")
	}
}
