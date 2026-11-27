import { sendPassword } from "@/lib/actions/password";
import { handleServerActionResponse } from "@/lib/client-utils";
import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { PasswordForm } from "./password-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/password", () => ({
  sendPassword: vi.fn(),
  resetPassword: vi.fn(),
}));

vi.mock("@/lib/client-utils", () => ({
  handleServerActionResponse: vi.fn(),
}));

describe("PasswordForm", () => {
  afterEach(cleanup);

  test("should autofocus the password input on mount", () => {
    const { getByTestId } = render(<PasswordForm loginSettings={undefined} loginName="test@example.com" />);
    expect(getByTestId("password-text-input")).toHaveFocus();
  });

  test("should hand the external callback of the action to the response handler once the action has settled", async () => {
    const response = { redirect: "https://app.example.com/callback?code=abc" };
    vi.mocked(sendPassword).mockResolvedValue(response);

    const { getByTestId } = render(
      <PasswordForm loginSettings={undefined} loginName="test@example.com" requestId="oidc_123" />,
    );

    fireEvent.input(getByTestId("password-text-input"), { target: { value: "Password1!" } });
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
    fireEvent.click(getByTestId("submit-button"));

    await waitFor(() =>
      expect(handleServerActionResponse).toHaveBeenCalledWith(response, expect.anything(), expect.any(Function)),
    );
    expect(sendPassword).toHaveBeenCalledWith(
      expect.objectContaining({ loginName: "test@example.com", requestId: "oidc_123" }),
    );
  });
});
