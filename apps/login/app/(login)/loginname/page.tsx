import {
  createSession,
  getLoginSettings,
  listAuthenticationMethodTypes,
  server,
} from "#/lib/zitadel";
import UsernameForm from "#/ui/UsernameForm";
import { AuthenticationMethodType, Factors } from "@zitadel/server";

type SessionAuthMethods = {
  authMethodTypes: AuthenticationMethodType[];
  sessionId: string;
  factors: Factors;
};

async function updateCookie(loginName: string) {
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
    }
  );

  const response = await res.json();

  if (!res.ok) {
    console.log("damn");
    return Promise.reject(response.details);
  }
  return response;
}

async function getSessionAndAuthMethods(
  loginName: string,
  domain: string
): Promise<SessionAuthMethods> {
  const createdSession = await createSession(
    server,
    loginName,
    domain,
    undefined,
    undefined
  );

  if (createdSession) {
    return updateCookie(loginName)
      .then((resp) => {
        return listAuthenticationMethodTypes(resp.factors.user.id)
          .then((methods) => {
            return {
              authMethodTypes: methods.authMethodTypes,
              sessionId: createdSession.sessionId,
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
  } else {
    throw "Could not create session";
  }
}

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const domain: string = process.env.VERCEL_URL ?? "localhost";

  const loginName = searchParams?.loginName;
  if (loginName) {
    const login = await getLoginSettings(server);
    console.log(login);
    const sessionAndAuthMethods = await getSessionAndAuthMethods(
      loginName,
      domain
    );
    console.log(sessionAndAuthMethods);
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
