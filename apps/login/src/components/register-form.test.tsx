import { cleanup, fireEvent, render, waitFor } from "@testing-library/react";
import { create } from "@zitadel/client";
import { LegalAndSupportSettingsSchema } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import { LoginSettingsSchema, PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { afterEach, describe, expect, test, vi } from "vitest";
import { RegisterForm } from "./register-form";

const push = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/server/register", () => ({
  registerUser: vi.fn(),
}));

vi.mock("@/lib/client", () => ({
  handleServerActionResponse: vi.fn(),
}));

const defaultLegal = create(LegalAndSupportSettingsSchema, {});

describe("RegisterForm", () => {
  afterEach(() => {
    cleanup();
    push.mockClear();
    window.history.replaceState(null, "", "/");
  });

  test("should autofocus the firstname input on mount", () => {
    const { getByTestId } = render(<RegisterForm legal={defaultLegal} organization="org-1" idpCount={0} />);
    expect(getByTestId("firstname-text-input")).toHaveFocus();
  });

  test("should prefill the form from the URL when it is restored without props", async () => {
    window.history.replaceState(
      null,
      "",
      "/register?firstname=Ada&lastname=Lovelace&email=ada%40example.com&organization=org-1",
    );
    const { getByTestId } = render(<RegisterForm legal={defaultLegal} organization="org-1" idpCount={0} />);

    await waitFor(() => expect(getByTestId("email-text-input")).toHaveValue("ada@example.com"));
    expect(getByTestId("firstname-text-input")).toHaveValue("Ada");
    expect(getByTestId("lastname-text-input")).toHaveValue("Lovelace");
  });

  test("should keep the entered values in the current history entry before going to the password page", async () => {
    window.history.replaceState(null, "", "/register?organization=org-1");
    const passwordOnly = create(LoginSettingsSchema, {
      allowLocalAuthentication: true,
      passkeysType: PasskeysType.NOT_ALLOWED,
    });
    const { getByTestId } = render(
      <RegisterForm legal={defaultLegal} organization="org-1" idpCount={0} loginSettings={passwordOnly} />,
    );

    fireEvent.change(getByTestId("firstname-text-input"), { target: { value: "Ada" } });
    fireEvent.change(getByTestId("lastname-text-input"), { target: { value: "Lovelace" } });
    fireEvent.change(getByTestId("email-text-input"), { target: { value: "ada@example.com" } });
    await waitFor(() => expect(getByTestId("submit-button")).not.toBeDisabled());
    fireEvent.click(getByTestId("submit-button"));

    await waitFor(() => expect(push).toHaveBeenCalledTimes(1));
    const expected = { firstname: "Ada", lastname: "Lovelace", email: "ada@example.com", organization: "org-1" };
    const [pushed] = push.mock.calls[0];
    expect(pushed.split("?")[0]).toBe("/register/password");
    expect(Object.fromEntries(new URLSearchParams(pushed.split("?")[1]))).toEqual(expected);
    // the entry the back button returns to now carries the values the register page prefills from
    expect(window.location.pathname).toBe("/register");
    expect(Object.fromEntries(new URLSearchParams(window.location.search))).toEqual(expected);
  });
});
