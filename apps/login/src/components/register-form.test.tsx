import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { create } from "@zitadel/client";
import { LegalAndSupportSettingsSchema } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import { PasswordComplexitySettingsSchema } from "@zitadel/proto/zitadel/settings/v2/password_settings_pb";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { RegisterForm } from "./register-form";

const push = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push, back: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
  useLocale: () => "en",
}));

vi.mock("@/lib/actions/register", () => ({
  registerUser: vi.fn(),
}));

const defaultLegal = create(LegalAndSupportSettingsSchema, {});

const complexity = create(PasswordComplexitySettingsSchema, {
  minLength: BigInt(8),
  requiresUppercase: true,
  requiresLowercase: true,
  requiresNumber: true,
  requiresSymbol: false,
});

function renderForm(props: Partial<Parameters<typeof RegisterForm>[0]> = {}) {
  return render(
    <RegisterForm legal={defaultLegal} passwordComplexitySettings={complexity} organization="org-1" {...props} />,
  );
}

function type(input: HTMLElement, value: string) {
  fireEvent.change(input, { target: { value } });
}

function fill(getByTestId: (id: string) => HTMLElement, password: string, confirmPassword = password) {
  type(getByTestId("firstname-text-input"), "Jane");
  type(getByTestId("lastname-text-input"), "Doe");
  type(getByTestId("email-text-input"), "jane@example.com");
  type(getByTestId("password-text-input"), password);
  type(getByTestId("password-confirm-text-input"), confirmPassword);
}

describe("RegisterForm", () => {
  let registerUser: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    registerUser = vi.mocked((await import("@/lib/actions/register")).registerUser);
  });

  afterEach(cleanup);

  test("should autofocus the firstname input on mount", () => {
    const { getByTestId } = renderForm();
    expect(getByTestId("firstname-text-input")).toHaveFocus();
  });

  test("asks for the name, the email and the password on one screen without a method choice", () => {
    const { getByTestId, queryByTestId } = renderForm();

    expect(getByTestId("email-text-input")).toBeTruthy();
    expect(getByTestId("password-text-input")).toBeTruthy();
    expect(getByTestId("password-confirm-text-input")).toBeTruthy();
    expect(queryByTestId("passkey-radio")).toBeNull();
    expect(queryByTestId("password-radio")).toBeNull();
  });

  test("prefills the values passed in the link", () => {
    const { getByTestId } = renderForm({ email: "jane@example.com", firstname: "Jane", lastname: "Doe" });

    expect(getByTestId("email-text-input")).toHaveValue("jane@example.com");
    expect(getByTestId("firstname-text-input")).toHaveValue("Jane");
    expect(getByTestId("lastname-text-input")).toHaveValue("Doe");
  });

  test("registers the user with the password in one step", async () => {
    registerUser.mockResolvedValue({ redirect: "/signedin?loginName=jane%40example.com" });
    const { getByTestId } = renderForm({ requestId: "oidc_1" });

    fill(getByTestId, "Password1");
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
    fireEvent.click(getByTestId("submit-button"));

    await waitFor(() => expect(push).toHaveBeenCalledWith("/signedin?loginName=jane%40example.com"));
    expect(registerUser).toHaveBeenCalledWith({
      email: "jane@example.com",
      firstName: "Jane",
      lastName: "Doe",
      password: "Password1",
      organization: "org-1",
      requestId: "oidc_1",
    });
  });

  test("disables the submit while the password misses the complexity policy", async () => {
    const { getByTestId } = renderForm();

    fill(getByTestId, "Password1");
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());

    type(getByTestId("password-text-input"), "password");
    type(getByTestId("password-confirm-text-input"), "password");
    await waitFor(() => expect(getByTestId("submit-button")).toBeDisabled());
  });

  test("disables the submit while the confirmation does not match", async () => {
    const { getByTestId } = renderForm();

    fill(getByTestId, "Password1");
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());

    type(getByTestId("password-confirm-text-input"), "Password2");
    await waitFor(() => expect(getByTestId("submit-button")).toBeDisabled());
  });

  test("requires accepting the terms when the instance links them", async () => {
    const legal = create(LegalAndSupportSettingsSchema, { tosLink: "https://example.com/tos" });
    const { getByTestId } = renderForm({ legal });

    fill(getByTestId, "Password1");
    // the fields are valid, so only the unchecked terms keep the submit disabled
    await waitFor(() => expect(getByTestId("password-confirm-text-input")).toHaveValue("Password1"));
    expect(getByTestId("submit-button")).toBeDisabled();

    fireEvent.click(getByTestId("tos-checkbox"));
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
  });

  test("shows the error of a failed registration", async () => {
    registerUser.mockResolvedValue({ error: "errors.couldNotCreateUser" });
    const { getByTestId, findByText } = renderForm();

    fill(getByTestId, "Password1");
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
    fireEvent.click(getByTestId("submit-button"));

    expect(await findByText("errors.couldNotCreateUser")).toBeTruthy();
    expect(push).not.toHaveBeenCalled();
  });
});
