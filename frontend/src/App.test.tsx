import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import App from "./App";

// Mock the API module so we don't make real fetch calls
vi.mock("./api/calculator", () => ({
  HttpCalculatorApi: class {
    calculate = vi.fn();
    evaluateExpression = vi.fn();
    getOperations = vi.fn();
  },
}));

describe("App", () => {
  it("renders without crashing", () => {
    render(<App />);
    // Calculator is present — verify by checking for the display
    expect(screen.getByLabelText("display value")).toBeInTheDocument();
  });

  it("calculator component is present with buttons", () => {
    render(<App />);
    expect(screen.getByText("=")).toBeInTheDocument();
    expect(screen.getByText("C")).toBeInTheDocument();
  });
});
