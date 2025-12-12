import { Alert } from "@/components/alert";
import { BackButton } from "@/components/back-button";
import { Button, ButtonVariants } from "@/components/button";
import { DynamicTheme } from "@/components/dynamic-theme";
import { TotpRegister } from "@/components/totp-register";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getServiceConfig } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { addOTPEmail, addOTPSMS, getBrandingSettings, getLoginSettings, getUserByID, registerTOTP } from "@/lib/zitadel";
import { RegisterTOTPResponse } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { headers } from "next/headers";
import Link from "next/link";
import { redirect } from "next/navigation";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const params = await props.params;
  const searchParams = await props.searchParams;

  const { loginName, organization, sessionId, requestId, checkAfter } = searchParams;
  const { method } = params;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const branding = await getBrandingSettings({ serviceConfig, organization,
  });
  const loginSettings = await getLoginSettings({ serviceConfig, organization,
  });

  const session = await loadMostRecentSession({ serviceConfig, sessionParams: {
      loginName,
      organization,
    },
  });

  // Get user information to check verification status
  let phoneVerified = false;
  let emailVerified = false;
  if (session?.factors?.user?.id) {
    const userResponse = await getUserByID({ serviceConfig, userId: session.factors.user.id });
    if (userResponse?.user?.type.case === "human") {
      const humanUser = userResponse.user.type.value;
      phoneVerified = humanUser.phone?.isVerified ?? false;
      emailVerified = humanUser.email?.isVerified ?? false;
    }
  }

  let totpResponse: RegisterTOTPResponse | undefined, error: Error | undefined;
  if (session && session.factors?.user?.id) {
    if (method === "time-based") {
      await registerTOTP({ serviceConfig, userId: session.factors.user.id,
      })
        .then((resp) => {
          if (resp) {
            totpResponse = resp;
          }
        })
        .catch((err) => {
          error = err;
        });
    } else if (method === "sms") {
      await addOTPSMS({ serviceConfig, userId: session.factors.user.id,
      }).catch((_error) => {
        // TODO: Throw this error?
        new Error("Could not add OTP via SMS");
      });
    } else if (method === "email") {
      await addOTPEmail({ serviceConfig, userId: session.factors.user.id,
      }).catch((_error) => {
        // TODO: Throw this error?
        new Error("Could not add OTP via Email");
      });
    } else {
      throw new Error("Invalid method");
    }
  } else {
    throw new Error("No session found");
  }

  const paramsToContinue = new URLSearchParams({});
  let urlToContinue = "/accounts";

  if (sessionId) {
    paramsToContinue.append("sessionId", sessionId);
  }
  if (loginName) {
    paramsToContinue.append("loginName", loginName);
  }
  if (organization) {
    paramsToContinue.append("organization", organization);
  }

  if (checkAfter) {
    if (requestId) {
      paramsToContinue.append("requestId", requestId);
    }

    // Check if contact method needs verification
    const needsVerification = 
      (method === "sms" && !phoneVerified) || 
      (method === "email" && !emailVerified);

    if (needsVerification) {
      // Contact method is not verified, redirect to OTP verification
      urlToContinue = `/otp/${method}?` + paramsToContinue;
      return redirect(urlToContinue);
    } else {
      // Contact is already verified, skip OTP verification and go to login flow
      if (requestId && sessionId) {
        const loginParams = new URLSearchParams();
        if (sessionId) {
          loginParams.append("sessionId", sessionId);
        }
        if (loginName) {
          loginParams.append("loginName", loginName);
        }
        if (organization) {
          loginParams.append("organization", organization);
        }
        if (requestId) {
          loginParams.append("authRequest", requestId);
        }
        urlToContinue = `/login?` + loginParams;
      } else if (loginName) {
        urlToContinue = `/signedin?` + paramsToContinue;
      }
    }
  } else if (requestId && sessionId) {
    if (requestId) {
      paramsToContinue.append("authRequest", requestId);
    }
    urlToContinue = `/login?` + paramsToContinue;
  } else if (loginName) {
    if (requestId) {
      paramsToContinue.append("requestId", requestId);
    }
    urlToContinue = `/signedin?` + paramsToContinue;
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="set.title" namespace="otp" />
        </h1>

        {totpResponse && "uri" in totpResponse && "secret" in totpResponse ? (
          <p className="ztdl-p">
            <Translated i18nKey="set.totpRegisterDescription" namespace="otp" />
          </p>
        ) : (
          <p className="ztdl-p">
            {method === "email" ? (
              <Translated i18nKey="set.emailOtpAdded" namespace="otp" />
            ) : method === "sms" ? (
              <Translated i18nKey="set.smsOtpAdded" namespace="otp" />
            ) : (
              ""
            )}
          </p>
        )}

        {!session && (
          <div className="py-4">
            <Alert>
              <Translated i18nKey="unknownContext" namespace="error" />
            </Alert>
          </div>
        )}

        {error && (
          <div className="py-4">
            <Alert>{error?.message}</Alert>
          </div>
        )}

        {session && (
          <UserAvatar
            loginName={loginName ?? session.factors?.user?.loginName}
            displayName={session.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}
      </div>

      <div className="w-full">
        {totpResponse && "uri" in totpResponse && "secret" in totpResponse ? (
          <div>
            <TotpRegister
              uri={totpResponse.uri as string}
              secret={totpResponse.secret as string}
              loginName={loginName}
              sessionId={sessionId}
              requestId={requestId}
              organization={organization}
              checkAfter={checkAfter === "true"}
              loginSettings={loginSettings}
            ></TotpRegister>
          </div>
        ) : (
          <div className="mt-8 flex w-full flex-col items-center gap-2">
            <Link href={urlToContinue} className="self-end w-full">
              <Button type="submit" className="self-end w-full" variant={ButtonVariants.Primary}>
                <Translated i18nKey="set.submit" namespace="otp" />
              </Button>
            </Link>
            <BackButton data-testid="back-button" />
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}
