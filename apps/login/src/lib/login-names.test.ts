import { describe, expect, test } from "vitest";
import { isLoginNameOfUser } from "./login-names";

describe("isLoginNameOfUser", () => {
  const user = {
    preferredLoginName: "user@example.com",
    loginNames: ["user@example.com", "user@organization.id.example.com"],
  };

  test("accepts the preferred login name", () => {
    expect(isLoginNameOfUser(user, "user@example.com")).toBe(true);
  });

  test("accepts a login name other than the preferred one", () => {
    expect(isLoginNameOfUser(user, "user@organization.id.example.com")).toBe(true);
  });

  test("ignores case, as the login name search does", () => {
    expect(isLoginNameOfUser(user, "User@Organization.ID.example.com")).toBe(true);
  });

  test("rejects a value that is not one of the user's login names", () => {
    expect(isLoginNameOfUser(user, "user@gmail.com")).toBe(false);
  });

  test("falls back to the preferred login name when login names are absent", () => {
    const partial = { preferredLoginName: "user@example.com" } as Parameters<typeof isLoginNameOfUser>[0];
    expect(isLoginNameOfUser(partial, "user@example.com")).toBe(true);
    expect(isLoginNameOfUser(partial, "other@example.com")).toBe(false);
  });
});
