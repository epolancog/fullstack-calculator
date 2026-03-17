import type {
  CalculateResponse,
  ExpressionResponse,
  OperationsResponse,
  ApiError,
} from "../types";

export interface CalculatorApi {
  calculate(
    operandA: number,
    operator: string,
    operandB: number
  ): Promise<CalculateResponse>;
  evaluateExpression(expression: string): Promise<ExpressionResponse>;
  getOperations(): Promise<OperationsResponse>;
}

export class CalculatorApiError extends Error {
  code: string;

  constructor(code: string, message: string) {
    super(message);
    this.code = code;
    this.name = "CalculatorApiError";
  }
}

export class HttpCalculatorApi implements CalculatorApi {
  private baseUrl: string;

  constructor(baseUrl = "/api") {
    this.baseUrl = baseUrl;
  }

  async calculate(
    operandA: number,
    operator: string,
    operandB: number
  ): Promise<CalculateResponse> {
    const response = await fetch(`${this.baseUrl}/calculate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        operand_a: operandA,
        operator,
        operand_b: operandB,
      }),
    });

    return this.handleResponse<CalculateResponse>(response);
  }

  async evaluateExpression(
    expression: string
  ): Promise<ExpressionResponse> {
    const response = await fetch(`${this.baseUrl}/calculate/expression`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ expression }),
    });

    return this.handleResponse<ExpressionResponse>(response);
  }

  async getOperations(): Promise<OperationsResponse> {
    const response = await fetch(`${this.baseUrl}/operations`);
    return this.handleResponse<OperationsResponse>(response);
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
      let apiError: ApiError;
      try {
        apiError = (await response.json()) as ApiError;
      } catch {
        throw new CalculatorApiError(
          "NETWORK_ERROR",
          `Request failed with status ${response.status}`
        );
      }
      throw new CalculatorApiError(
        apiError.error.code,
        apiError.error.message
      );
    }

    return (await response.json()) as T;
  }
}
