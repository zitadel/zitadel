import { cleanup, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { PasswordForm } from "./password-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/password", () => ({
  sendPassword: vi.fn(),
  resetPassword: vi.fn(),
}));

describe("PasswordForm", () => {
  afterEach(cleanup);

  test("should autofocus the password input on mount", () => {
    const { getByTestId } = render(<PasswordForm loginSettings={undefined} loginName="test@example.com" />);
    expect(getByTestId("password-text-input")).toHaveFocus();
  });

  test("should set autocomplete=current-password on the password input", () => {
    const { getByTestId } = render(<PasswordForm loginSettings={undefined} loginName="test@example.com" />);
    expect(getByTestId("password-text-input")).toHaveAttribute("autocomplete", "current-password");
  });

  test("should render the hidden username field before the password input when loginName is provided", () => {
    const { container, getByTestId } = render(<PasswordForm loginSettings={undefined} loginName="test@example.com" />);

    const username = container.querySelector('input[autocomplete="username"]');
    expect(username).not.toBeNull();
    expect(username).toHaveValue("test@example.com");

    // password managers pair the password with a preceding username field
    const password = getByTestId("password-text-input");
    expect(username!.compareDocumentPosition(password) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
  });
});
