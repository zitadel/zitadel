import { cleanup, fireEvent, render } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { TextButton, textLinkClasses } from "./text-button";

describe("TextButton", () => {
  afterEach(cleanup);

  it("renders a button that does not submit the surrounding form", () => {
    const { getByRole } = render(<TextButton>Reset Password</TextButton>);
    const button = getByRole("button", { name: "Reset Password" });
    expect(button.getAttribute("type")).toBe("button");
  });

  it("shows the shared text link look and keeps extra classes", () => {
    const { getByRole } = render(<TextButton className="ml-4">Resend</TextButton>);
    const { className } = getByRole("button");
    expect(className).toContain(textLinkClasses);
    expect(className).toContain("cursor-pointer");
    expect(className).toContain("font-semibold");
    expect(className).toContain("ml-4");
  });

  it("passes native button props through", () => {
    const onClick = vi.fn();
    const { getByTestId } = render(
      <TextButton onClick={onClick} data-testid="reset-button">
        Reset Password
      </TextButton>,
    );
    fireEvent.click(getByTestId("reset-button"));
    expect(onClick).toHaveBeenCalledTimes(1);
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
});
