import { useReducer, useCallback } from "react";
import type { CalculatorApi } from "../api/calculator";

interface CalculatorState {
  expression: string;
  currentInput: string;
  result: string | null;
  error: string | null;
  isLoading: boolean;
}

type CalculatorAction =
  | { type: "DIGIT"; payload: string }
  | { type: "DECIMAL" }
  | { type: "OPERATOR"; payload: string }
  | { type: "EQUALS" }
  | { type: "CLEAR" }
  | { type: "BACKSPACE" }
  | { type: "CLEAR_ENTRY" }
  | { type: "SQRT" }
  | { type: "PERCENT" }
  | { type: "OPEN_PAREN" }
  | { type: "CLOSE_PAREN" }
  | { type: "SET_RESULT"; payload: string }
  | { type: "SET_ERROR"; payload: string }
  | { type: "SET_LOADING" };

const initialState: CalculatorState = {
  expression: "",
  currentInput: "",
  result: null,
  error: null,
  isLoading: false,
};

function countOpenParens(expression: string): number {
  let count = 0;
  for (const ch of expression) {
    if (ch === "(") count++;
    if (ch === ")") count--;
  }
  return count;
}

function expressionEndsWithOperator(expression: string): boolean {
  const trimmed = expression.trimEnd();
  if (trimmed.length === 0) return false;
  const lastChar = trimmed[trimmed.length - 1];
  return ["+", "-", "*", "/", "^"].includes(lastChar);
}

function calculatorReducer(
  state: CalculatorState,
  action: CalculatorAction
): CalculatorState {
  switch (action.type) {
    case "DIGIT": {
      const digit = action.payload;

      // After result, new digit starts fresh
      if (state.result !== null) {
        return {
          ...initialState,
          currentInput: digit,
        };
      }

      // Prevent leading zeros (except "0.")
      if (state.currentInput === "0" && digit !== "0") {
        return { ...state, currentInput: digit, error: null };
      }
      if (state.currentInput === "0" && digit === "0") {
        return state;
      }

      return {
        ...state,
        currentInput: state.currentInput + digit,
        error: null,
      };
    }

    case "DECIMAL": {
      // After result, start fresh with "0."
      if (state.result !== null) {
        return {
          ...initialState,
          currentInput: "0.",
        };
      }

      // Prevent double decimal
      if (state.currentInput.includes(".")) return state;

      // Start with "0." if no current input
      if (state.currentInput === "" || state.currentInput === "-") {
        return {
          ...state,
          currentInput: state.currentInput === "-" ? "-0." : "0.",
          error: null,
        };
      }

      return {
        ...state,
        currentInput: state.currentInput + ".",
        error: null,
      };
    }

    case "OPERATOR": {
      const op = action.payload;

      // After result, continue with result
      if (state.result !== null) {
        return {
          ...initialState,
          expression: state.result + " " + op + " ",
          currentInput: "",
        };
      }

      // Negative number: pressing "-" with no input and empty/operator-ending expression
      if (
        op === "-" &&
        state.currentInput === "" &&
        (state.expression === "" || expressionEndsWithOperator(state.expression) || state.expression.trimEnd().endsWith("("))
      ) {
        return { ...state, currentInput: "-", error: null };
      }

      // No input and no expression — ignore operator (except "-" handled above)
      if (state.currentInput === "" && state.expression === "") {
        return state;
      }

      // Replace consecutive operator
      if (state.currentInput === "" && expressionEndsWithOperator(state.expression)) {
        const trimmed = state.expression.trimEnd();
        // Find last operator and replace
        const withoutOp = trimmed.slice(0, trimmed.length - 1).trimEnd();
        return {
          ...state,
          expression: withoutOp + " " + op + " ",
          error: null,
        };
      }

      // Flush currentInput to expression
      return {
        ...state,
        expression: state.expression + state.currentInput + " " + op + " ",
        currentInput: "",
        error: null,
      };
    }

    case "EQUALS": {
      // Nothing to evaluate
      if (state.expression === "" && state.currentInput === "") return state;

      const fullExpression = state.expression + state.currentInput;
      if (fullExpression.trim() === "") return state;

      return {
        ...state,
        expression: fullExpression,
        currentInput: "",
        isLoading: true,
        error: null,
      };
    }

    case "CLEAR":
      return { ...initialState };

    case "BACKSPACE": {
      if (state.result !== null) return state;
      if (state.currentInput === "") return state;

      return {
        ...state,
        currentInput: state.currentInput.slice(0, -1),
      };
    }

    case "CLEAR_ENTRY": {
      return { ...state, currentInput: "", error: null };
    }

    case "SQRT": {
      // After result, apply sqrt to result
      if (state.result !== null) {
        return {
          ...initialState,
          expression: "sqrt " + state.result + " ",
          currentInput: "",
          isLoading: false,
        };
      }

      // Implicit multiplication: if there are digits in currentInput, flush with "*"
      if (state.currentInput !== "" && state.currentInput !== "-") {
        return {
          ...state,
          expression: state.expression + state.currentInput + " * sqrt ",
          currentInput: "",
          error: null,
        };
      }

      return {
        ...state,
        expression: state.expression + "sqrt ",
        currentInput: "",
        error: null,
      };
    }

    case "PERCENT": {
      // After result, apply % to result
      if (state.result !== null) {
        return {
          ...initialState,
          currentInput: state.result + "%",
        };
      }

      // Append % to currentInput (even if empty — backend handles "%" → 0)
      return {
        ...state,
        currentInput: state.currentInput + "%",
        error: null,
      };
    }

    case "OPEN_PAREN": {
      // After result, start fresh with "("
      if (state.result !== null) {
        return {
          ...initialState,
          expression: "(",
        };
      }

      // If there's currentInput, flush with implicit multiply
      if (state.currentInput !== "" && state.currentInput !== "-") {
        return {
          ...state,
          expression: state.expression + state.currentInput + " * (",
          currentInput: "",
          error: null,
        };
      }

      return {
        ...state,
        expression: state.expression + "(",
        error: null,
      };
    }

    case "CLOSE_PAREN": {
      // Only allow if there are open parens to close
      const fullExpr = state.expression + state.currentInput;
      if (countOpenParens(fullExpr) <= 0) return state;

      return {
        ...state,
        expression: state.expression + state.currentInput + ")",
        currentInput: "",
        error: null,
      };
    }

    case "SET_RESULT":
      return {
        ...state,
        result: action.payload,
        isLoading: false,
        error: null,
      };

    case "SET_ERROR":
      return {
        ...state,
        error: action.payload,
        isLoading: false,
      };

    case "SET_LOADING":
      return { ...state, isLoading: true };

    default:
      return state;
  }
}

export function useCalculator(api: CalculatorApi) {
  const [state, dispatch] = useReducer(calculatorReducer, initialState);

  const handleDigit = useCallback((digit: string) => {
    dispatch({ type: "DIGIT", payload: digit });
  }, []);

  const handleDecimal = useCallback(() => {
    dispatch({ type: "DECIMAL" });
  }, []);

  const handleOperator = useCallback((op: string) => {
    dispatch({ type: "OPERATOR", payload: op });
  }, []);

  const handleEquals = useCallback(async () => {
    const fullExpression =
      state.expression + state.currentInput;
    if (fullExpression.trim() === "") return;

    dispatch({ type: "EQUALS" });

    try {
      const response = await api.evaluateExpression(
        fullExpression.trim()
      );
      dispatch({ type: "SET_RESULT", payload: String(response.result) });
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "An unexpected error occurred";
      dispatch({ type: "SET_ERROR", payload: message });
    }
  }, [api, state.expression, state.currentInput]);

  const handleClear = useCallback(() => {
    dispatch({ type: "CLEAR" });
  }, []);

  const handleBackspace = useCallback(() => {
    dispatch({ type: "BACKSPACE" });
  }, []);

  const handleClearEntry = useCallback(() => {
    dispatch({ type: "CLEAR_ENTRY" });
  }, []);

  const handleSqrt = useCallback(() => {
    dispatch({ type: "SQRT" });
  }, []);

  const handlePercent = useCallback(() => {
    dispatch({ type: "PERCENT" });
  }, []);

  const handleOpenParen = useCallback(() => {
    dispatch({ type: "OPEN_PAREN" });
  }, []);

  const handleCloseParen = useCallback(() => {
    dispatch({ type: "CLOSE_PAREN" });
  }, []);

  return {
    state,
    handleDigit,
    handleDecimal,
    handleOperator,
    handleEquals,
    handleClear,
    handleBackspace,
    handleClearEntry,
    handleSqrt,
    handlePercent,
    handleOpenParen,
    handleCloseParen,
  };
}
