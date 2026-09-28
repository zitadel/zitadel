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
    },
    {
      name: "U2F",
      card: U2F,
      translatedText: "authenticator.methods.u2f",
    },
    {
      name: "EMAIL",
      card: EMAIL,
      translatedText: "authenticator.methods.otpEmail",
    },
    {
      name: "SMS",
      card: SMS,
      translatedText: "authenticator.methods.otpSms",
    },
    {
      name: "PASSKEYS",
      card: PASSKEYS,
      translatedText: "authenticator.methods.passkey",
    },
    {
      name: "PASSWORD",
      card: PASSWORD,
      translatedText: "authenticator.methods.password",
    },
  ])("renders $name label from translation key", ({ card, translatedText }) => {
    render(card(false, "/auth/method"));

    const translatedLabel = screen.getByText(translatedText);
    expect(translatedLabel).toBeInTheDocument();
    expect(translatedLabel).toHaveAttribute("data-i18n-key", translatedText);
  });
});
