import { getSession, server } from "#/lib/zitadel";
import PasswordForm from "#/ui/PasswordForm";
import UserAvatar from "#/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";

async function loadSession(loginName: string) {
  try {
    const recent = await getMostRecentCookieWithLoginname(`${loginName}`);

    return getSession(server, recent.id, recent.token).then(({ session }) => {
      return session;
    });
  } catch (error) {
    throw new Error("Session could not be loaded!");
  }
}

export default async function Page({ searchParams }: { searchParams: any }) {
  const { loginName } = searchParams;
  const sessionFactors = await loadSession(loginName);

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>{`Welcome ${sessionFactors.factors?.user?.displayName}`}</h1>
      <p className="ztdl-p mb-6 block">You are signed in.</p>

      <UserAvatar
        loginName={loginName ?? sessionFactors.factors?.user?.loginName}
        displayName={sessionFactors.factors?.user?.displayName}
        showDropdown
      ></UserAvatar>
    </div>
  );
}
