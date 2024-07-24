import { createCallback, getBrandingSettings, getSession } from "@/lib/zitadel";
import DynamicTheme from "@/ui/DynamicTheme";
import UserAvatar from "@/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "@/utils/cookies";
import { redirect } from "next/navigation";

async function loadSession(loginName: string, authRequestId?: string) {
  const recent = await getMostRecentCookieWithLoginname(`${loginName}`);

  if (authRequestId) {
    return createCallback({
      authRequestId,
      callbackKind: {
        case: "session",
        value: { sessionId: recent.id, sessionToken: recent.token },
      },
    }).then(({ callbackUrl }) => {
      return redirect(callbackUrl);
    });
  }
  return getSession(recent.id, recent.token).then((response) => {
    if (response?.session) {
      return response.session;
    }
  });
}

export default async function Page({ searchParams }: { searchParams: any }) {
  const { loginName, authRequestId, organization } = searchParams;
  const sessionFactors = await loadSession(loginName, authRequestId);

  const branding = await getBrandingSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{`Welcome ${sessionFactors?.factors?.user?.displayName}`}</h1>
        <p className="ztdl-p mb-6 block">You are signed in.</p>

        <UserAvatar
          loginName={loginName ?? sessionFactors?.factors?.user?.loginName}
          displayName={sessionFactors?.factors?.user?.displayName}
          showDropdown
          searchParams={searchParams}
        ></UserAvatar>
      </div>
    </DynamicTheme>
  );
}
