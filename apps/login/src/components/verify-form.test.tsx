import { cleanup, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { VerifyForm } from "./verify-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
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

    const { sendVerification } = await import("@/lib/server/verify");
    mockSendVerification = vi.mocked(sendVerification);
    mockSendVerification.mockResolvedValue({ redirect: "/success" });
  });

  afterEach(cleanup);

  describe("Input Focus", () => {
    test("should autofocus the code input on mount", () => {
      const { getByTestId } = render(
        <VerifyForm userId="user-1" code="" isInvite={false} submit={false} />,
      );
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
});
