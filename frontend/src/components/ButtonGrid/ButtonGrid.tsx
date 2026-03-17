import { Button } from "../Button/Button";
import type { CalculatorButton } from "../../types";

interface ButtonGridProps {
  onDigit: (digit: string) => void;
  onDecimal: () => void;
  onOperator: (op: string) => void;
  onEquals: () => void;
  onClear: () => void;
  onBackspace: () => void;
  onSqrt: () => void;
  onPercent: () => void;
  onOpenParen: () => void;
  onCloseParen: () => void;
  disabled?: boolean;
}

const buttons: CalculatorButton[] = [
  // Row 1
  { label: "C", value: "clear", variant: "action" },
  { label: "⌫", value: "backspace", variant: "action" },
  { label: "(", value: "(", variant: "action" },
  { label: ")", value: ")", variant: "action" },
  // Row 2
  { label: "√", value: "sqrt", variant: "operator" },
  { label: "^", value: "^", variant: "operator" },
  { label: "%", value: "%", variant: "operator" },
  { label: "÷", value: "/", variant: "operator" },
  // Row 3
  { label: "7", value: "7", variant: "number" },
  { label: "8", value: "8", variant: "number" },
  { label: "9", value: "9", variant: "number" },
  { label: "×", value: "*", variant: "operator" },
  // Row 4
  { label: "4", value: "4", variant: "number" },
  { label: "5", value: "5", variant: "number" },
  { label: "6", value: "6", variant: "number" },
  { label: "−", value: "-", variant: "operator" },
  // Row 5
  { label: "1", value: "1", variant: "number" },
  { label: "2", value: "2", variant: "number" },
  { label: "3", value: "3", variant: "number" },
  { label: "+", value: "+", variant: "operator" },
  // Row 6
  { label: "0", value: "0", variant: "number", size: "wide" },
  { label: ".", value: ".", variant: "number" },
  { label: "=", value: "=", variant: "equals" },
];

export function ButtonGrid({
  onDigit,
  onDecimal,
  onOperator,
  onEquals,
  onClear,
  onBackspace,
  onSqrt,
  onPercent,
  onOpenParen,
  onCloseParen,
  disabled = false,
}: ButtonGridProps) {
  function handleClick(button: CalculatorButton) {
    if (disabled) return;

    switch (button.value) {
      case "clear":
        onClear();
        break;
      case "backspace":
        onBackspace();
        break;
      case "sqrt":
        onSqrt();
        break;
      case "%":
        onPercent();
        break;
      case "(":
        onOpenParen();
        break;
      case ")":
        onCloseParen();
        break;
      case ".":
        onDecimal();
        break;
      case "=":
        onEquals();
        break;
      case "+":
      case "-":
      case "*":
      case "/":
      case "^":
        onOperator(button.value);
        break;
      default:
        onDigit(button.value);
    }
  }

  return (
    <div className="grid grid-cols-4 gap-2">
      {buttons.map((button) => (
        <Button
          key={button.value === "0" ? "zero" : button.value}
          label={button.label}
          variant={button.variant}
          size={button.size}
          onClick={() => handleClick(button)}
          disabled={disabled}
          ariaLabel={
            button.label === "×"
              ? "multiply"
              : button.label === "÷"
                ? "divide"
                : button.label === "−"
                  ? "subtract"
                  : button.label === "⌫"
                    ? "backspace"
                    : button.label === "√"
                      ? "square root"
                      : undefined
          }
        />
      ))}
    </div>
  );
}
