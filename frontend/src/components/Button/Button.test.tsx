import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Button } from "./Button";

describe("Button", () => {
  it("renders label text", () => {
    render(<Button label="7" onClick={() => {}} />);
    expect(screen.getByText("7")).toBeInTheDocument();
  });

  it("applies number variant classes by default", () => {
    render(<Button label="5" onClick={() => {}} />);
    const button = screen.getByRole("button");
    expect(button.className).toContain("bg-glass-button");
  });

  it("applies operator variant classes", () => {
    render(<Button label="+" variant="operator" onClick={() => {}} />);
    const button = screen.getByRole("button");
    expect(button.className).toContain("bg-accent-muted");
  });

  it("applies action variant classes", () => {
    render(<Button label="C" variant="action" onClick={() => {}} />);
    const button = screen.getByRole("button");
    expect(button.className).toContain("bg-white/5");
  });

  it("applies equals variant classes", () => {
    render(<Button label="=" variant="equals" onClick={() => {}} />);
    const button = screen.getByRole("button");
    expect(button.className).toContain("bg-accent");
  });

  it("applies wide class for size='wide'", () => {
    render(<Button label="0" size="wide" onClick={() => {}} />);
    const button = screen.getByRole("button");
    expect(button.className).toContain("col-span-2");
  });

  it("calls onClick when clicked", async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();
    render(<Button label="1" onClick={handleClick} />);

    await user.click(screen.getByRole("button"));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it("renders with aria-label when provided", () => {
    render(<Button label="×" ariaLabel="multiply" onClick={() => {}} />);
    expect(screen.getByRole("button", { name: "multiply" })).toBeInTheDocument();
  });

  it("uses label as aria-label by default", () => {
    render(<Button label="7" onClick={() => {}} />);
    expect(screen.getByRole("button", { name: "7" })).toBeInTheDocument();
  });

  it("disabled button does not fire onClick", async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();
    render(<Button label="1" onClick={handleClick} disabled />);

    await user.click(screen.getByRole("button"));
    expect(handleClick).not.toHaveBeenCalled();
  });
});
