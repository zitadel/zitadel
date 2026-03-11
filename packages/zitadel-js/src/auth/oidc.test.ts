import { describe, expect, test } from "vitest";
import {
  createOIDCAuthorizationUrl as createOIDCAuthorizationUrlFromAuthLane,
  discoverOIDCAuthorizationServer as discoverOIDCAuthorizationServerFromAuthLane,
  exchangeOIDCAuthorizationCode as exchangeOIDCAuthorizationCodeFromAuthLane,
  refreshOIDCTokens as refreshOIDCTokensFromAuthLane,
  createOIDCEndSessionUrl as createOIDCEndSessionUrlFromAuthLane,
  generatePKCE as generatePKCEFromAuthLane,
  generateState as generateStateFromAuthLane,
} from "./oidc.js";
import {
  createOIDCAuthorizationUrl,
  discoverOIDCAuthorizationServer,
  exchangeOIDCAuthorizationCode,
  refreshOIDCTokens,
  createOIDCEndSessionUrl,
} from "../oidc.js";
import { generatePKCE, generateState } from "../pkce.js";

describe("auth/oidc exports", () => {
  test("re-export oidc and pkce helpers", () => {
    expect(discoverOIDCAuthorizationServerFromAuthLane).toBe(
      discoverOIDCAuthorizationServer,
    );
    expect(createOIDCAuthorizationUrlFromAuthLane).toBe(createOIDCAuthorizationUrl);
    expect(exchangeOIDCAuthorizationCodeFromAuthLane).toBe(
      exchangeOIDCAuthorizationCode,
    );
    expect(refreshOIDCTokensFromAuthLane).toBe(refreshOIDCTokens);
    expect(createOIDCEndSessionUrlFromAuthLane).toBe(createOIDCEndSessionUrl);
    expect(generatePKCEFromAuthLane).toBe(generatePKCE);
    expect(generateStateFromAuthLane).toBe(generateState);
  });
});
