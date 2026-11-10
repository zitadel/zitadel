import { create } from "@zitadel/client";
import {
  RequestChallenges,
  RequestChallenges_OTPEmail_SendCodeSchema,
} from "@zitadel/proto/zitadel/session/v2/challenge_pb";

/**
 * Strips the `returnCode` delivery type from a client-supplied challenge request.
 *
 * `returnCode` makes Zitadel return the freshly generated OTP in plaintext in the response
 * body instead of sending it to the user's mailbox or phone. It is a legitimate API feature
 * for callers that deliver the code themselves, but this app proxies requests on behalf of
 * an untrusted browser: anything it forwards is attacker-controllable, and anything it
 * returns is attacker-readable.
 *
 * Without this, a party who knows only a victim's login name could request both OTP-Email
 * and OTP-SMS challenges on an identify-only session, read both codes straight out of the
 * server-action response and submit them back as two verified factors — GHSA-3gwm-5wx8-4gm6.
 *
 * OTP challenges are rewritten to their out-of-band equivalent rather than dropped, so the
 * legitimate 2FA flow still completes and the code reaches the actual user. WebAuthN
 * challenges carry no code and are passed through untouched.
 */
export function sanitizeChallenges(challenges?: RequestChallenges): RequestChallenges | undefined {
  if (!challenges) {
    return challenges;
  }

  const sanitized = { ...challenges };

  if (sanitized.otpSms) {
    sanitized.otpSms = { ...sanitized.otpSms, returnCode: false };
  }

  if (sanitized.otpEmail?.deliveryType?.case === "returnCode") {
    sanitized.otpEmail = {
      ...sanitized.otpEmail,
      // fall back to the default Zitadel verification url
      deliveryType: { case: "sendCode", value: create(RequestChallenges_OTPEmail_SendCodeSchema, {}) },
    };
  }

  return sanitized;
}
