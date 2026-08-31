import { cleanup, render, screen } from "@testing-library/react";
import type { AnchorHTMLAttributes } from "react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { EMAIL, PASSKEYS, PASSWORD, SMS, TOTP, U2F } from "./auth-methods";

vi.mock("next-intl", () => ({
  useTranslations: (namespace?: string) => (key: string) => `${namespace ? `${namespace}.` : ""}${key}`,
}));

vi.mock("next/link", () => ({
  default: ({ href, children, ...props }: AnchorHTMLAttributes<HTMLAnchorElement> & { href: string }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

describe("auth-methods", () => {
  afterEach(cleanup);

  test.each([
    {
      name: "TOTP",
      card: TOTP,
      translatedText: "authenticator.methods.totp",
      previousLiteral: "Authenticator App",
    },
    {
      name: "U2F",
      card: U2F,
      translatedText: "authenticator.methods.u2f",
      previousLiteral: "Universal Second Factor",
    },
    {
      name: "EMAIL",
      card: EMAIL,
      translatedText: "authenticator.methods.otpEmail",
      previousLiteral: "Code via Email",
    },
    {
      name: "SMS",
      card: SMS,
      translatedText: "authenticator.methods.otpSms",
      previousLiteral: "Code via SMS",
    },
    {
      name: "PASSKEYS",
      card: PASSKEYS,
      translatedText: "authenticator.methods.passkey",
      previousLiteral: "Passkeys",
    },
    {
      name: "PASSWORD",
      card: PASSWORD,
      translatedText: "authenticator.methods.password",
      previousLiteral: "Password",
    },
  ])("renders $name label from translation key", ({ card, translatedText, previousLiteral }) => {
    render(card(false, "/auth/method"));

    const translatedLabel = screen.getByText(translatedText);
    expect(translatedLabel).toBeInTheDocument();
    expect(translatedLabel).toHaveAttribute("data-i18n-key", translatedText);
    expect(screen.queryByText(previousLiteral)).not.toBeInTheDocument();
  });
});
