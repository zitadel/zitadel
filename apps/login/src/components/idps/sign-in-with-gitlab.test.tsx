import { afterEach, describe, expect, test } from "vitest";

import { cleanup, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { SignInWithGitlab } from "./sign-in-with-gitlab";

afterEach(cleanup);

describe("<SignInWithGitlab />", async () => {
  const messages = await getMessages({ locale: "en" });

  test("renders without crashing", () => {
    const { container } = render(<SignInWithGitlab />);
    expect(container.firstChild).toBeDefined();
  });

  test("displays the default text", () => {
    render(
      <NextIntlClientProvider messages={messages}>
        <SignInWithGitlab />
      </NextIntlClientProvider>,
    );
    const signInText = screen.getByText(/Sign in with Gitlab/i);
    expect(signInText).toBeInTheDocument();
  });

  test("displays the given text", () => {
    render(
      <NextIntlClientProvider messages={messages}>
        <SignInWithGitlab name={"Gitlab"} />
      </NextIntlClientProvider>,
    );
    const signInText = screen.getByText(/Gitlab/i);
    expect(signInText).toBeInTheDocument();
  });
});
