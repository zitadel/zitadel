import { cleanup, render, screen } from "@testing-library/react";
import { create } from "@zitadel/client";
import { LegalAndSupportSettingsSchema } from "@zitadel/proto/zitadel/settings/v2/legal_settings_pb";
import type { AnchorHTMLAttributes } from "react";
import { afterEach, describe, expect, test, vi } from "vitest";
import { PrivacyPolicyCheckboxes } from "./privacy-policy-checkboxes";

vi.mock("next-intl", () => ({
  useLocale: () => "fr",
  useTranslations: () => (key: string) => key,
}));

vi.mock("next/link", () => ({
  default: ({ href, children, ...props }: AnchorHTMLAttributes<HTMLAnchorElement> & { href: string }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

describe("PrivacyPolicyCheckboxes", () => {
  afterEach(cleanup);

  test("resolves language placeholders in legal links with the active locale", () => {
    const legal = create(LegalAndSupportSettingsSchema, {
      tosLink: "https://demo.com/tos-{{.Lang}}",
      privacyPolicyLink: "https://demo.com/privacy-{{.Lang}}",
      helpLink: "https://demo.com/help-{{.Lang}}",
    });

    render(<PrivacyPolicyCheckboxes legal={legal} onChange={vi.fn()} />);

    expect(screen.getByTestId("tos-link")).toHaveAttribute("href", "https://demo.com/tos-fr");
    expect(screen.getByTestId("privacy-policy-link")).toHaveAttribute("href", "https://demo.com/privacy-fr");
    expect(screen.getByTestId("help-link")).toHaveAttribute("href", "https://demo.com/help-fr");
  });
});
