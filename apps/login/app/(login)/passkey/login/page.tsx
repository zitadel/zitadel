import Alert from "#/ui/Alert";
import LoginPasskey from "#/ui/LoginPasskey";
import { ChallengeKind } from "@zitadel/server";

async function updateSessionAndCookie(loginName: string) {
  const res = await fetch(
    `${process.env.VERCEL_URL ?? "http://localhost:3000"}/session`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        loginName,
        challenges: [ChallengeKind.CHALLENGE_KIND_PASSKEY],
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

const title = "Authenticate with a passkey";
const description =
  "Your device will ask for your fingerprint, face, or screen lock";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName } = searchParams;
  if (loginName) {
    console.log(loginName);
    const session = await updateSessionAndCookie(loginName);

    console.log("sess", session);
    const challenge = session?.challenges?.passkey;

    console.log(challenge);

    return (
      <div className="flex flex-col items-center space-y-4">
        <h1>{title}</h1>

        {/* {sessionFactors && (
        <UserAvatar
          loginName={loginName ?? sessionFactors.factors?.user?.loginName}
          displayName={sessionFactors.factors?.user?.displayName}
          showDropdown
        ></UserAvatar>
      )}
      <p className="ztdl-p mb-6 block">{description}</p>

      {!sessionFactors && (
        <div className="py-4">
          <Alert>
            Could not get the context of the user. Make sure to enter the
            username first or provide a loginName as searchParam.
          </Alert>
        </div>
      )} */}

        {challenge && <LoginPasskey challenge={challenge} />}
      </div>
    );
  } else {
    return (
      <div className="flex flex-col items-center space-y-4">
        <h1>{title}</h1>

        {/* {sessionFactors && (
              <UserAvatar
                loginName={loginName ?? sessionFactors.factors?.user?.loginName}
                displayName={sessionFactors.factors?.user?.displayName}
                showDropdown
              ></UserAvatar>
            )}
            <p className="ztdl-p mb-6 block">{description}</p>
      
            {!sessionFactors && (
              <div className="py-4">
                
              </div>
            )} */}

        <Alert>Provide your active session as loginName param</Alert>
      </div>
    );
  }
}
