import { cleanup, render } from "@testing-library/react";
import { create } from "@zitadel/client";
import { PasswordComplexitySettingsSchema } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { afterEach, describe, expect, test, vi } from "vitest";
import { SetPasswordForm } from "./set-password-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/password", () => ({
  changePassword: vi.fn(),
  resetPassword: vi.fn(),
  sendPassword: vi.fn(),
}));

vi.mock("@/lib/client", () => ({
  handleServerActionResponse: vi.fn(),
}));

const defaultComplexitySettings = create(PasswordComplexitySettingsSchema, {
  minLength: 8n,
  requiresUppercase: false,
  requiresLowercase: false,
  requiresNumber: false,
  requiresSymbol: false,
});

describe("SetPasswordForm", () => {
  afterEach(cleanup);

  test("should autofocus the code input when codeRequired is true", () => {
    const { getByTestId } = render(
      <SetPasswordForm
        passwordComplexitySettings={defaultComplexitySettings}
        loginName="test@example.com"
        userId="user-1"
        codeRequired={true}
      />,
    );
    expect(getByTestId("code-text-input")).toHaveFocus();
  });

  test("should autofocus the password input when codeRequired is false", () => {
    const { getByTestId } = render(
      <SetPasswordForm
        passwordComplexitySettings={defaultComplexitySettings}
        loginName="test@example.com"
        userId="user-1"
        codeRequired={false}
      />,
    );
    expect(getByTestId("password-set-text-input")).toHaveFocus();
  });

  test.each([true, false])(
    "should render a hidden username field immediately before the password input (codeRequired=%s)",
    (codeRequired) => {
      const { container, getByTestId } = render(
        <SetPasswordForm
          passwordComplexitySettings={defaultComplexitySettings}
          loginName="test@example.com"
          userId="user-1"
          codeRequired={codeRequired}
        />,
      );

      const username = container.querySelector('input[autocomplete="username"]');
      expect(username).not.toBeNull();
      expect(username).toHaveValue("test@example.com");

      const password = getByTestId("password-set-text-input") as HTMLInputElement;
      const inputs = Array.from(container.querySelectorAll("input"));
      expect(inputs[inputs.indexOf(password) - 1]).toBe(username);
    },
  );
});
