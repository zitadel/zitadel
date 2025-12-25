import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, test } from "vitest";
import { PasswordComplexity } from "./password-complexity";
import { NextIntlClientProvider } from "next-intl";

describe("<PasswordComplexity/>", () => {
  const messages = {
    password: {
      complexity: {
        length: "Must be at least {minLength} characters long.",
        hasSymbol: "Must include a symbol.",
        hasNumber: "Must include a number.",
        hasUppercase: "Must include an uppercase letter.",
        hasLowercase: "Must include a lowercase letter.",
        equals: "Password confirmation matched.",
        matches: "Matches",
        doesNotMatch: "Doesn't match",
      },
    },
  };

  afterEach(cleanup);

  test("should render length check when minLength is defined", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="Password1!"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(5),
              requiresLowercase: false,
              requiresUppercase: false,
              requiresNumber: false,
              requiresSymbol: false,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    const lengthCheck = screen.getByTestId("length-check");
    expect(lengthCheck).toBeInTheDocument();
    expect(lengthCheck.querySelector("svg")).toBeInTheDocument();
  });

  test("should not render length check when minLength is undefined", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="Password1!"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(0),
              requiresLowercase: false,
              requiresUppercase: false,
              requiresNumber: false,
              requiresSymbol: false,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    expect(screen.queryByTestId("length-check")).toBeInTheDocument();
  });

  test("should render check icon when password meets length requirement", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="Password1!"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(5),
              requiresLowercase: false,
              requiresUppercase: false,
              requiresNumber: false,
              requiresSymbol: false,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    const lengthCheck = screen.getByTestId("length-check");
    const svg = lengthCheck.querySelector("svg");
    expect(svg).toHaveClass("text-green-500");
  });

  test("should render cross icon when password does not meet length requirement", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="Pass"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(10),
              requiresLowercase: false,
              requiresUppercase: false,
              requiresNumber: false,
              requiresSymbol: false,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    const lengthCheck = screen.getByTestId("length-check");
    const svg = lengthCheck.querySelector("svg");
    expect(svg).toHaveClass("text-warn-light-500");
  });

  test("should render all complexity checks", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="Password1!"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(8),
              requiresLowercase: true,
              requiresUppercase: true,
              requiresNumber: true,
              requiresSymbol: true,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    expect(screen.getByTestId("length-check")).toBeInTheDocument();
    expect(screen.getByTestId("symbol-check")).toBeInTheDocument();
    expect(screen.getByTestId("number-check")).toBeInTheDocument();
    expect(screen.getByTestId("uppercase-check")).toBeInTheDocument();
    expect(screen.getByTestId("lowercase-check")).toBeInTheDocument();
  });

  test("should not render optional checks when policy flags are false", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="abc"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(3),
              requiresLowercase: false,
              requiresUppercase: false,
              requiresNumber: false,
              requiresSymbol: false,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    expect(screen.queryByTestId("symbol-check")).not.toBeInTheDocument();
    expect(screen.queryByTestId("number-check")).not.toBeInTheDocument();
    expect(screen.queryByTestId("uppercase-check")).not.toBeInTheDocument();
    expect(screen.queryByTestId("lowercase-check")).not.toBeInTheDocument();
  });

  test("should render an enabled rule in failing state when password does not match", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="password1!"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(8),
              requiresLowercase: true,
              requiresUppercase: true,
              requiresNumber: true,
              requiresSymbol: true,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    const uppercaseCheck = screen.getByTestId("uppercase-check");
    const uppercaseIcon = uppercaseCheck.querySelector("svg");
    expect(uppercaseIcon).toHaveClass("text-warn-light-500");
  });
});
