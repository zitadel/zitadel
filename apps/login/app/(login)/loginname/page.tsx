import {
  getLoginSettings,
  listAuthenticationMethodTypes,
  server,
} from "#/lib/zitadel";
import UsernameForm from "#/ui/UsernameForm";
import { AuthenticationMethodType, Factors } from "@zitadel/server";
import { redirect } from "next/navigation";

type SessionAuthMethods = {
  authMethodTypes: AuthenticationMethodType[];
  sessionId: string;
  factors: Factors;
};

async function createSessionAndCookie(loginName: string) {
  const res = await fetch(
    `${process.env.VERCEL_URL ?? "http://localhost:3000"}/session`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
      }),
      next: { revalidate: 0 },
    }
  );

  const response = await res.json();

  if (!res.ok) {
    return Promise.reject(response.details);
  }
  return response;
}

async function createSessionAndGetAuthMethods(
  loginName: string
): Promise<SessionAuthMethods> {
  return createSessionAndCookie(loginName)
    .then((resp) => {
      return listAuthenticationMethodTypes(resp.factors.user.id)
        .then((methods) => {
          return {
            authMethodTypes: methods.authMethodTypes,
            sessionId: resp.sessionId,
            factors: resp?.factors,
          };
        })
        .catch((error) => {
          throw "Could not get auth methods";
        });
    })
    .catch((error) => {
      console.log(error);
      throw "Could not add session to cookie";
    });
}

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const loginName = searchParams?.loginName;
  if (loginName) {
    const login = await getLoginSettings(server);
    const sessionAndAuthMethods = await createSessionAndGetAuthMethods(
      loginName
    );
    if (sessionAndAuthMethods.authMethodTypes.length == 1) {
      const method = sessionAndAuthMethods.authMethodTypes[0];
      switch (method) {
        case AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSWORD:
          return redirect("/password?" + new URLSearchParams({ loginName }));
        case AuthenticationMethodType.AUTHENTICATION_METHOD_TYPE_PASSKEY:
          return redirect(
            "/passkey/login?" + new URLSearchParams({ loginName })
          );
        default:
          return redirect("/password?" + new URLSearchParams({ loginName }));
      }
    }
    return (
      <div className="flex flex-col items-center space-y-4">
        <h1>Welcome back!</h1>
        <p className="ztdl-p">Enter your login data.</p>

        <UsernameForm />
      </div>
    );
  } else {
    return (
      <div className="flex flex-col items-center space-y-4">
        <h1>Welcome back!</h1>
        <p className="ztdl-p">Enter your login data.</p>

        <UsernameForm />
      </div>
    );
  }
}
