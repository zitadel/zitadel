import { create } from "@zitadel/client";
import { RequestChallengesSchema } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { afterEach, assert, describe, expect, it, vi } from "vitest";
import {
  getOTPResendCooldownSeconds,
  getOTPRetryAfterSeconds,
  getOTPSentTimestamps,
  sanitizeChallenges,
} from "./challenges";

describe("sanitizeChallenges", () => {
  it("returns undefined when no challenges are requested", () => {
    expect(sanitizeChallenges(undefined)).toBeUndefined();
  });

  it("forces returnCode off for an OTP SMS challenge", () => {
    const challenges = create(RequestChallengesSchema, {
      otpSms: { returnCode: true },
    });

    expect(sanitizeChallenges(challenges)?.otpSms?.returnCode).toBe(false);
  });

  it("leaves an OTP SMS challenge that already sends the code alone", () => {
    const challenges = create(RequestChallengesSchema, {
      otpSms: {},
    });

    expect(sanitizeChallenges(challenges)?.otpSms?.returnCode).toBe(false);
  });

  it("rewrites an OTP email returnCode challenge to sendCode", () => {
    const challenges = create(RequestChallengesSchema, {
      otpEmail: { deliveryType: { case: "returnCode", value: {} } },
    });

    const deliveryType = sanitizeChallenges(challenges)?.otpEmail?.deliveryType;

    assert(deliveryType?.case === "sendCode");
    // no url template: Zitadel falls back to its default verification url
    expect(deliveryType.value.urlTemplate).toBeUndefined();
  });

  it("preserves the url template of an OTP email sendCode challenge", () => {
    const urlTemplate = "https://example.com/otp/verify?code={{.Code}}";
    const challenges = create(RequestChallengesSchema, {
      otpEmail: { deliveryType: { case: "sendCode", value: { urlTemplate } } },
    });

    const deliveryType = sanitizeChallenges(challenges)?.otpEmail?.deliveryType;

    assert(deliveryType?.case === "sendCode");
    expect(deliveryType.value.urlTemplate).toBe(urlTemplate);
  });

  it("passes a WebAuthN challenge through untouched", () => {
    const challenges = create(RequestChallengesSchema, {
      webAuthN: { domain: "login.example.com", userVerificationRequirement: 1 },
    });

    expect(sanitizeChallenges(challenges)?.webAuthN).toEqual(challenges.webAuthN);
  });

  it("strips returnCode from both OTP methods requested together", () => {
    const challenges = create(RequestChallengesSchema, {
      otpSms: { returnCode: true },
      otpEmail: { deliveryType: { case: "returnCode", value: {} } },
    });

    const sanitized = sanitizeChallenges(challenges);

    expect(sanitized?.otpSms?.returnCode).toBe(false);
    expect(sanitized?.otpEmail?.deliveryType?.case).toBe("sendCode");
  });

  it("does not mutate the caller's challenges", () => {
    const challenges = create(RequestChallengesSchema, {
      otpSms: { returnCode: true },
      otpEmail: { deliveryType: { case: "returnCode", value: {} } },
    });

    sanitizeChallenges(challenges);

    expect(challenges.otpSms?.returnCode).toBe(true);
    expect(challenges.otpEmail?.deliveryType?.case).toBe("returnCode");
  });
});

const sms = { otpSms: {} } as any;
const email = { otpEmail: {} } as any;
const both = { otpSms: {}, otpEmail: {} } as any;

describe("getOTPResendCooldownSeconds", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it("defaults to 30 seconds for each channel", () => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", undefined);
    vi.stubEnv("OTP_EMAIL_RESEND_COOLDOWN_SECONDS", undefined);
    expect(getOTPResendCooldownSeconds(sms)).toBe(30);
    expect(getOTPResendCooldownSeconds(email)).toBe(30);
  });

  it("uses each channel's configured value", () => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", "60");
    vi.stubEnv("OTP_EMAIL_RESEND_COOLDOWN_SECONDS", "45");
    expect(getOTPResendCooldownSeconds(sms)).toBe(60);
    expect(getOTPResendCooldownSeconds(email)).toBe(45);
  });

  it("allows disabling the cooldown with 0", () => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", "0");
    expect(getOTPResendCooldownSeconds(sms)).toBe(0);
  });

  it.each(["", "abc", "-5", "1.5"])("falls back to the default for the invalid value %j", (value) => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", value);
    vi.stubEnv("OTP_EMAIL_RESEND_COOLDOWN_SECONDS", value);
    expect(getOTPResendCooldownSeconds(sms)).toBe(30);
    expect(getOTPResendCooldownSeconds(email)).toBe(30);
  });

  it("is 0 when no OTP challenge is requested", () => {
    expect(getOTPResendCooldownSeconds(undefined)).toBe(0);
    expect(getOTPResendCooldownSeconds({ webAuthN: {} } as any)).toBe(0);
  });

  it("takes the longer cooldown when both channels are requested", () => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", "10");
    vi.stubEnv("OTP_EMAIL_RESEND_COOLDOWN_SECONDS", "45");
    expect(getOTPResendCooldownSeconds(both)).toBe(45);
  });
});

describe("getOTPRetryAfterSeconds", () => {
  const now = 1_000_000;

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it("allows sending when no code was sent yet", () => {
    expect(getOTPRetryAfterSeconds(sms, {}, now)).toBe(0);
    expect(getOTPRetryAfterSeconds(email, {}, now)).toBe(0);
  });

  it("allows sending when the timestamp is not a number", () => {
    expect(getOTPRetryAfterSeconds(sms, { otpSmsSentTs: "garbage" }, now)).toBe(0);
  });

  it("returns the remaining seconds, rounded up, within the cooldown", () => {
    expect(getOTPRetryAfterSeconds(sms, { otpSmsSentTs: `${now}` }, now)).toBe(30);
    expect(getOTPRetryAfterSeconds(sms, { otpSmsSentTs: `${now - 10_500}` }, now)).toBe(20);
    expect(getOTPRetryAfterSeconds(email, { otpEmailSentTs: `${now - 29_999}` }, now)).toBe(1);
  });

  it("allows sending once the cooldown has passed", () => {
    expect(getOTPRetryAfterSeconds(sms, { otpSmsSentTs: `${now - 30_000}` }, now)).toBe(0);
  });

  it("caps a timestamp in the future at the full cooldown", () => {
    expect(getOTPRetryAfterSeconds(sms, { otpSmsSentTs: `${now + 3_600_000}` }, now)).toBe(30);
  });

  it("follows the configured cooldown", () => {
    vi.stubEnv("OTP_EMAIL_RESEND_COOLDOWN_SECONDS", "0");
    expect(getOTPRetryAfterSeconds(email, { otpEmailSentTs: `${now}` }, now)).toBe(0);
  });

  it("only looks at the channels requested", () => {
    expect(getOTPRetryAfterSeconds(email, { otpSmsSentTs: `${now}` }, now)).toBe(0);
    expect(getOTPRetryAfterSeconds(undefined, { otpSmsSentTs: `${now}`, otpEmailSentTs: `${now}` }, now)).toBe(0);
  });
});

describe("getOTPSentTimestamps", () => {
  const now = 1_000_000;

  it("stamps the channels requested and carries the others over", () => {
    expect(getOTPSentTimestamps(email, { otpSmsSentTs: "999000", otpEmailSentTs: "998000" }, now)).toEqual({
      otpSmsSentTs: "999000",
      otpEmailSentTs: "1000000",
    });
  });

  it("stamps both channels when both are requested", () => {
    expect(getOTPSentTimestamps(both, {}, now)).toEqual({ otpSmsSentTs: "1000000", otpEmailSentTs: "1000000" });
  });

  it("adds no field for a channel that was never requested", () => {
    expect(getOTPSentTimestamps(undefined, {}, now)).toEqual({});
  });
});
