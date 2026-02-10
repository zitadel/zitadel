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

  test("should render all complexity checks when all requirements are enabled", () => {
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
    expect(screen.getByTestId("equal-check")).toBeInTheDocument();
  });

  test("should not render symbol check when requiresSymbol is false", () => {
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
              requiresSymbol: false,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    expect(screen.queryByTestId("symbol-check")).not.toBeInTheDocument();
    expect(screen.getByTestId("number-check")).toBeInTheDocument();
    expect(screen.getByTestId("uppercase-check")).toBeInTheDocument();
    expect(screen.getByTestId("lowercase-check")).toBeInTheDocument();
  });

  test("should not render number check when requiresNumber is false", () => {
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
              requiresNumber: false,
              requiresSymbol: true,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    expect(screen.getByTestId("symbol-check")).toBeInTheDocument();
    expect(screen.queryByTestId("number-check")).not.toBeInTheDocument();
    expect(screen.getByTestId("uppercase-check")).toBeInTheDocument();
    expect(screen.getByTestId("lowercase-check")).toBeInTheDocument();
  });

  test("should not render uppercase check when requiresUppercase is false", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="Password1!"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(8),
              requiresLowercase: true,
              requiresUppercase: false,
              requiresNumber: true,
              requiresSymbol: true,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    expect(screen.getByTestId("symbol-check")).toBeInTheDocument();
    expect(screen.getByTestId("number-check")).toBeInTheDocument();
    expect(screen.queryByTestId("uppercase-check")).not.toBeInTheDocument();
    expect(screen.getByTestId("lowercase-check")).toBeInTheDocument();
  });

  test("should not render lowercase check when requiresLowercase is false", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="Password1!"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(8),
              requiresLowercase: false,
              requiresUppercase: true,
              requiresNumber: true,
              requiresSymbol: true,
              resourceOwnerType: 0,
            } as any
          }
        />
      </NextIntlClientProvider>,
    );

    expect(screen.getByTestId("symbol-check")).toBeInTheDocument();
    expect(screen.getByTestId("number-check")).toBeInTheDocument();
    expect(screen.getByTestId("uppercase-check")).toBeInTheDocument();
    expect(screen.queryByTestId("lowercase-check")).not.toBeInTheDocument();
  });

  test("should only render length and equals checks when all other requirements are disabled", () => {
    render(
      <NextIntlClientProvider locale="en" messages={messages}>
        <PasswordComplexity
          password="password"
          equals
          passwordComplexitySettings={
            {
              minLength: BigInt(8),
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

    expect(screen.getByTestId("length-check")).toBeInTheDocument();
    expect(screen.getByTestId("equal-check")).toBeInTheDocument();
    expect(screen.queryByTestId("symbol-check")).not.toBeInTheDocument();
    expect(screen.queryByTestId("number-check")).not.toBeInTheDocument();
    expect(screen.queryByTestId("uppercase-check")).not.toBeInTheDocument();
    expect(screen.queryByTestId("lowercase-check")).not.toBeInTheDocument();
  });
});
