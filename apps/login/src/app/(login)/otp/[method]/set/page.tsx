import {
  addOTPEmail,
  addOTPSMS,
  getBrandingSettings,
  getSession,
  registerTOTP,
} from "@/lib/zitadel";
import Alert from "@/ui/Alert";
import { Button, ButtonVariants } from "@/ui/Button";
import DynamicTheme from "@/ui/DynamicTheme";
import { Spinner } from "@/ui/Spinner";
import TOTPRegister from "@/ui/TOTPRegister";
import UserAvatar from "@/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "@/utils/cookies";
import Link from "next/link";
import { RegisterTOTPResponse } from "@zitadel/proto/zitadel/user/v2beta/user_service_pb";

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, organization, sessionId, authRequestId, checkAfter } =
    searchParams;
  const { method } = params;

  const branding = await getBrandingSettings(organization);
  const { session, token } = await loadSession(loginName, organization);

  let totpResponse: RegisterTOTPResponse | undefined,
    totpError: Error | undefined;
  if (session && session.factors?.user?.id) {
    if (method === "time-based") {
      await registerTOTP(session.factors.user.id)
        .then((resp) => {
          if (resp) {
            totpResponse = resp;
          }
        })
        .catch((error) => {
          totpError = error;
        });
    } else if (method === "sms") {
      // does not work
      await addOTPSMS(session.factors.user.id);
    } else if (method === "email") {
      // works
      await addOTPEmail(session.factors.user.id);
    } else {
      throw new Error("Invalid method");
    }
  } else {
    throw new Error("No session found");
  }

  async function loadSession(loginName?: string, organization?: string) {
    const recent = await getMostRecentCookieWithLoginname(
      loginName,
      organization,
    );

    return getSession(recent.id, recent.token).then((response) => {
      return { session: response?.session, token: recent.token };
    });
  }

  const paramsToContinue = new URLSearchParams({});
  let urlToContinue = "/accounts";

  if (authRequestId && sessionId) {
    if (sessionId) {
      paramsToContinue.append("sessionId", sessionId);
    }
    if (authRequestId) {
      paramsToContinue.append("authRequestId", authRequestId);
    }
    if (organization) {
      paramsToContinue.append("organization", organization);
    }
    urlToContinue = `/login?` + paramsToContinue;
  } else if (loginName) {
    if (loginName) {
      paramsToContinue.append("loginName", loginName);
    }
    if (authRequestId) {
      paramsToContinue.append("authRequestId", authRequestId);
    }
    if (organization) {
      paramsToContinue.append("organization", organization);
    }

    urlToContinue = `/signedin?` + paramsToContinue;
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Register 2-factor</h1>
        {!session && (
          <div className="py-4">
            <Alert>
              Could not get the context of the user. Make sure to enter the
              username first or provide a loginName as searchParam.
            </Alert>
          </div>
        )}

        {totpError && (
          <div className="py-4">
            <Alert>{totpError?.message}</Alert>
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
              Scan the QR Code or navigate to the URL manually.
            </p>
            <div>
              {/* {auth && <div>{auth.to}</div>} */}

              <TOTPRegister
                uri={totpResponse.uri as string}
                secret={totpResponse.secret as string}
                loginName={loginName}
                sessionId={sessionId}
                authRequestId={authRequestId}
                organization={organization}
                checkAfter={checkAfter === "true"}
              ></TOTPRegister>
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
              <span className="flex-grow"></span>
              <Link
                href={
                  checkAfter
                    ? `/otp/${method}?` + new URLSearchParams()
                    : urlToContinue
                }
              >
                <Button
                  type="submit"
                  className="self-end"
                  variant={ButtonVariants.Primary}
                >
                  continue
                </Button>
              </Link>
            </div>
          </>
        )}
      </div>
    </DynamicTheme>
  );
}
