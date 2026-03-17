import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Display } from "./Display";

describe("Display", () => {
  it("renders current input", () => {
    render(
      <Display expression="" currentInput="123" result={null} isLoading={false} />
    );
    expect(screen.getByText("123")).toBeInTheDocument();
  });

  it("renders 0 when no input", () => {
    render(
      <Display expression="" currentInput="" result={null} isLoading={false} />
    );
    expect(screen.getByText("0")).toBeInTheDocument();
  });

  it("renders expression", () => {
    render(
      <Display expression="5 + 3" currentInput="" result={null} isLoading={false} />
    );
    expect(screen.getByText("5 + 3")).toBeInTheDocument();
  });

  it("renders result when present", () => {
    render(
      <Display expression="5 + 3" currentInput="" result="8" isLoading={false} />
    );
    expect(screen.getByText("8")).toBeInTheDocument();
  });

  it("shows loading indicator", () => {
    render(
      <Display expression="5 + 3" currentInput="" result={null} isLoading={true} />
    );
    expect(screen.getByText("...")).toBeInTheDocument();
  });

  it("formats large numbers with commas", () => {
    render(
      <Display expression="" currentInput="" result="1234567" isLoading={false} />
    );
    expect(screen.getByText("1,234,567")).toBeInTheDocument();
  });

  it("handles long numbers without overflow", () => {
    render(
      <Display
        expression=""
        currentInput="12345678901234"
        result={null}
        isLoading={false}
      />
    );
    // Should render without crashing — the overflow is handled by CSS
    expect(screen.getByText("12345678901234")).toBeInTheDocument();
  });

  it("has aria-live for result announcements", () => {
    render(
      <Display expression="" currentInput="5" result={null} isLoading={false} />
    );
    const displayValue = screen.getByLabelText("display value");
    expect(displayValue).toHaveAttribute("aria-live", "polite");
  });
});
