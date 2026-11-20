import "server-only";

import { RESEND_COOLDOWN_SECONDS } from "@/lib/constants";
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

/** When an OTP was last requested for a session, in ms since epoch: fields of the session cookie entry. */
export type OTPSentTimestamps = {
  otpSmsSentTs?: string;
  otpEmailSentTs?: string;
};

type OTPChallenges = Pick<RequestChallenges, "otpSms" | "otpEmail">;

/**
 * The OTP challenges that send a code, each with the cookie field that records when it was last
 * requested and its resend cooldown. Every such challenge makes Zitadel send a text message or a
 * mail; without a cooldown a user (or a script) can request codes in a loop, and each new code
 * invalidates the previous one.
 *
 * The cooldowns are configured with `OTP_SMS_RESEND_COOLDOWN_SECONDS` and
 * `OTP_EMAIL_RESEND_COOLDOWN_SECONDS`, read on every call.
 */
const OTP_CHANNELS = [
  {
    challenge: "otpSms",
    sentTsField: "otpSmsSentTs",
    configuredCooldown: () => process.env.OTP_SMS_RESEND_COOLDOWN_SECONDS,
  },
  {
    challenge: "otpEmail",
    sentTsField: "otpEmailSentTs",
    configuredCooldown: () => process.env.OTP_EMAIL_RESEND_COOLDOWN_SECONDS,
  },
] as const;

// a value that is not a non-negative integer falls back to RESEND_COOLDOWN_SECONDS, `0` disables the cooldown
function parseCooldownSeconds(configured: string | undefined): number {
  const trimmed = configured?.trim();
  if (!trimmed || !/^\d+$/.test(trimmed)) {
    return RESEND_COOLDOWN_SECONDS;
  }
  return Number(trimmed);
}

function getRetryAfterSeconds(sentTs: string | undefined, cooldownSeconds: number, now: number): number {
  const sentAt = Number(sentTs);
  if (!sentTs || !Number.isFinite(sentAt)) {
    return 0;
  }

  const remainingSeconds = Math.ceil((sentAt + cooldownSeconds * 1000 - now) / 1000);
  // a timestamp in the future can only come from clock skew between replicas, cap it at the full cooldown
  return Math.min(Math.max(remainingSeconds, 0), cooldownSeconds);
}

function requestedChannels(challenges: OTPChallenges | undefined) {
  return OTP_CHANNELS.filter((channel) => !!challenges?.[channel.challenge]);
}

/**
 * Returns the resend cooldown of the OTP challenges requested, the longest if both are, or 0 if none
 * is requested or the cooldown is disabled.
 */
export function getOTPResendCooldownSeconds(challenges: OTPChallenges | undefined): number {
  return Math.max(0, ...requestedChannels(challenges).map((channel) => parseCooldownSeconds(channel.configuredCooldown())));
}

/**
 * Returns how many seconds are left until the OTP challenges requested may be sent again, or 0 if
 * they may be sent now.
 *
 * @param sentTimestamps the session cookie entry, which records when each OTP was last requested
 */
export function getOTPRetryAfterSeconds(
  challenges: OTPChallenges | undefined,
  sentTimestamps: OTPSentTimestamps,
  now: number = Date.now(),
): number {
  return Math.max(
    0,
    ...requestedChannels(challenges).map((channel) =>
      getRetryAfterSeconds(sentTimestamps[channel.sentTsField], parseCooldownSeconds(channel.configuredCooldown()), now),
    ),
  );
}

/**
 * Returns the timestamps a session cookie entry carries after an update: now for each OTP challenge
 * requested, the previous value for the others. They are carried over on every update, otherwise any
 * other update of the session would reset the resend cooldown.
 */
export function getOTPSentTimestamps(
  challenges: OTPChallenges | undefined,
  previous: OTPSentTimestamps = {},
  now: number = Date.now(),
): OTPSentTimestamps {
  const timestamps: OTPSentTimestamps = {};
  for (const channel of OTP_CHANNELS) {
    const sentTs = challenges?.[channel.challenge] ? `${now}` : previous[channel.sentTsField];
    if (sentTs) {
      timestamps[channel.sentTsField] = sentTs;
    }
  }
  return timestamps;
}

/**
 * Thrown when a session was updated, so the requested challenges were already sent, but a step after
 * the update failed. Callers must not request the same challenges again, the user already has the code.
 */
export class SessionUpdatedError extends Error {
  constructor(cause: unknown) {
    super("session was updated, but the steps after the update failed", { cause });
    this.name = "SessionUpdatedError";
  }
}
