import { cleanup, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { ChangePasswordForm } from "./change-password-form";
import { create } from "@zitadel/client";
import { PasswordComplexitySettingsSchema } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/password", () => ({
  checkSessionAndSetPassword: vi.fn(),
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

describe("ChangePasswordForm", () => {
  afterEach(cleanup);

  test("should autofocus the current password input on mount", () => {
    const { getByTestId } = render(
      <ChangePasswordForm
        passwordComplexitySettings={defaultComplexitySettings}
        sessionId="session-1"
        loginName="test@example.com"
      />,
    );
    expect(getByTestId("password-change-current-text-input")).toHaveFocus();
  });
});
