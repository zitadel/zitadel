import { createHash } from "crypto";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { createIdpLinkProof, IDP_LINK_PROOF_MAX_AGE_SECONDS, isIdpLinkProofValid } from "./idp-link-proof";
import { createVerificationCheck } from "./verification-check";

vi.mock("next/headers", () => ({ cookies: vi.fn() }));
vi.mock("./fingerprint", () => ({ getFingerprintIdCookie: vi.fn(), getOrSetFingerprintId: vi.fn() }));

describe("idp-link-proof", () => {
  beforeEach(() => {
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it("accepts a proof issued for the same session and fingerprint", () => {
    const value = createIdpLinkProof("session-1", "fp-1");
    expect(isIdpLinkProofValid(value, "session-1", "fp-1")).toBe(true);
  });

  it("rejects a proof for another session or another fingerprint", () => {
    const value = createIdpLinkProof("session-1", "fp-1");
    expect(isIdpLinkProofValid(value, "session-2", "fp-1")).toBe(false);
    expect(isIdpLinkProofValid(value, "session-1", "fp-2")).toBe(false);
    expect(isIdpLinkProofValid(value, undefined, "fp-1")).toBe(false);
  });

  it("rejects the legacy, client computable sha256(sessionId + fingerprintId) value", () => {
    const forged = createHash("sha256")
      .update("session-1" + "fp-1")
      .digest("hex");
    expect(isIdpLinkProofValid(forged, "session-1", "fp-1")).toBe(false);
  });

  it("rejects a value signed with a secret the server does not know", () => {
    const value = createIdpLinkProof("session-1", "fp-1");
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "another-secret-".padEnd(32, "x"));
    expect(isIdpLinkProofValid(value, "session-1", "fp-1")).toBe(false);
  });

  it("enforces the expiry server side", () => {
    const now = Date.now();
    const value = createIdpLinkProof("session-1", "fp-1", now);
    expect(isIdpLinkProofValid(value, "session-1", "fp-1", now + (IDP_LINK_PROOF_MAX_AGE_SECONDS - 1) * 1000)).toBe(true);
    expect(isIdpLinkProofValid(value, "session-1", "fp-1", now + (IDP_LINK_PROOF_MAX_AGE_SECONDS + 1) * 1000)).toBe(false);
  });

  it("is not interchangeable with a verification check for the same values", () => {
    const now = Date.now();
    // a verification check with the same fields and the same expiry must not be accepted as link proof
    const verificationCheck = createVerificationCheck(
      "session-1",
      "fp-1",
      now + (IDP_LINK_PROOF_MAX_AGE_SECONDS - 300) * 1000,
    );
    expect(isIdpLinkProofValid(verificationCheck, "session-1", "fp-1", now)).toBe(false);
  });
});
