import { cleanup, render } from "@testing-library/react";
import { create } from "@zitadel/client";
import { LoginSettingsSchema } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { afterEach, describe, expect, test, vi } from "vitest";
import { PasswordForm } from "./password-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/password", () => ({
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

  test.each([false, true])(
    "should render one username field before the password input (hidePasswordReset=%s)",
    (hidePasswordReset) => {
      const { container, getByTestId } = render(
        <PasswordForm loginSettings={create(LoginSettingsSchema, { hidePasswordReset })} loginName="test@example.com" />,
      );

      const usernames = container.querySelectorAll('input[autocomplete="username"]');
      expect(usernames).toHaveLength(1);
      expect(usernames[0]).toHaveValue("test@example.com");

      // Opera ignores type="hidden" username fields, so this must be a visually hidden text input
      expect(usernames[0]).toHaveAttribute("type", "text");
      expect(usernames[0]).toHaveAttribute("readonly");
      expect(usernames[0]).toHaveAttribute("tabindex", "-1");

      // password managers pair the password with a preceding username field
      const password = getByTestId("password-text-input");
      expect(usernames[0].compareDocumentPosition(password) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
    },
  );
});
