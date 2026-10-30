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

const DEFAULT_OTP_SMS_RESEND_COOLDOWN_SECONDS = 30;

/**
 * Minimum number of seconds between two OTP SMS challenges on the same session.
 *
 * Every OTP SMS challenge makes Zitadel send a text message. Without a cooldown a user (or a script)
 * can hit "resend" in a loop and flood the user's phone with codes.
 *
 * Configurable via `OTP_SMS_RESEND_COOLDOWN_SECONDS`, defaults to 30. `0` disables the cooldown;
 * anything that is not a non-negative integer falls back to the default.
 */
export function getOTPSMSResendCooldownSeconds(): number {
  const configured = process.env.OTP_SMS_RESEND_COOLDOWN_SECONDS?.trim();
  if (!configured || !/^\d+$/.test(configured)) {
    return DEFAULT_OTP_SMS_RESEND_COOLDOWN_SECONDS;
  }
  return Number(configured);
}

/**
 * Returns how many seconds are left until another OTP SMS may be sent, or 0 if one may be sent now.
 *
 * @param otpSmsSentTs the `otpSmsSentTs` (ms since epoch) of the session cookie entry
 */
export function getOTPSMSRetryAfterSeconds(otpSmsSentTs: string | undefined, now: number = Date.now()): number {
  const sentAt = Number(otpSmsSentTs);
  if (!otpSmsSentTs || !Number.isFinite(sentAt)) {
    return 0;
  }

  const remainingMs = sentAt + getOTPSMSResendCooldownSeconds() * 1000 - now;
  // a timestamp in the future can only come from clock skew between replicas, cap it at the full cooldown
  return Math.min(Math.max(Math.ceil(remainingMs / 1000), 0), getOTPSMSResendCooldownSeconds());
}
