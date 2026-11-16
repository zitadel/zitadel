import { createIdpFlowProof, IDP_FLOW_PROOF_MAX_AGE_SECONDS, isIdpFlowProofValid } from "@/lib/server/idp-flow-proof";
import { createVerificationCheck } from "@/lib/server/verification-check";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("next/headers", () => ({ cookies: vi.fn() }));
vi.mock("@/lib/server/fingerprint", () => ({ getBindingIdCookie: vi.fn(), getOrSetBindingId: vi.fn() }));

describe("idp-flow-proof", () => {
  beforeEach(() => {
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it("accepts a proof checked with the same binding id and session", () => {
    expect(isIdpFlowProofValid(createIdpFlowProof("fp-1"), "fp-1")).toBe(true);
    expect(isIdpFlowProofValid(createIdpFlowProof("fp-1", "session-1"), "fp-1", "session-1")).toBe(true);
  });

  it("rejects a proof for another binding id or a missing value", () => {
    const value = createIdpFlowProof("fp-1");
    expect(isIdpFlowProofValid(value, "fp-2")).toBe(false);
    expect(isIdpFlowProofValid(value, undefined)).toBe(false);
    expect(isIdpFlowProofValid(undefined, "fp-1")).toBe(false);
  });

  it("rejects a sign-in proof checked with a session and a link proof checked without one", () => {
    expect(isIdpFlowProofValid(createIdpFlowProof("fp-1"), "fp-1", "session-1")).toBe(false);
    expect(isIdpFlowProofValid(createIdpFlowProof("fp-1", "session-1"), "fp-1")).toBe(false);
  });

  it("treats an empty session the same as no session", () => {
    expect(isIdpFlowProofValid(createIdpFlowProof("fp-1", ""), "fp-1", undefined)).toBe(true);
    expect(isIdpFlowProofValid(createIdpFlowProof("fp-1", undefined), "fp-1", "")).toBe(true);
  });

  it("rejects a proof for another session", () => {
    expect(isIdpFlowProofValid(createIdpFlowProof("fp-1", "session-1"), "fp-1", "session-2")).toBe(false);
  });

  it("rejects a value signed with a secret the server does not know", () => {
    const value = createIdpFlowProof("fp-1", "session-1");
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "another-secret-".padEnd(32, "x"));
    expect(isIdpFlowProofValid(value, "fp-1", "session-1")).toBe(false);
  });

  it("enforces the expiry server side", () => {
    const now = Date.now();
    const value = createIdpFlowProof("fp-1", "session-1", now);
    expect(isIdpFlowProofValid(value, "fp-1", "session-1", now + (IDP_FLOW_PROOF_MAX_AGE_SECONDS - 1) * 1000)).toBe(true);
    expect(isIdpFlowProofValid(value, "fp-1", "session-1", now + (IDP_FLOW_PROOF_MAX_AGE_SECONDS + 1) * 1000)).toBe(false);
  });

  it("is not interchangeable with a verification check for the same values", () => {
    const now = Date.now();
    // a proof of another purpose with the same two fields and expiry must not be accepted as flow proof
    const verificationCheck = createVerificationCheck(
      "fp-1",
      "session-1",
      now + (IDP_FLOW_PROOF_MAX_AGE_SECONDS - 300) * 1000,
    );
    expect(isIdpFlowProofValid(verificationCheck, "fp-1", "session-1", now)).toBe(false);
  });
});
