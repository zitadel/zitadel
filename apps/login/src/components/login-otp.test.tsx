import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { LoginOTP } from "./login-otp";

const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
    back: vi.fn(),
  }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/session", () => ({
  updateOrCreateSession: vi.fn(),
}));

vi.mock("@/lib/client", () => ({
  handleServerActionResponse: vi.fn(),
  completeFlowOrGetUrl: vi.fn(),
}));

describe("LoginOTP", () => {
  let mockUpdateOrCreateSession: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    vi.clearAllMocks();
    const { updateOrCreateSession } = await import("@/lib/server/session");
    mockUpdateOrCreateSession = vi.mocked(updateOrCreateSession);
  });

  afterEach(cleanup);

  test("should autofocus the code input on mount", () => {
    const { getByTestId } = render(<LoginOTP host={null} method="time-based" />);
    expect(getByTestId("code-text-input")).toHaveFocus();
  });

  test("should display translated fallback error when OTP verification request rejects", async () => {
    mockUpdateOrCreateSession.mockRejectedValueOnce(new Error("network down"));

    render(<LoginOTP host={null} method="time-based" />);

    const input = screen.getByTestId("code-text-input");
    fireEvent.change(input, { target: { value: "123456" } });

    const submitButton = screen.getByTestId("submit-button");
    await waitFor(() => {
      expect(submitButton).not.toBeDisabled();
    });

    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText("errors.couldNotVerifyCode")).toBeInTheDocument();
    });
  });

  test("should display translated fallback error when OTP challenge request rejects", async () => {
    mockUpdateOrCreateSession.mockRejectedValueOnce(new Error("challenge request failed"));

    render(<LoginOTP host={null} method="sms" />);

    expect(screen.getByRole("button", { name: "verify.resendCode" })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText("errors.couldNotRequestChallenge")).toBeInTheDocument();
    });
  });

  test("should render the back button when no recovery code alternative is offered", () => {
    const { queryByTestId } = render(<LoginOTP host={null} method="time-based" loginName="test@example.com" />);
    expect(queryByTestId("recovery-code-button")).not.toBeInTheDocument();
  });

  test("should render the recovery code button when altRecoveryCode is set", () => {
    const { getByTestId } = render(
      <LoginOTP host={null} method="time-based" loginName="test@example.com" altRecoveryCode />,
    );
    expect(getByTestId("recovery-code-button")).toBeInTheDocument();
  });

  test("should navigate to the recovery code page with the session params", () => {
    const { getByTestId } = render(
      <LoginOTP
        host={null}
        method="time-based"
        loginName="test@example.com"
        sessionId="session-123"
        requestId="oidc_123"
        organization="org-123"
        altRecoveryCode
      />,
    );

    fireEvent.click(getByTestId("recovery-code-button"));

    expect(mockPush).toHaveBeenCalledTimes(1);
    const target = mockPush.mock.calls[0][0] as string;
    expect(target.startsWith("/recovery-code?")).toBe(true);

    const params = new URLSearchParams(target.split("?")[1]);
    expect(params.get("loginName")).toBe("test@example.com");
    expect(params.get("sessionId")).toBe("session-123");
    expect(params.get("requestId")).toBe("oidc_123");
    expect(params.get("organization")).toBe("org-123");
  });
});
