import { Code } from "@zitadel/client";
import { describe, expect, it } from "vitest";
import { invalidSessionTokenRedirect, isStaleSessionToken } from "./invalid-session-token";

describe("stale session token on the account button", () => {
  it("recognises the Zitadel error for a rotated session token", () => {
    expect(isStaleSessionToken(Code.PermissionDenied, "Session Token ist ungültig (COMMAND-sGr42)")).toBe(true);
    expect(isStaleSessionToken(Code.FailedPrecondition, "Session Token ist ungültig (COMMAND-sGr42)")).toBe(false);
    expect(isStaleSessionToken(Code.PermissionDenied, "User not granted")).toBe(false);
  });

  it("sends the account click back to the password form", () => {
    expect(
      invalidSessionTokenRedirect({
        loginName: "admin@example.test",
        organization: "123",
        requestId: "3927",
      }),
    ).toEqual({
      redirect: "/password?loginName=admin%40example.test&organization=123&requestId=oidc_3927",
    });
  });

  it("keeps an existing oidc prefix", () => {
    expect(invalidSessionTokenRedirect({ loginName: "a@b.test", requestId: "oidc_V2_1" })).toEqual({
      redirect: "/password?loginName=a%40b.test&requestId=oidc_V2_1",
    });
  });

  it("does not invent a login name", () => {
    expect(invalidSessionTokenRedirect({ requestId: "oidc_1" })).toEqual({ error: "Session not found or invalid" });
  });
});
