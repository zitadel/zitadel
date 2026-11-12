import { changePassword, sendPassword } from "@/lib/server/password";
import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { create } from "@zitadel/client";
import { PasswordComplexitySettingsSchema } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { SetPasswordForm } from "./set-password-form";

const mockPush = vi.hoisted(() => vi.fn());

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
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
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    cleanup();
    vi.useRealTimers();
  });

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

  test("should use translated resend accessible name when codeRequired is true", () => {
    render(
      <SetPasswordForm
        passwordComplexitySettings={defaultComplexitySettings}
        loginName="test@example.com"
        userId="user-1"
        codeRequired={true}
      />,
    );

    expect(screen.getByRole("button", { name: "set.resend" })).toBeInTheDocument();
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

  // A second submit sets and checks the password again and rotates the session token while the
  // browser already loads the next step with the previous one.
  test("should not submit the password again while the redirect is pending", async () => {
    // the form waits 2 seconds between setting and checking the password
    vi.useFakeTimers({ shouldAdvanceTime: true });
    vi.mocked(changePassword).mockResolvedValue({} as any);
    vi.mocked(sendPassword).mockResolvedValue({ redirect: "/otp/sms?loginName=test%40example.com" });
    // the navigation has not finished yet
    mockPush.mockReturnValue(new Promise(() => {}));

    const { getByTestId } = render(
      <SetPasswordForm
        passwordComplexitySettings={defaultComplexitySettings}
        loginName="test@example.com"
        userId="user-1"
        codeRequired={false}
      />,
    );

    fireEvent.input(getByTestId("password-set-text-input"), { target: { value: "Password1!" } });
    fireEvent.input(getByTestId("password-set-confirm-text-input"), { target: { value: "Password1!" } });
    await waitFor(() => expect(getByTestId("submit-button")).toBeEnabled());

    fireEvent.click(getByTestId("submit-button"));
    await waitFor(() => expect(changePassword).toHaveBeenCalledTimes(1));
    await vi.advanceTimersByTimeAsync(2000);
    await waitFor(() => expect(mockPush).toHaveBeenCalledWith("/otp/sms?loginName=test%40example.com"));

    expect(getByTestId("submit-button")).toBeDisabled();

    fireEvent.click(getByTestId("submit-button"));
    await vi.advanceTimersByTimeAsync(2000);

    expect(changePassword).toHaveBeenCalledTimes(1);
    expect(sendPassword).toHaveBeenCalledTimes(1);
  });
});
