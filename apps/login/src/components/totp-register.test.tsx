import { cleanup, render } from "@testing-library/react";
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

  test("should show the secret as a manual setup key", () => {
    const { getByTestId, getByText } = render(
      <TotpRegister uri="otpauth://totp/test?secret=EXAMPLESECRET" secret="EXAMPLESECRET" />,
    );
    expect(getByText("set.manualSetupDescription")).toBeInTheDocument();
    expect(getByTestId("totp-secret")).toHaveTextContent(/^EXAMPLESECRET$/);
  });

  test("should not show a manual setup key without a secret", () => {
    const { queryByTestId, queryByText } = render(<TotpRegister uri="otpauth://totp/test" secret="" />);
    expect(queryByText("set.manualSetupDescription")).not.toBeInTheDocument();
    expect(queryByTestId("totp-secret")).not.toBeInTheDocument();
  });
});
