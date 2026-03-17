import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useCalculator } from "./useCalculator";
import type { CalculatorApi } from "../api/calculator";

function createMockApi(): CalculatorApi {
  return {
    calculate: vi.fn(),
    evaluateExpression: vi.fn(),
    getOperations: vi.fn(),
  };
}

describe("useCalculator", () => {
  let mockApi: CalculatorApi;

  beforeEach(() => {
    mockApi = createMockApi();
  });

  describe("digit input", () => {
    it("appends digit to currentInput", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      expect(result.current.state.currentInput).toBe("5");

      act(() => result.current.handleDigit("3"));
      expect(result.current.state.currentInput).toBe("53");
    });

    it("prevents leading zeros", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("0"));
      expect(result.current.state.currentInput).toBe("0");

      act(() => result.current.handleDigit("0"));
      expect(result.current.state.currentInput).toBe("0");
    });

    it("replaces leading zero with non-zero digit", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("0"));
      act(() => result.current.handleDigit("5"));
      expect(result.current.state.currentInput).toBe("5");
    });
  });

  describe("decimal input", () => {
    it("appends decimal point", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("1"));
      act(() => result.current.handleDecimal());
      expect(result.current.state.currentInput).toBe("1.");
    });

    it("prevents double decimal", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("1"));
      act(() => result.current.handleDecimal());
      act(() => result.current.handleDecimal());
      expect(result.current.state.currentInput).toBe("1.");
    });

    it("starts with 0. when no current input", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDecimal());
      expect(result.current.state.currentInput).toBe("0.");
    });
  });

  describe("operator input", () => {
    it("flushes currentInput to expression", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));

      expect(result.current.state.expression).toBe("5 + ");
      expect(result.current.state.currentInput).toBe("");
    });

    it("replaces consecutive operators (non-minus)", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleOperator("*"));

      expect(result.current.state.expression).toBe("5 * ");
    });

    it("minus after operator starts negative number instead of replacing", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleOperator("-"));

      // "-" starts a negative number, doesn't replace "+"
      expect(result.current.state.expression).toBe("5 + ");
      expect(result.current.state.currentInput).toBe("-");
    });

    it("ignores operator with no input and no expression", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleOperator("+"));
      expect(result.current.state.expression).toBe("");
    });
  });

  describe("equals", () => {
    it("calls API and displays result", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: 8,
        expression: "5 + 3",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleDigit("3"));

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(mockApi.evaluateExpression).toHaveBeenCalledWith("5 + 3");
      expect(result.current.state.result).toBe("8");
      expect(result.current.state.isLoading).toBe(false);
    });

    it("does nothing with empty expression", async () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(mockApi.evaluateExpression).not.toHaveBeenCalled();
    });
  });

  describe("expression building", () => {
    it("builds 5 + 3 * 2 correctly", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: 11,
        expression: "5 + 3 * 2",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleDigit("3"));
      act(() => result.current.handleOperator("*"));
      act(() => result.current.handleDigit("2"));

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(mockApi.evaluateExpression).toHaveBeenCalledWith("5 + 3 * 2");
      expect(result.current.state.result).toBe("11");
    });
  });

  describe("clear and backspace", () => {
    it("clear resets all state", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleDigit("3"));
      act(() => result.current.handleClear());

      expect(result.current.state.expression).toBe("");
      expect(result.current.state.currentInput).toBe("");
      expect(result.current.state.result).toBeNull();
    });

    it("backspace removes last digit", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("1"));
      act(() => result.current.handleDigit("2"));
      act(() => result.current.handleDigit("3"));
      act(() => result.current.handleBackspace());

      expect(result.current.state.currentInput).toBe("12");
    });

    it("backspace does nothing when currentInput is empty", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleBackspace());
      expect(result.current.state.currentInput).toBe("");
    });
  });

  describe("error handling", () => {
    it("sets error when API returns error", async () => {
      vi.mocked(mockApi.evaluateExpression).mockRejectedValue(
        new Error("division by zero is not allowed")
      );

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("1"));
      act(() => result.current.handleOperator("/"));
      act(() => result.current.handleDigit("0"));

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(result.current.state.error).toBe("division by zero is not allowed");
      expect(result.current.state.isLoading).toBe(false);
    });

    it("sets error on network failure", async () => {
      vi.mocked(mockApi.evaluateExpression).mockRejectedValue(
        new TypeError("Failed to fetch")
      );

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("1"));

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(result.current.state.error).toBe("Failed to fetch");
    });
  });

  describe("after result behavior", () => {
    it("new digit starts fresh expression", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: 8,
        expression: "5 + 3",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleDigit("3"));

      await act(async () => {
        await result.current.handleEquals();
      });

      act(() => result.current.handleDigit("7"));

      expect(result.current.state.currentInput).toBe("7");
      expect(result.current.state.expression).toBe("");
      expect(result.current.state.result).toBeNull();
    });

    it("new operator continues with result", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: 8,
        expression: "5 + 3",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleDigit("3"));

      await act(async () => {
        await result.current.handleEquals();
      });

      act(() => result.current.handleOperator("+"));

      expect(result.current.state.expression).toBe("8 + ");
      expect(result.current.state.result).toBeNull();
    });
  });

  describe("negative number input", () => {
    it("pressing - as first input starts negative number", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleOperator("-"));
      expect(result.current.state.currentInput).toBe("-");
    });

    it("builds expression with negative first operand", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: -2,
        expression: "-5 + 3",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleOperator("-"));
      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleDigit("3"));

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(mockApi.evaluateExpression).toHaveBeenCalledWith("-5 + 3");
      expect(result.current.state.result).toBe("-2");
    });

    it("builds expression with negative after operator", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: 2,
        expression: "5 + -3",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleOperator("-"));
      act(() => result.current.handleDigit("3"));

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(mockApi.evaluateExpression).toHaveBeenCalledWith("5 + -3");
      expect(result.current.state.result).toBe("2");
    });
  });

  describe("sqrt", () => {
    it("appends sqrt to expression with no prior input", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleSqrt());
      expect(result.current.state.expression).toBe("sqrt ");
    });

    it("implicit multiplication when currentInput has digits", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: 15,
        expression: "5 * sqrt 9",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleSqrt());
      act(() => result.current.handleDigit("9"));

      expect(result.current.state.expression).toBe("5 * sqrt ");

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(mockApi.evaluateExpression).toHaveBeenCalledWith("5 * sqrt 9");
      expect(result.current.state.result).toBe("15");
    });

    it("no implicit multiply when expression ends with operator", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: 8,
        expression: "5 + sqrt 9",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleSqrt());
      act(() => result.current.handleDigit("9"));

      expect(result.current.state.expression).toBe("5 + sqrt ");

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(mockApi.evaluateExpression).toHaveBeenCalledWith("5 + sqrt 9");
      expect(result.current.state.result).toBe("8");
    });

    it("sqrt with no prior input evaluates correctly", async () => {
      vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
        result: 3,
        expression: "sqrt 9",
      });

      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleSqrt());
      act(() => result.current.handleDigit("9"));

      await act(async () => {
        await result.current.handleEquals();
      });

      expect(mockApi.evaluateExpression).toHaveBeenCalledWith("sqrt 9");
      expect(result.current.state.result).toBe("3");
    });
  });

  describe("percent", () => {
    it("appends % to currentInput", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleDigit("0"));
      act(() => result.current.handlePercent());

      expect(result.current.state.currentInput).toBe("50%");
    });

    it("allows multiple % presses", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleDigit("0"));
      act(() => result.current.handlePercent());
      act(() => result.current.handlePercent());

      expect(result.current.state.currentInput).toBe("50%%");
    });
  });

  describe("parentheses", () => {
    it("appends open paren to expression", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleOpenParen());
      expect(result.current.state.expression).toBe("(");
    });

    it("close paren flushes currentInput", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleOpenParen());
      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleOperator("+"));
      act(() => result.current.handleDigit("3"));
      act(() => result.current.handleCloseParen());

      expect(result.current.state.expression).toBe("(5 + 3)");
      expect(result.current.state.currentInput).toBe("");
    });

    it("prevents close paren when no matching open paren", () => {
      const { result } = renderHook(() => useCalculator(mockApi));

      act(() => result.current.handleDigit("5"));
      act(() => result.current.handleCloseParen());

      expect(result.current.state.expression).toBe("");
      expect(result.current.state.currentInput).toBe("5");
    });
  });
});
