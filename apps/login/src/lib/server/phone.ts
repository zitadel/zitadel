"use server";

import { completeFlowOrGetUrl } from "@/lib/client";
import { addOTPSMS, getLoginSettings, resendPhoneCode, setPhone, verifyPhone } from "@/lib/zitadel";
import { create } from "@zitadel/client";
import { RequestChallengesSchema } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { headers } from "next/headers";
import { updateOrCreateSession } from "./session";
import { getServiceConfig } from "../service-url";

type PhoneFlowCommand = {
  userId: string;
  loginName?: string;
  sessionId?: string;
  requestId?: string;
  organization?: string;
  checkAfter?: string;
  send?: string;
};

function buildCommonParams(command: PhoneFlowCommand) {
  const params = new URLSearchParams({
    userId: command.userId,
  });

  if (command.loginName) {
    params.set("loginName", command.loginName);
  }
  if (command.sessionId) {
    params.set("sessionId", command.sessionId);
  }
  if (command.requestId) {
    params.set("requestId", command.requestId);
  }
  if (command.organization) {
    params.set("organization", command.organization);
  }
  if (command.checkAfter) {
    params.set("checkAfter", command.checkAfter);
  }
  if (command.send) {
    params.set("send", command.send);
  }

  return params;
}

export async function setPhoneAndContinue(command: PhoneFlowCommand & { phone: string }) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  return setPhone({
    serviceConfig,
    userId: command.userId,
    phone: command.phone,
  })
    .then(() => {
      const params = buildCommonParams(command);
      return { redirect: `/phone/verify?${params}` };
    })
    .catch((error) => {
      console.warn(error);
      return { error: "Could not save phone number" };
    });
}

export async function resendPhoneVerification(command: { userId: string }) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  return resendPhoneCode({
    serviceConfig,
    userId: command.userId,
  }).catch((error) => {
    console.warn(error);
    return { error: "Could not resend SMS code" };
  });
}

export async function verifyPhoneAndContinue(command: PhoneFlowCommand & { code: string }) {
  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  return verifyPhone({
    serviceConfig,
    userId: command.userId,
    verificationCode: command.code,
  })
    .then(async () => {
      // During MFA setup with checkAfter=true, auto-setup + auto-check OTP SMS
      // right after phone verification to avoid asking for two consecutive codes.
      if (command.checkAfter === "true") {
        const addOtpResponse = await addOTPSMS({
          serviceConfig,
          userId: command.userId,
        }).catch((error) => {
          console.warn(error);
          return { error: "Could not add OTP via SMS" };
        });

        if (addOtpResponse && "error" in addOtpResponse && addOtpResponse.error) {
          return { error: "Could not add OTP via SMS" };
        }

        const challenges = create(RequestChallengesSchema, {
          otpSms: { returnCode: true },
        });

        const challengeResponse = await updateOrCreateSession({
          loginName: command.loginName,
          sessionId: command.sessionId,
          organization: command.organization,
          requestId: command.requestId,
          challenges,
        }).catch((error) => {
          console.warn(error);
          return { error: "Could not request OTP challenge" };
        });

        if (challengeResponse && "error" in challengeResponse && challengeResponse.error) {
          return { error: challengeResponse.error };
        }

        const otpCode =
          challengeResponse && "challenges" in challengeResponse ? challengeResponse?.challenges?.otpSms : undefined;

        if (!otpCode) {
          return { error: "Could not request OTP challenge" };
        }

        const checks = create(ChecksSchema, {
          otpSms: { code: otpCode },
        });

        const checkResponse = await updateOrCreateSession({
          loginName: command.loginName,
          sessionId: command.sessionId,
          organization: command.organization,
          requestId: command.requestId,
          checks,
        }).catch((error) => {
          console.warn(error);
          return { error: "Could not verify OTP code" };
        });

        if (checkResponse && "error" in checkResponse && checkResponse.error) {
          return { error: checkResponse.error };
        }

        if (!checkResponse || !("sessionId" in checkResponse) || !checkResponse?.factors?.user) {
          return { error: "Could not continue session" };
        }

        // Keep consistency with the existing OTP flow completion behavior.
        await new Promise((resolve) => setTimeout(resolve, 2000));

        const loginSettings = await getLoginSettings({
          serviceConfig,
          organization: checkResponse.factors.user.organizationId ?? command.organization,
        });

        if (command.requestId && checkResponse.sessionId) {
          return completeFlowOrGetUrl(
            {
              sessionId: checkResponse.sessionId,
              requestId: command.requestId,
              organization: checkResponse.factors.user.organizationId ?? command.organization,
            },
            loginSettings?.defaultRedirectUri,
          );
        }

        return completeFlowOrGetUrl(
          {
            loginName: checkResponse.factors.user.loginName,
            organization: checkResponse.factors.user.organizationId ?? command.organization,
          },
          loginSettings?.defaultRedirectUri,
        );
      }

      const params = buildCommonParams(command);
      return { redirect: `/otp/sms/set?${params}` };
    })
    .catch((error) => {
      console.warn(error);
      return { error: "Could not verify phone number" };
    });
}
