import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { create } from "@zitadel/client";
import { LoginSettingsSchema, PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { LoginForm } from "./login-form";

const push = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push, back: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/loginname", () => ({
  sendLoginname: vi.fn(),
}));

vi.mock("@/lib/actions/password", () => ({
  sendPassword: vi.fn(),
  resetPassword: vi.fn(),
}));

const passkeySettings = create(LoginSettingsSchema, { passkeysType: PasskeysType.ALLOWED });

function renderForm(props: Partial<Parameters<typeof LoginForm>[0]> = {}) {
  return render(<LoginForm loginName="" requestId={undefined} loginSettings={passkeySettings} submit={false} {...props} />);
}

function type(input: HTMLElement, value: string) {
  fireEvent.change(input, { target: { value } });
}

describe("LoginForm", () => {
  let sendLoginname: any;
  let sendPassword: any;
  let resetPassword: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    sendLoginname = vi.mocked((await import("@/lib/actions/loginname")).sendLoginname);
    const passwordActions = await import("@/lib/actions/password");
    sendPassword = vi.mocked(passwordActions.sendPassword);
    resetPassword = vi.mocked(passwordActions.resetPassword);
  });

  afterEach(cleanup);

  test("should autofocus the email input on mount", () => {
    const { getByTestId } = renderForm();
    expect(getByTestId("username-text-input")).toHaveFocus();
  });

  test("asks for the email and the password on one screen", () => {
    const { getByTestId, getByText } = renderForm();

    expect(getByTestId("username-text-input").getAttribute("type")).toBe("email");
    expect(getByTestId("password-text-input").getAttribute("type")).toBe("password");
    expect(getByText("labels.email")).toBeTruthy();
  });

  test("labels the field by the login names the settings allow", () => {
    const settings = create(LoginSettingsSchema, { disableLoginWithEmail: true, disableLoginWithPhone: true });
    const { getByTestId, getByText } = renderForm({ loginSettings: settings });

    expect(getByText("labels.username")).toBeTruthy();
    expect(getByTestId("username-text-input").getAttribute("type")).toBe("text");
  });

  test("signs in with the email and the password", async () => {
    sendPassword.mockResolvedValue({ redirect: "/signedin" });
    const { getByTestId } = renderForm({ organization: "org-1", requestId: "oidc_1" });

    type(getByTestId("username-text-input"), "john@example.com");
    type(getByTestId("password-text-input"), "Password1!");
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
    fireEvent.click(getByTestId("submit-button"));

    await waitFor(() => expect(push).toHaveBeenCalledWith("/signedin"));
    expect(sendPassword).toHaveBeenCalledWith(
      expect.objectContaining({ loginName: "john@example.com", organization: "org-1", requestId: "oidc_1" }),
    );
    expect(sendPassword.mock.calls[0][0].checks.password.password).toBe("Password1!");
    expect(sendLoginname).not.toHaveBeenCalled();
  });

  test("adds the organization domain to the login name before checking the password", async () => {
    sendPassword.mockResolvedValue({ redirect: "/signedin" });
    const { getByTestId } = renderForm({ suffix: "acme.test" });

    type(getByTestId("username-text-input"), "amy");
    type(getByTestId("password-text-input"), "Password1!");
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
    fireEvent.click(getByTestId("submit-button"));

    await waitFor(() => expect(sendPassword).toHaveBeenCalled());
    expect(sendPassword.mock.calls[0][0].loginName).toBe("amy@acme.test");
  });

  test("shows the error of a rejected password", async () => {
    sendPassword.mockResolvedValue({ error: "errors.failedToAuthenticateNoLimit" });
    const { getByTestId, findByText } = renderForm();

    type(getByTestId("username-text-input"), "john@example.com");
    type(getByTestId("password-text-input"), "wrong");
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
    fireEvent.click(getByTestId("submit-button"));

    expect(await findByText("errors.failedToAuthenticateNoLimit")).toBeTruthy();
    expect(push).not.toHaveBeenCalled();
  });

  test("continues with a passkey without a password", async () => {
    sendLoginname.mockResolvedValue({ redirect: "/passkey?loginName=pat%40example.com" });
    const { getByTestId } = renderForm({ organization: "org-1" });

    type(getByTestId("username-text-input"), "pat@example.com");
    fireEvent.click(getByTestId("passkey-button"));

    await waitFor(() => expect(push).toHaveBeenCalledWith("/passkey?loginName=pat%40example.com"));
    expect(sendLoginname).toHaveBeenCalledWith(
      expect.objectContaining({ loginName: "pat@example.com", organization: "org-1" }),
    );
    expect(sendPassword).not.toHaveBeenCalled();
  });

  test("requires the email before continuing with a passkey", async () => {
    const { getByTestId, findByText } = renderForm();

    fireEvent.click(getByTestId("passkey-button"));

    expect(await findByText("required.loginName")).toBeTruthy();
    expect(sendLoginname).not.toHaveBeenCalled();
  });

  test("hides the passkey button when passkeys are not allowed", () => {
    const settings = create(LoginSettingsSchema, { passkeysType: PasskeysType.NOT_ALLOWED });
    const { queryByTestId } = renderForm({ loginSettings: settings });

    expect(queryByTestId("passkey-button")).toBeNull();
  });

  test("sends the reset link for the typed email without a password", async () => {
    resetPassword.mockResolvedValue({});
    const { getByTestId } = renderForm({ organization: "org-1" });

    type(getByTestId("username-text-input"), "john@example.com");
    fireEvent.click(getByTestId("reset-button"));

    await waitFor(() => expect(push).toHaveBeenCalledWith("/password/set?loginName=john%40example.com&organization=org-1"));
    expect(resetPassword).toHaveBeenCalledWith(expect.objectContaining({ loginName: "john@example.com" }));
  });

  test("discovers the user when the page asks to submit the login hint", async () => {
    sendLoginname.mockResolvedValue({ redirect: "/password?loginName=john%40example.com" });
    renderForm({ loginName: "john@example.com", submit: true });

    await waitFor(() =>
      expect(sendLoginname).toHaveBeenCalledWith(expect.objectContaining({ loginName: "john@example.com" })),
    );
    expect(sendPassword).not.toHaveBeenCalled();
  });
});
