// API types
export interface CalculateRequest {
  operand_a: number;
  operator: string;
  operand_b: number;
}

export interface ExpressionRequest {
  expression: string;
}

export interface CalculateResponse {
  result: number;
}

export interface ExpressionResponse {
  result: number;
  expression: string;
}

export interface OperationsResponse {
  operations: string[];
}

export interface ApiError {
  error: {
    code: string;
    message: string;
  };
}

// Calculator UI types
export type ButtonVariant = "number" | "operator" | "action" | "equals";
export type ButtonSize = "default" | "wide";

export interface CalculatorButton {
  label: string;
  value: string;
  variant: ButtonVariant;
  size?: ButtonSize;
}
