import { cleanup, fireEvent, render } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { TextButton } from "./text-button";

describe("TextButton", () => {
  afterEach(cleanup);

  it("renders a button that does not submit the surrounding form", () => {
    const { getByRole } = render(<TextButton>Reset Password</TextButton>);
    const button = getByRole("button", { name: "Reset Password" });
    expect(button.getAttribute("type")).toBe("button");
  });

  it("passes native button props and extra classes through", () => {
    const onClick = vi.fn();
    const { getByTestId } = render(
      <TextButton onClick={onClick} className="ml-4" data-testid="resend-button">
        Resend code
      </TextButton>,
    );
    const button = getByTestId("resend-button");
    fireEvent.click(button);
    expect(onClick).toHaveBeenCalledTimes(1);
    expect(button.className).toContain("ml-4");
  });

  it("does not react to clicks while disabled", () => {
    const onClick = vi.fn();
    const { getByRole } = render(
      <TextButton onClick={onClick} disabled>
        Reset Password
      </TextButton>,
    );
    fireEvent.click(getByRole("button"));
    expect(onClick).not.toHaveBeenCalled();
  });

  it("colours a button next to other text like a link", () => {
    const { getByRole, rerender } = render(<TextButton>Skip</TextButton>);
    const standalone = getByRole("button").className;
    rerender(<TextButton inline>Resend code</TextButton>);
    expect(getByRole("button").className).not.toBe(standalone);
  });
});
