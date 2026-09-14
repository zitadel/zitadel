import { create } from "@zitadel/client";
import { RequestChallengesSchema } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { assert, describe, expect, it } from "vitest";
import { sanitizeChallenges } from "./challenges";

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
