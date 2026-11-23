import { create } from "@zitadel/client";
import { RequestChallengesSchema } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { afterEach, assert, describe, expect, it, vi } from "vitest";
import { getOTPSMSResendCooldownSeconds, getOTPSMSRetryAfterSeconds, sanitizeChallenges } from "./challenges";

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

describe("getOTPSMSResendCooldownSeconds", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it("defaults to 30 seconds", () => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", undefined);
    expect(getOTPSMSResendCooldownSeconds()).toBe(30);
  });

  it("uses the configured value", () => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", "60");
    expect(getOTPSMSResendCooldownSeconds()).toBe(60);
  });

  it("allows disabling the cooldown with 0", () => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", "0");
    expect(getOTPSMSResendCooldownSeconds()).toBe(0);
  });

  it.each(["", "abc", "-5", "1.5"])("falls back to the default for the invalid value %j", (value) => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", value);
    expect(getOTPSMSResendCooldownSeconds()).toBe(30);
  });
});

describe("getOTPSMSRetryAfterSeconds", () => {
  const now = 1_000_000;

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it("allows sending when no SMS was sent yet", () => {
    expect(getOTPSMSRetryAfterSeconds(undefined, now)).toBe(0);
  });

  it("allows sending when the timestamp is not a number", () => {
    expect(getOTPSMSRetryAfterSeconds("garbage", now)).toBe(0);
  });

  it("returns the remaining seconds, rounded up, within the cooldown", () => {
    expect(getOTPSMSRetryAfterSeconds(`${now}`, now)).toBe(30);
    expect(getOTPSMSRetryAfterSeconds(`${now - 10_500}`, now)).toBe(20);
    expect(getOTPSMSRetryAfterSeconds(`${now - 29_999}`, now)).toBe(1);
  });

  it("allows sending once the cooldown has passed", () => {
    expect(getOTPSMSRetryAfterSeconds(`${now - 30_000}`, now)).toBe(0);
  });

  it("caps a timestamp in the future at the full cooldown", () => {
    expect(getOTPSMSRetryAfterSeconds(`${now + 3_600_000}`, now)).toBe(30);
  });

  it("follows the configured cooldown", () => {
    vi.stubEnv("OTP_SMS_RESEND_COOLDOWN_SECONDS", "0");
    expect(getOTPSMSRetryAfterSeconds(`${now}`, now)).toBe(0);
  });
});
