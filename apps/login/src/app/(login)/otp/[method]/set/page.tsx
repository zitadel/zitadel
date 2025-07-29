import { Alert } from "@/components/alert";
import { BackButton } from "@/components/back-button";
import { Button, ButtonVariants } from "@/components/button";
import { DynamicTheme } from "@/components/dynamic-theme";
import { TotpRegister } from "@/components/totp-register";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import {
  addOTPEmail,
  addOTPSMS,
  getBrandingSettings,
  getLoginSettings,
  registerTOTP,
} from "@/lib/zitadel";
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

  const { loginName, organization, sessionId, requestId, checkAfter } =
    searchParams;
  const { method } = params;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });
  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization,
  });

  const session = await loadMostRecentSession({
    serviceUrl,
    sessionParams: {
      loginName,
      organization,
    },
  });

  let totpResponse: RegisterTOTPResponse | undefined, error: Error | undefined;
  if (session && session.factors?.user?.id) {
    if (method === "time-based") {
      await registerTOTP({
        serviceUrl,
        userId: session.factors.user.id,
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
      await addOTPSMS({
        serviceUrl,
        userId: session.factors.user.id,
      }).catch((_error) => {
        // TODO: Throw this error?
        new Error("Could not add OTP via SMS");
      });
    } else if (method === "email") {
      await addOTPEmail({
        serviceUrl,
        userId: session.factors.user.id,
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
    urlToContinue = `/otp/${method}?` + paramsToContinue;

    // immediately check the OTP on the next page if sms or email was set up
    if (["email", "sms"].includes(method)) {
      return redirect(urlToContinue);
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
      <div className="flex flex-col items-center space-y-4">
        <h1>
          <Translated i18nKey="set.title" namespace="otp" />
        </h1>
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

        {totpResponse && "uri" in totpResponse && "secret" in totpResponse ? (
          <>
            <p className="ztdl-p">
              <Translated
                i18nKey="set.totpRegisterDescription"
                namespace="otp"
              />
            </p>
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
            </div>{" "}
          </>
        ) : (
          <>
            <p className="ztdl-p">
              {method === "email"
                ? "Code via email was successfully added."
                : method === "sms"
                  ? "Code via SMS was successfully added."
                  : ""}
            </p>

            <div className="mt-8 flex w-full flex-row items-center">
              <BackButton />
              <span className="flex-grow"></span>

              <Link href={urlToContinue}>
                <Button
                  type="submit"
                  className="self-end"
                  variant={ButtonVariants.Primary}
                >
                  <Translated i18nKey="set.submit" namespace="otp" />
                </Button>
              </Link>
            </div>
          </>
        )}
      </div>
    </DynamicTheme>
  );
}
