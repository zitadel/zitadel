import { verifyTOTP } from "@/lib/server/verify";
import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { TotpRegister } from "./totp-register";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/verify", () => ({
  verifyTOTP: vi.fn(),
}));

vi.mock("@/lib/client", () => ({
  handleServerActionResponse: vi.fn(),
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("qrcode.react", () => ({
  QRCodeSVG: ({ value }: { value: string }) => <svg data-testid="qr-code" data-value={value} />,
}));

describe("TotpRegister", () => {
  afterEach(cleanup);

  test("should autofocus the code input on mount", () => {
    const { getByTestId } = render(<TotpRegister uri="otpauth://totp/test" secret="SECRET" />);
    expect(getByTestId("code-text-input")).toHaveFocus();
  });

  test("should show the error returned by verifyTOTP when code verification fails", async () => {
    vi.mocked(verifyTOTP).mockResolvedValue({ error: "The code is invalid. Please try again." });

    const { getByTestId, findByText } = render(
      <TotpRegister uri="otpauth://totp/test" secret="SECRET" loginName="user@example.com" />,
    );

    fireEvent.input(getByTestId("code-text-input"), { target: { value: "123456" } });
    await waitFor(() => expect(getByTestId("submit-button")).toBeEnabled());
    fireEvent.click(getByTestId("submit-button"));

    expect(await findByText("The code is invalid. Please try again.")).toBeInTheDocument();
  });
});
