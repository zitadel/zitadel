import { createCallback, getSession, server } from "#/lib/zitadel";
import UserAvatar from "#/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";
import { redirect } from "next/navigation";
import { NextResponse } from "next/server";

async function loadSession(loginName: string, authRequestId?: string) {
  const recent = await getMostRecentCookieWithLoginname(`${loginName}`);

  if (authRequestId) {
    console.log(authRequestId);
    return createCallback(server, {
      authRequestId,
      session: { sessionId: recent.id, sessionToken: recent.token },
    }).then(({ callbackUrl }) => {
      return redirect(callbackUrl);
    });
  }
  return getSession(server, recent.id, recent.token).then((response) => {
    if (response?.session) {
      return response.session;
    }
  });
}

export default async function Page({ searchParams }: { searchParams: any }) {
  const { loginName, authRequestId } = searchParams;
  const sessionFactors = await loadSession(loginName, authRequestId);

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>{`Welcome ${sessionFactors?.factors?.user?.displayName}`}</h1>
      <p className="ztdl-p mb-6 block">You are signed in.</p>

      <UserAvatar
        loginName={loginName ?? sessionFactors?.factors?.user?.loginName}
        displayName={sessionFactors?.factors?.user?.displayName}
        showDropdown
      ></UserAvatar>
    </div>
  );
}
