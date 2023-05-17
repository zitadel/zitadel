import { getSession, server } from "#/lib/zitadel";
import PasswordForm from "#/ui/PasswordForm";
import UserAvatar from "#/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";

async function loadSession(loginName: string) {
  const recent = await getMostRecentCookieWithLoginname(loginName);

  return getSession(server, recent.id, recent.token).then(({ session }) => {
    console.log("ss", session);
    return session;
  });
}

export default async function Page({ searchParams }: { searchParams: any }) {
  const { loginName } = searchParams;

  const sessionFactors = await loadSession(loginName);

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>{sessionFactors.factors?.user?.displayName ?? "Password"}</h1>
      <p className="ztdl-p mb-6 block">Enter your password.</p>

      <UserAvatar
        loginName={loginName ?? sessionFactors.factors?.user?.loginName}
        displayName={sessionFactors.factors?.user?.displayName}
        showDropdown
      ></UserAvatar>

      <PasswordForm />
    </div>
  );
}
