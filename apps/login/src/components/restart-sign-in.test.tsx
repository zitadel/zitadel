import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { RestartSignIn } from "./restart-sign-in";

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

describe("RestartSignIn", () => {
  afterEach(cleanup);

  test("links to the login name step and keeps the login name, organization and auth request", () => {
    render(<RestartSignIn loginName="user+tag@example.com" organization="org-1" requestId="oidc_123" />);

    const link = screen.getByTestId("restart-sign-in");
    expect(link.getAttribute("href")).toBe(
      "/loginname?loginName=user%2Btag%40example.com&organization=org-1&requestId=oidc_123",
    );
    expect(link.textContent).toBe("sessionResolution.restart");
  });

  test("leaves out what is not known", () => {
    render(<RestartSignIn requestId="saml_456" />);

    expect(screen.getByTestId("restart-sign-in").getAttribute("href")).toBe("/loginname?requestId=saml_456");
  });
});
