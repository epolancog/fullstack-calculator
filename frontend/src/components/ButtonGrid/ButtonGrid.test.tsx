import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ButtonGrid } from "./ButtonGrid";

function createHandlers() {
  return {
    onDigit: vi.fn(),
    onDecimal: vi.fn(),
    onOperator: vi.fn(),
    onEquals: vi.fn(),
    onClear: vi.fn(),
    onBackspace: vi.fn(),
    onSqrt: vi.fn(),
    onPercent: vi.fn(),
    onOpenParen: vi.fn(),
    onCloseParen: vi.fn(),
  };
}

describe("ButtonGrid", () => {
  it("renders all number buttons (0-9)", () => {
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    for (let i = 0; i <= 9; i++) {
      expect(screen.getByText(String(i))).toBeInTheDocument();
    }
  });

  it("renders all operator buttons", () => {
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    expect(screen.getByRole("button", { name: "multiply" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "divide" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "subtract" })).toBeInTheDocument();
    expect(screen.getByText("+")).toBeInTheDocument();
    expect(screen.getByText("^")).toBeInTheDocument();
  });

  it("renders action buttons", () => {
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    expect(screen.getByText("C")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "backspace" })).toBeInTheDocument();
    expect(screen.getByText("=")).toBeInTheDocument();
  });

  it("renders advanced operation buttons", () => {
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    expect(screen.getByRole("button", { name: "square root" })).toBeInTheDocument();
    expect(screen.getByText("%")).toBeInTheDocument();
  });

  it("renders parentheses buttons", () => {
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    expect(screen.getByText("(")).toBeInTheDocument();
    expect(screen.getByText(")")).toBeInTheDocument();
  });

  it("clicking a number button fires onDigit", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByText("7"));
    expect(handlers.onDigit).toHaveBeenCalledWith("7");
  });

  it("clicking decimal fires onDecimal", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByText("."));
    expect(handlers.onDecimal).toHaveBeenCalled();
  });

  it("clicking operator fires onOperator", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByText("+"));
    expect(handlers.onOperator).toHaveBeenCalledWith("+");
  });

  it("clicking equals fires onEquals", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByText("="));
    expect(handlers.onEquals).toHaveBeenCalled();
  });

  it("clicking C fires onClear", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByText("C"));
    expect(handlers.onClear).toHaveBeenCalled();
  });

  it("clicking sqrt fires onSqrt", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByRole("button", { name: "square root" }));
    expect(handlers.onSqrt).toHaveBeenCalled();
  });

  it("clicking % fires onPercent", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByText("%"));
    expect(handlers.onPercent).toHaveBeenCalled();
  });

  it("clicking ( fires onOpenParen", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByText("("));
    expect(handlers.onOpenParen).toHaveBeenCalled();
  });

  it("clicking ) fires onCloseParen", async () => {
    const user = userEvent.setup();
    const handlers = createHandlers();
    render(<ButtonGrid {...handlers} />);

    await user.click(screen.getByText(")"));
    expect(handlers.onCloseParen).toHaveBeenCalled();
  });
});
