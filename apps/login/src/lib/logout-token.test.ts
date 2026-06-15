import { describe, expect, test } from "vitest";
import { getLogoutTokenVerificationOptions } from "./logout-token";

describe("getLogoutTokenVerificationOptions", () => {
  test("ignores publicHost-only fallback from request headers", () => {
    const options = getLogoutTokenVerificationOptions(
      {
        baseUrl: "http://zitadel.ludocare.local:8080",
        publicHost: "login.ludocare.local:3021",
      },
      undefined,
    );

    expect(options).toBeUndefined();
  });

  test("uses explicit custom host headers and applies service port", () => {
    const options = getLogoutTokenVerificationOptions(
      {
        baseUrl: "http://zitadel.ludocare.local:8080",
        publicHost: "login.ludocare.local:3021",
      },
      "x-zitadel-instance-host:zitadel.ludocare.local,x-zitadel-public-host:zitadel.ludocare.local",
    );

    expect(options).toEqual({
      instanceHost: "zitadel.ludocare.local:8080",
      publicHost: "zitadel.ludocare.local:8080",
    });
  });

  test("keeps existing instance/public host when present", () => {
    const options = getLogoutTokenVerificationOptions(
      {
        baseUrl: "https://api.zitadel.cloud",
        instanceHost: "customer.zitadel.cloud",
        publicHost: "customer.zitadel.cloud",
      },
      undefined,
    );

    expect(options).toEqual({
      instanceHost: "customer.zitadel.cloud",
      publicHost: "customer.zitadel.cloud",
    });
  });
});
