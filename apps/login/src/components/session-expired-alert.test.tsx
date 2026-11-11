import { cleanup, render } from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { SessionExpiredAlert } from "./session-expired-alert";

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

vi.mock("next/link", () => ({
  default: ({ href, children, ...props }: { href: string; children: React.ReactNode }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

describe("SessionExpiredAlert", () => {
  afterEach(cleanup);

  test("links back to the login name page with the auth request kept", () => {
    const { getByTestId, getByText } = render(
      <SessionExpiredAlert hasContext loginName="ana@acme.test" requestId="oidc_V2_1" organization="org1" />,
    );
    expect(getByText("sessionExpired")).toBeInTheDocument();
    expect(getByTestId("sign-in-again-link")).toHaveAttribute(
      "href",
      "/loginname?loginName=ana%40acme.test&requestId=oidc_V2_1&organization=org1",
    );
  });

  test("falls back to a bare login name page without context", () => {
    const { getByTestId, getByText } = render(<SessionExpiredAlert hasContext={false} />);
    expect(getByText("unknownContext")).toBeInTheDocument();
    expect(getByTestId("sign-in-again-link")).toHaveAttribute("href", "/loginname");
  });
});
