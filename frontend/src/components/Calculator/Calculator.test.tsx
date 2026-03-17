import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Calculator } from "./Calculator";
import type { CalculatorApi } from "../../api/calculator";

function createMockApi(): CalculatorApi {
  return {
    calculate: vi.fn(),
    evaluateExpression: vi.fn(),
    getOperations: vi.fn(),
  };
}

describe("Calculator", () => {
  let mockApi: CalculatorApi;

  beforeEach(() => {
    mockApi = createMockApi();
  });

  it("clicking digits updates display", async () => {
    const user = userEvent.setup();
    render(<Calculator api={mockApi} />);

    await user.click(screen.getByText("5"));
    await user.click(screen.getByText("3"));

    expect(screen.getByLabelText("display value")).toHaveTextContent("53");
  });

  it("clicking operator builds expression", async () => {
    const user = userEvent.setup();
    render(<Calculator api={mockApi} />);

    await user.click(screen.getByText("5"));
    await user.click(screen.getByText("+"));

    expect(screen.getByLabelText("expression")).toHaveTextContent("5 +");
  });

  it("clicking equals calls API and shows result", async () => {
    vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
      result: 8,
      expression: "5 + 3",
    });

    const user = userEvent.setup();
    render(<Calculator api={mockApi} />);

    await user.click(screen.getByText("5"));
    await user.click(screen.getByText("+"));
    await user.click(screen.getByText("3"));
    await user.click(screen.getByText("="));

    await waitFor(() => {
      expect(screen.getByLabelText("display value")).toHaveTextContent("8");
    });

    expect(mockApi.evaluateExpression).toHaveBeenCalledWith("5 + 3");
  });

  it("clicking clear resets the display", async () => {
    const user = userEvent.setup();
    render(<Calculator api={mockApi} />);

    await user.click(screen.getByText("5"));
    await user.click(screen.getByText("+"));
    await user.click(screen.getByText("3"));
    await user.click(screen.getByText("C"));

    expect(screen.getByLabelText("display value")).toHaveTextContent("0");
  });

  it("keyboard input works", async () => {
    vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
      result: 8,
      expression: "5 + 3",
    });

    const user = userEvent.setup();
    render(<Calculator api={mockApi} />);

    await user.keyboard("5+3{Enter}");

    await waitFor(() => {
      expect(screen.getByLabelText("display value")).toHaveTextContent("8");
    });
  });

  it("error from API shows error message", async () => {
    vi.mocked(mockApi.evaluateExpression).mockRejectedValue(
      new Error("division by zero is not allowed")
    );

    const user = userEvent.setup();
    render(<Calculator api={mockApi} />);

    await user.click(screen.getByText("1"));
    await user.click(screen.getByRole("button", { name: "0" }));
    await user.click(screen.getByRole("button", { name: "divide" }));
    await user.click(screen.getByRole("button", { name: "0" }));
    await user.click(screen.getByText("="));

    await waitFor(() => {
      expect(screen.getByRole("alert")).toHaveTextContent(
        "division by zero is not allowed"
      );
    });
  });

  it("full calculation flow: 5 + 3 * 2 = 11", async () => {
    vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
      result: 11,
      expression: "5 + 3 * 2",
    });

    const user = userEvent.setup();
    render(<Calculator api={mockApi} />);

    await user.click(screen.getByText("5"));
    await user.click(screen.getByText("+"));
    await user.click(screen.getByText("3"));
    await user.click(screen.getByRole("button", { name: "multiply" }));
    await user.click(screen.getByText("2"));
    await user.click(screen.getByText("="));

    await waitFor(() => {
      expect(screen.getByLabelText("display value")).toHaveTextContent("11");
    });

    expect(mockApi.evaluateExpression).toHaveBeenCalledWith("5 + 3 * 2");
  });

  it("sqrt with implicit multiplication: 5, √, 9 = 15", async () => {
    vi.mocked(mockApi.evaluateExpression).mockResolvedValue({
      result: 15,
      expression: "5 * sqrt 9",
    });

    const user = userEvent.setup();
    render(<Calculator api={mockApi} />);

    await user.click(screen.getByText("5"));
    await user.click(screen.getByRole("button", { name: "square root" }));
    await user.click(screen.getByText("9"));
    await user.click(screen.getByText("="));

    await waitFor(() => {
      expect(screen.getByLabelText("display value")).toHaveTextContent("15");
    });

    expect(mockApi.evaluateExpression).toHaveBeenCalledWith("5 * sqrt 9");
  });
});
