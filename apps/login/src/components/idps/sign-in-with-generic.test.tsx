import { afterEach, describe, expect, test } from "vitest";

import { cleanup, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { SignInWithGeneric } from "./sign-in-with-generic";

afterEach(cleanup);

describe("<SignInWithGeneric />", () => {
  const messages = {};

  function renderButton(name: string) {
    return render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <SignInWithGeneric name={name} />
      </NextIntlClientProvider>,
    );
  }

  test("renders a button labelled with the provider name", () => {
    renderButton("Acme SSO");
    expect(screen.getByRole("button", { name: /Acme SSO/i })).toBeInTheDocument();
  });

  test("renders the brand logo for a recognised provider", () => {
    const { container } = renderButton("Facebook");
    // Accessible name comes from the visible text; the glyph is decorative.
    expect(screen.getByRole("button", { name: /Facebook/i })).toBeInTheDocument();
    const svg = container.querySelector("svg");
    expect(svg).toBeInTheDocument();
    expect(svg).toHaveAttribute("aria-hidden", "true");
  });

  test("falls back to a name-only button for an unknown provider", () => {
    const { container } = renderButton("Some Corporate IdP");
    expect(screen.getByRole("button", { name: /Some Corporate IdP/i })).toBeInTheDocument();
    expect(container.querySelector("svg")).not.toBeInTheDocument();
  });

  test("does not match a brand on an unrelated substring", () => {
    // "Metadata SSO" contains "meta" but must not render the Facebook glyph.
    const { container } = renderButton("Metadata SSO");
    expect(container.querySelector("svg")).not.toBeInTheDocument();
  });
});
