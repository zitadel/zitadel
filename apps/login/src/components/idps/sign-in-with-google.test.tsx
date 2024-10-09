import { afterEach, describe, expect, test } from "vitest";

import { cleanup, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { SignInWithGoogle } from "./sign-in-with-google";

afterEach(cleanup);

describe("<SignInWithGoogle />", async () => {
  const messages = await getMessages({ locale: "en" });

  test("renders without crashing", () => {
    const { container } = render(
      <NextIntlClientProvider messages={messages}>
        <SignInWithGoogle />
      </NextIntlClientProvider>,
    );
    expect(container.firstChild).toBeDefined();
  });

  test("displays the default text", () => {
    render(
      <NextIntlClientProvider messages={messages}>
        <SignInWithGoogle />
      </NextIntlClientProvider>,
    );
    const signInText = screen.getByText(/Sign in with Google/i);
    expect(signInText).toBeInTheDocument();
  });

  test("displays the given text", () => {
    render(
      <NextIntlClientProvider messages={messages}>
        <SignInWithGoogle name={"Google"} />
      </NextIntlClientProvider>,
    );
    const signInText = screen.getByText(/Google/i);
    expect(signInText).toBeInTheDocument();
  });
});
