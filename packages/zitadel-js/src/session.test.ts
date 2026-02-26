import { describe, expect, test } from "vitest";
import { isSessionExpired, isSessionValid } from "./session.js";

describe("session utilities", () => {
  test("isSessionExpired returns true for expired session", () => {
    const expired = { expiresAt: Date.now() - 1000 };
    expect(isSessionExpired(expired)).toBe(true);
  });

  test("isSessionExpired returns false for valid session", () => {
    const valid = { expiresAt: Date.now() + 60000 };
    expect(isSessionExpired(valid)).toBe(false);
  });

  test("isSessionValid returns true for valid session", () => {
    const valid = { expiresAt: Date.now() + 60000 };
    expect(isSessionValid(valid)).toBe(true);
  });

  test("isSessionValid returns false for expired session", () => {
    const expired = { expiresAt: Date.now() - 1000 };
    expect(isSessionValid(expired)).toBe(false);
  });

  test("isSessionExpired handles ISO string dates", () => {
    const past = { expiresAt: new Date(Date.now() - 1000).toISOString() };
    expect(isSessionExpired(past)).toBe(true);

    const future = { expiresAt: new Date(Date.now() + 60000).toISOString() };
    expect(isSessionExpired(future)).toBe(false);
  });
});
