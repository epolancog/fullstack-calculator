import { describe, it, expect, vi, beforeEach } from "vitest";
import { HttpCalculatorApi, CalculatorApiError } from "./calculator";

const mockFetch = vi.fn();
vi.stubGlobal("fetch", mockFetch);

function jsonResponse(body: unknown, status = 200) {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  };
}

describe("HttpCalculatorApi", () => {
  let api: HttpCalculatorApi;

  beforeEach(() => {
    api = new HttpCalculatorApi("/api");
    mockFetch.mockReset();
  });

  describe("calculate", () => {
    it("returns result on success", async () => {
      mockFetch.mockResolvedValue(jsonResponse({ result: 8 }));

      const result = await api.calculate(5, "+", 3);

      expect(result).toEqual({ result: 8 });
      expect(mockFetch).toHaveBeenCalledWith("/api/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ operand_a: 5, operator: "+", operand_b: 3 }),
      });
    });

    it("throws CalculatorApiError on API error", async () => {
      mockFetch.mockResolvedValue(
        jsonResponse(
          { error: { code: "DIVISION_BY_ZERO", message: "division by zero is not allowed" } },
          400
        )
      );

      await expect(api.calculate(10, "/", 0)).rejects.toThrow(CalculatorApiError);
      await expect(api.calculate(10, "/", 0)).rejects.toMatchObject({
        code: "DIVISION_BY_ZERO",
        message: "division by zero is not allowed",
      });
    });
  });

  describe("evaluateExpression", () => {
    it("returns result and expression on success", async () => {
      mockFetch.mockResolvedValue(
        jsonResponse({ result: 11, expression: "5 + 3 * 2" })
      );

      const result = await api.evaluateExpression("5 + 3 * 2");

      expect(result).toEqual({ result: 11, expression: "5 + 3 * 2" });
      expect(mockFetch).toHaveBeenCalledWith("/api/calculate/expression", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ expression: "5 + 3 * 2" }),
      });
    });

    it("throws CalculatorApiError on invalid expression", async () => {
      mockFetch.mockResolvedValue(
        jsonResponse(
          { error: { code: "INVALID_EXPRESSION", message: "invalid expression" } },
          400
        )
      );

      await expect(api.evaluateExpression("")).rejects.toThrow(CalculatorApiError);
      await expect(api.evaluateExpression("")).rejects.toMatchObject({
        code: "INVALID_EXPRESSION",
      });
    });
  });

  describe("getOperations", () => {
    it("returns operations list on success", async () => {
      mockFetch.mockResolvedValue(
        jsonResponse({ operations: ["+", "-", "*", "/", "^", "sqrt", "%"] })
      );

      const result = await api.getOperations();

      expect(result).toEqual({
        operations: ["+", "-", "*", "/", "^", "sqrt", "%"],
      });
      expect(mockFetch).toHaveBeenCalledWith("/api/operations");
    });
  });

  describe("error handling", () => {
    it("throws with NETWORK_ERROR when response is not valid JSON", async () => {
      mockFetch.mockResolvedValue({
        ok: false,
        status: 500,
        json: () => Promise.reject(new Error("invalid json")),
      });

      await expect(api.calculate(1, "+", 1)).rejects.toMatchObject({
        code: "NETWORK_ERROR",
        message: "Request failed with status 500",
      });
    });

    it("throws when fetch itself fails (network error)", async () => {
      mockFetch.mockRejectedValue(new TypeError("Failed to fetch"));

      await expect(api.calculate(1, "+", 1)).rejects.toThrow("Failed to fetch");
    });
  });
});
