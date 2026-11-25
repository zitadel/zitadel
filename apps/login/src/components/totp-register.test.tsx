import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { TotpRegister } from "./totp-register";

const { push } = vi.hoisted(() => ({ push: vi.fn() }));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/verify", () => ({
  verifyTOTP: vi.fn(),
}));

vi.mock("@/lib/actions/flow", () => ({
  handleServerActionResponse: vi.fn(),
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("qrcode.react", () => ({
  QRCodeSVG: ({ value }: { value: string }) => <svg data-testid="qr-code" data-value={value} />,
}));

describe("TotpRegister", () => {
  let mockVerifyTOTP: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    vi.clearAllMocks();
    mockVerifyTOTP = vi.mocked((await import("@/lib/actions/verify")).verifyTOTP);
  });

  afterEach(cleanup);

  async function submitCode() {
    render(<TotpRegister uri="otpauth://totp/test" secret="SECRET" loginName="user@example.com" checkAfter />);

    fireEvent.change(screen.getByTestId("code-text-input"), { target: { value: "123456" } });
    const submitButton = screen.getByTestId("submit-button");
    await waitFor(() => expect(submitButton).not.toBeDisabled());
    fireEvent.click(submitButton);
  }

  test("should autofocus the code input on mount", () => {
    const { getByTestId } = render(<TotpRegister uri="otpauth://totp/test" secret="SECRET" />);
    expect(getByTestId("code-text-input")).toHaveFocus();
  });

  test("shows the message of a refused code and stays on the page", async () => {
    mockVerifyTOTP.mockResolvedValue({ error: "Invalid code. Please try again." });

    await submitCode();

    await waitFor(() => expect(screen.getByText("Invalid code. Please try again.")).toBeInTheDocument());
    expect(push).not.toHaveBeenCalled();
  });

  test("continues after a verified code", async () => {
    mockVerifyTOTP.mockResolvedValue({});

    await submitCode();

    await waitFor(() => expect(push).toHaveBeenCalledWith("/otp/time-based?loginName=user%40example.com"));
  });
});
