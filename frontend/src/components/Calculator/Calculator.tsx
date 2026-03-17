import { useEffect, useCallback } from "react";
import type { CalculatorApi } from "../../api/calculator";
import { useCalculator } from "../../hooks/useCalculator";
import { Display } from "../Display/Display";
import { ButtonGrid } from "../ButtonGrid/ButtonGrid";
import { ErrorMessage } from "../ErrorMessage/ErrorMessage";

interface CalculatorProps {
  api: CalculatorApi;
}

export function Calculator({ api }: CalculatorProps) {
  const {
    state,
    handleDigit,
    handleDecimal,
    handleOperator,
    handleEquals,
    handleClear,
    handleBackspace,
    handleSqrt,
    handlePercent,
    handleOpenParen,
    handleCloseParen,
  } = useCalculator(api);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key >= "0" && e.key <= "9") {
        handleDigit(e.key);
      } else if (e.key === ".") {
        handleDecimal();
      } else if (e.key === "+" || e.key === "-" || e.key === "*" || e.key === "/" || e.key === "^") {
        handleOperator(e.key);
      } else if (e.key === "%") {
        handlePercent();
      } else if (e.key === "(" ) {
        handleOpenParen();
      } else if (e.key === ")") {
        handleCloseParen();
      } else if (e.key === "Enter" || e.key === "=") {
        e.preventDefault();
        handleEquals();
      } else if (e.key === "Escape") {
        handleClear();
      } else if (e.key === "Backspace") {
        handleBackspace();
      }
    },
    [
      handleDigit,
      handleDecimal,
      handleOperator,
      handleEquals,
      handleClear,
      handleBackspace,
      handlePercent,
      handleOpenParen,
      handleCloseParen,
    ]
  );

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [handleKeyDown]);

  return (
    <div className="flex flex-col gap-4 w-full max-w-sm p-4">
      <Display
        expression={state.expression}
        currentInput={state.currentInput}
        result={state.result}
        isLoading={state.isLoading}
      />
      <ErrorMessage message={state.error} />
      <ButtonGrid
        onDigit={handleDigit}
        onDecimal={handleDecimal}
        onOperator={handleOperator}
        onEquals={handleEquals}
        onClear={handleClear}
        onBackspace={handleBackspace}
        onSqrt={handleSqrt}
        onPercent={handlePercent}
        onOpenParen={handleOpenParen}
        onCloseParen={handleCloseParen}
        disabled={state.isLoading}
      />
    </div>
  );
}
