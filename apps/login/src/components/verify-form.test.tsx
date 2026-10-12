import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { VerifyForm } from "./verify-form";

let mockSearchParams = new URLSearchParams();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  usePathname: () => "/verify",
  useSearchParams: () => mockSearchParams,
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/verify", () => ({
  sendVerification: vi.fn(),
  resendVerification: vi.fn(),
}));

describe("VerifyForm", () => {
  let mockSendVerification: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    vi.clearAllMocks();
    mockSearchParams = new URLSearchParams();

    const { sendVerification } = await import("@/lib/server/verify");
    mockSendVerification = vi.mocked(sendVerification);
    mockSendVerification.mockResolvedValue({ redirect: "/success" });
  });

  afterEach(cleanup);

  describe("Input Focus", () => {
    test("should autofocus the code input on mount", () => {
      const { getByTestId } = render(<VerifyForm userId="user-1" code="" isInvite={false} submit={false} />);
      expect(getByTestId("code-text-input")).toHaveFocus();
    });
  });

  describe("Auto-submit Behavior", () => {
    test("should call sendVerification automatically when submit=true", async () => {
      render(<VerifyForm userId="user-1" code="123456" isInvite={false} submit={true} />);

      await waitFor(() => {
        expect(mockSendVerification).toHaveBeenCalledWith(
          expect.objectContaining({
            code: "123456",
            userId: "user-1",
          }),
        );
      });
    });

    test("should prefill code but not auto-submit when submit=false", () => {
      render(<VerifyForm userId="user-1" code="123456" isInvite={false} submit={false} />);

      const input = screen.getByTestId("code-text-input");
      expect(input).toHaveValue("123456");

      const submitButton = screen.getByTestId("submit-button");
      expect(submitButton).toBeInTheDocument();

      expect(mockSendVerification).not.toHaveBeenCalled();
    });
  });

  describe("Resend cooldown", () => {
    test("should start the cooldown on arrival when a code was just sent (codeSent=true)", () => {
      mockSearchParams = new URLSearchParams({ codeSent: "true" });

      render(<VerifyForm userId="user-1" code="" isInvite={false} submit={false} />);

      expect(screen.getByTestId("resend-button")).toBeDisabled();
      expect(screen.getByTestId("resend-button")).toHaveTextContent("verify.resendCode (30s)");
    });

    test("should not start the cooldown on arrival without codeSent", () => {
      render(<VerifyForm userId="user-1" code="" isInvite={false} submit={false} />);

      expect(screen.getByTestId("resend-button")).not.toBeDisabled();
    });

    test("should include the remaining seconds in the accessible name while cooling down", () => {
      mockSearchParams = new URLSearchParams({ codeSent: "true" });

      render(<VerifyForm userId="user-1" code="" isInvite={false} submit={false} />);

      expect(screen.getByRole("button", { name: "verify.resendCode (30s)" })).toBeInTheDocument();
    });

    test("should disable the resend button and show the countdown after a successful resend", async () => {
      const { resendVerification } = await import("@/lib/server/verify");
      vi.mocked(resendVerification).mockResolvedValue({} as never);

      render(<VerifyForm userId="user-1" code="" isInvite={false} submit={false} />);

      const button = screen.getByTestId("resend-button");
      expect(button).not.toBeDisabled();
      expect(button).not.toHaveTextContent("(30s)");

      fireEvent.click(button);

      await waitFor(() => {
        expect(screen.getByTestId("resend-button")).toBeDisabled();
      });
      expect(screen.getByTestId("resend-button")).toHaveTextContent("verify.resendCode (30s)");
    });

    test("should not start the cooldown when the resend fails", async () => {
      const { resendVerification } = await import("@/lib/server/verify");
      vi.mocked(resendVerification).mockResolvedValue({ error: "boom" });

      render(<VerifyForm userId="user-1" code="" isInvite={false} submit={false} />);
      fireEvent.click(screen.getByTestId("resend-button"));

      await waitFor(() => {
        expect(screen.getByTestId("error")).toBeInTheDocument();
      });
      expect(screen.getByTestId("resend-button")).not.toBeDisabled();
    });
  });
});
