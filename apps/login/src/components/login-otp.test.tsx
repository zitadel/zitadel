import { act, cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { LoginOTP } from "./login-otp";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string, data?: Record<string, unknown>) => (data ? `${key} ${JSON.stringify(data)}` : key),
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

  describe("SMS resend cooldown", () => {
    afterEach(() => {
      vi.useRealTimers();
    });

    test("disables resend and counts down after the code was sent", async () => {
      vi.useFakeTimers({ shouldAdvanceTime: true });
      mockUpdateOrCreateSession.mockResolvedValue({ sessionId: "session-1", retryAfterSeconds: 2 });

      render(<LoginOTP host={null} method="sms" />);

      const resendButton = screen.getByTestId("resend-button");
      await waitFor(() => expect(resendButton).toHaveTextContent('verify.resendCodeIn {"seconds":2}'));
      expect(resendButton).toBeDisabled();

      await act(async () => {
        await vi.advanceTimersByTimeAsync(1000);
      });
      expect(resendButton).toHaveTextContent('verify.resendCodeIn {"seconds":1}');

      await act(async () => {
        await vi.advanceTimersByTimeAsync(1000);
      });
      expect(resendButton).toHaveTextContent(/^verify\.resendCode$/);
      expect(resendButton).not.toBeDisabled();
    });

    test("shows the countdown instead of an error when the server refuses within the cooldown", async () => {
      mockUpdateOrCreateSession.mockResolvedValue({ error: "errors.resendCooldown", retryAfterSeconds: 25 });

      render(<LoginOTP host={null} method="sms" />);

      const resendButton = screen.getByTestId("resend-button");
      await waitFor(() => expect(resendButton).toHaveTextContent('verify.resendCodeIn {"seconds":25}'));
      expect(resendButton).toBeDisabled();
      expect(screen.queryByTestId("error")).not.toBeInTheDocument();
    });

    test("does not start a countdown for email codes", async () => {
      mockUpdateOrCreateSession.mockResolvedValue({ sessionId: "session-1" });

      render(<LoginOTP host="example.com" method="email" />);

      await waitFor(() => expect(mockUpdateOrCreateSession).toHaveBeenCalledTimes(1));
      const resendButton = screen.getByTestId("resend-button");
      await waitFor(() => expect(resendButton).not.toBeDisabled());
      expect(resendButton).toHaveTextContent(/^verify\.resendCode$/);
    });
  });
});
