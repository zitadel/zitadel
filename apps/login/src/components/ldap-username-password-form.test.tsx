import { cleanup, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { LDAPUsernamePasswordForm } from "./ldap-username-password-form";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("@/lib/actions/idp", () => ({
  createNewSessionForLDAP: vi.fn(),
}));

describe("LDAPUsernamePasswordForm", () => {
  afterEach(cleanup);

  test("should autofocus the username input on mount", () => {
    const { getByTestId } = render(<LDAPUsernamePasswordForm idpId="idp-1" link={false} />);
    expect(getByTestId("username-text-input")).toHaveFocus();
  });

  test("should set autocomplete=current-password on the password input", () => {
    const { getByTestId } = render(<LDAPUsernamePasswordForm idpId="idp-1" link={false} />);
    expect(getByTestId("password-text-input")).toHaveAttribute("autocomplete", "current-password");
  });

  test("should render one username field before the password input", () => {
    const { container, getByTestId } = render(<LDAPUsernamePasswordForm idpId="idp-1" link={false} />);

    const usernames = container.querySelectorAll('input[autocomplete="username"]');
    expect(usernames).toHaveLength(1);

    // password managers pair the password with a preceding username field
    const password = getByTestId("password-text-input");
    expect(usernames[0].compareDocumentPosition(password) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
  });
});
