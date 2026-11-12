import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { LoginRecoveryCode } from "./login-recovery-code";

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

vi.mock("@/lib/server/recovery-code", () => ({
  sendRecoveryCode: vi.fn(),
}));

describe("LoginRecoveryCode", () => {
  let mockSendRecoveryCode: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    mockPush.mockClear();

    const { sendRecoveryCode } = await import("@/lib/server/recovery-code");
    mockSendRecoveryCode = vi.mocked(sendRecoveryCode);
    mockSendRecoveryCode.mockResolvedValue({ redirect: "/signedin" });
  });

  afterEach(cleanup);

  async function submitCode(code: string) {
    const input = screen.getByTestId("code-text-input");
    fireEvent.change(input, { target: { value: code } });

    const submitButton = screen.getByTestId("submit-button");
    await waitFor(() => expect(submitButton).not.toBeDisabled());

    fireEvent.click(submitButton);
  }

  test("should autofocus the code input on mount", () => {
    render(<LoginRecoveryCode loginName="test@example.com" />);
    expect(screen.getByTestId("code-text-input")).toHaveFocus();
  });

  test("should not autocorrect or autocapitalize the code input", () => {
    render(<LoginRecoveryCode loginName="test@example.com" />);

    const input = screen.getByTestId("code-text-input");
    expect(input).toHaveAttribute("autocapitalize", "off");
    expect(input).toHaveAttribute("autocomplete", "off");
    expect(input).toHaveAttribute("spellcheck", "false");
  });

  test("should keep the submit button disabled until a code is entered", () => {
    render(<LoginRecoveryCode loginName="test@example.com" />);
    expect(screen.getByTestId("submit-button")).toBeDisabled();
  });

  test("should send the code together with the session props", async () => {
    render(
      <LoginRecoveryCode loginName="test@example.com" sessionId="session-123" requestId="oidc_123" organization="org-123" />,
    );

    await submitCode("abcde-fghij");

    await waitFor(() => {
      expect(mockSendRecoveryCode).toHaveBeenCalledWith({
        loginName: "test@example.com",
        sessionId: "session-123",
        requestId: "oidc_123",
        organization: "org-123",
        code: "abcde-fghij",
      });
    });
  });

  test("should redirect when the server action returns a redirect", async () => {
    mockSendRecoveryCode.mockResolvedValue({ redirect: "/signedin" });

    render(<LoginRecoveryCode loginName="test@example.com" />);

    await submitCode("abcde-fghij");

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith("/signedin");
    });
  });

  test("should render the error returned by the server action", async () => {
    mockSendRecoveryCode.mockResolvedValue({ error: "The recovery code is invalid or has already been used." });

    render(<LoginRecoveryCode loginName="test@example.com" />);

    await submitCode("wrong");

    await waitFor(() => {
      expect(screen.getByTestId("error")).toHaveTextContent("The recovery code is invalid or has already been used.");
    });
    expect(mockPush).not.toHaveBeenCalled();
  });

  test("should render a generic error when the server action throws", async () => {
    mockSendRecoveryCode.mockRejectedValue(new Error("network down"));

    render(<LoginRecoveryCode loginName="test@example.com" />);

    await submitCode("abcde-fghij");

    await waitFor(() => {
      expect(screen.getByTestId("error")).toHaveTextContent("verify.errors.couldNotVerifyCode");
    });
  });

  test("should render a generic error when the server action returns nothing", async () => {
    mockSendRecoveryCode.mockResolvedValue(undefined);

    render(<LoginRecoveryCode loginName="test@example.com" />);

    await submitCode("abcde-fghij");

    await waitFor(() => {
      expect(screen.getByTestId("error")).toHaveTextContent("verify.errors.couldNotVerifyCode");
    });
  });
});
