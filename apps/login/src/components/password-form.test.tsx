import { sendPassword } from "@/lib/server/password";
import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { PasswordForm } from "./password-form";

const mockPush = vi.hoisted(() => vi.fn());

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/password", () => ({
  sendPassword: vi.fn(),
  resetPassword: vi.fn(),
}));

describe("PasswordForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(cleanup);

  async function enterPassword(getByTestId: (id: string) => HTMLElement) {
    fireEvent.input(getByTestId("password-text-input"), { target: { value: "Password1!" } });
    await waitFor(() => expect(getByTestId("submit-button")).toBeEnabled());
  }

  test("should autofocus the password input on mount", () => {
    const { getByTestId } = render(<PasswordForm loginSettings={undefined} loginName="test@example.com" />);
    expect(getByTestId("password-text-input")).toHaveFocus();
  });

  // A second submit checks the password again and rotates the session token while the browser
  // already loads the next step with the previous one.
  test("should not submit the password again while the redirect is pending", async () => {
    vi.mocked(sendPassword).mockResolvedValue({ redirect: "/otp/sms?loginName=test%40example.com" });
    // the navigation has not finished yet
    mockPush.mockReturnValue(new Promise(() => {}));

    const { getByTestId } = render(<PasswordForm loginSettings={undefined} loginName="test@example.com" />);
    await enterPassword(getByTestId);

    fireEvent.click(getByTestId("submit-button"));
    await waitFor(() => expect(mockPush).toHaveBeenCalledWith("/otp/sms?loginName=test%40example.com"));

    expect(getByTestId("submit-button")).toBeDisabled();
    expect(getByTestId("reset-button")).toBeDisabled();

    fireEvent.click(getByTestId("submit-button"));
    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(sendPassword).toHaveBeenCalledTimes(1);
  });

  test("should allow another attempt after a wrong password", async () => {
    vi.mocked(sendPassword).mockResolvedValue({ error: "errors.failedToAuthenticate" });

    const { getByTestId } = render(<PasswordForm loginSettings={undefined} loginName="test@example.com" />);
    await enterPassword(getByTestId);

    fireEvent.click(getByTestId("submit-button"));
    await waitFor(() => expect(getByTestId("error")).toHaveTextContent("errors.failedToAuthenticate"));

    expect(mockPush).not.toHaveBeenCalled();
    expect(getByTestId("submit-button")).toBeEnabled();

    fireEvent.click(getByTestId("submit-button"));
    await waitFor(() => expect(sendPassword).toHaveBeenCalledTimes(2));
  });
});
