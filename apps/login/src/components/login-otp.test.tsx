import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { LoginOTP } from "./login-otp";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/session", () => ({
  updateOrCreateSession: vi.fn(),
}));

vi.mock("@/lib/actions/flow", () => ({
  handleServerActionResponse: vi.fn(),
  completeFlowOrGetUrl: vi.fn(),
}));

describe("LoginOTP", () => {
  let mockUpdateOrCreateSession: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    vi.clearAllMocks();
    const { updateOrCreateSession } = await import("@/lib/actions/session");
    mockUpdateOrCreateSession = vi.mocked(updateOrCreateSession);
  });

  afterEach(cleanup);

  test("should autofocus the code input on mount", () => {
    const { getByTestId } = render(<LoginOTP method="time-based" />);
    expect(getByTestId("code-text-input")).toHaveFocus();
  });

  test("should display translated fallback error when OTP verification request rejects", async () => {
    mockUpdateOrCreateSession.mockRejectedValueOnce(new Error("network down"));

    render(<LoginOTP method="time-based" />);

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

    render(<LoginOTP method="sms" />);

    expect(screen.getByRole("button", { name: "verify.resendCode" })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText("errors.couldNotRequestChallenge")).toBeInTheDocument();
    });
  });

  test("should offer to restart the sign-in when the session of the step is lost", async () => {
    mockUpdateOrCreateSession.mockResolvedValueOnce({
      error: "sessionResolution.cookieMissing",
      sessionLost: true,
    });

    render(<LoginOTP method="sms" loginName="user@example.com" organization="org-1" requestId="oidc_123" />);

    await waitFor(() => {
      expect(screen.getByText("sessionResolution.cookieMissing")).toBeInTheDocument();
    });
    expect(screen.getByTestId("restart-sign-in")).toHaveAttribute(
      "href",
      "/loginname?loginName=user%40example.com&organization=org-1&requestId=oidc_123",
    );
  });

  test("should not offer to restart the sign-in for an error the user can retry", async () => {
    mockUpdateOrCreateSession.mockResolvedValueOnce({ error: "couldNotUpdateSession" });

    render(<LoginOTP method="sms" loginName="user@example.com" />);

    await waitFor(() => {
      expect(screen.getByText("couldNotUpdateSession")).toBeInTheDocument();
    });
    expect(screen.queryByTestId("restart-sign-in")).toBeNull();
  });

  // the restart link belongs to the answer that reported the lost session, not to a later error
  test.each([
    ["rejects", () => Promise.reject(new Error("network down")), "errors.couldNotVerifyCode"],
    [
      "returns an error the user can retry",
      () => Promise.resolve({ error: "couldNotUpdateSession" }),
      "couldNotUpdateSession",
    ],
  ])("should drop the restart link when a later request %s", async (_name, laterRequest, expectedError) => {
    mockUpdateOrCreateSession.mockResolvedValueOnce({ error: "sessionResolution.cookieMissing", sessionLost: true });

    render(<LoginOTP method="sms" loginName="user@example.com" />);

    await waitFor(() => {
      expect(screen.getByTestId("restart-sign-in")).toBeInTheDocument();
    });

    mockUpdateOrCreateSession.mockImplementationOnce(laterRequest);
    fireEvent.change(screen.getByTestId("code-text-input"), { target: { value: "123456" } });
    await waitFor(() => {
      expect(screen.getByTestId("submit-button")).not.toBeDisabled();
    });
    fireEvent.click(screen.getByTestId("submit-button"));

    await waitFor(() => {
      expect(screen.getByText(expectedError)).toBeInTheDocument();
    });
    expect(screen.queryByTestId("restart-sign-in")).toBeNull();
  });
});
