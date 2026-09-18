import { describe, expect, it } from "vitest";
import { loginNameEquals, normalizeLoginName } from "./login-name";

describe("normalizeLoginName", () => {
  it.each([
    ["user@example.com  ", "user@example.com"],
    ["  user@example.com", "user@example.com"],
    ["\tuser@example.com\n", "user@example.com"],
    ["user@example.com", "user@example.com"],
  ])("should strip surrounding whitespace from %j", (input, expected) => {
    expect(normalizeLoginName(input)).toBe(expected);
  });

  it("should keep the casing so it can be echoed back to the user", () => {
    expect(normalizeLoginName("  User@Example.com ")).toBe("User@Example.com");
  });

  it("should not touch inner characters, a username is not necessarily an email", () => {
    expect(normalizeLoginName(" my user ")).toBe("my user");
  });
});

describe("loginNameEquals", () => {
  it.each([
    ["user@example.com", "USER@EXAMPLE.COM"],
    ["User@Example.com", "user@example.com"],
    ["BREcloud@example.com", "brecloud@example.com"],
    ["user@example.com", "user@example.com  "],
    ["  user@example.com", "USER@example.com"],
  ])("should treat %j and %j as the same login name", (a, b) => {
    expect(loginNameEquals(a, b)).toBe(true);
  });

  it.each([
    ["user@example.com", "other@example.com"],
    ["user@example.com", "user@example.org"],
    ["user@example.com", "us er@example.com"],
  ])("should not match %j with %j", (a, b) => {
    expect(loginNameEquals(a, b)).toBe(false);
  });

  // The callers use this to check whether the input identified the user by email or
  // phone. An unset value must never count as a match, otherwise the enumeration
  // guards would accept a user the input did not actually name.
  it.each([
    [undefined, "user@example.com"],
    ["user@example.com", undefined],
    [undefined, undefined],
    ["", ""],
    ["", "user@example.com"],
  ])("should not match when a value is missing (%j, %j)", (a, b) => {
    expect(loginNameEquals(a, b)).toBe(false);
  });
});
