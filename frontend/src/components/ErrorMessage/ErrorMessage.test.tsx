import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { ErrorMessage } from "./ErrorMessage";

describe("ErrorMessage", () => {
  it("renders nothing when message is null", () => {
    const { container } = render(<ErrorMessage message={null} />);
    expect(container.firstChild).toBeNull();
  });

  it("renders error message text", () => {
    render(<ErrorMessage message="division by zero is not allowed" />);
    expect(screen.getByText("division by zero is not allowed")).toBeInTheDocument();
  });

  it("has role='alert' for accessibility", () => {
    render(<ErrorMessage message="something went wrong" />);
    expect(screen.getByRole("alert")).toBeInTheDocument();
  });
});
