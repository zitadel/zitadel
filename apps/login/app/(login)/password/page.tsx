import { getSession, server } from "#/lib/zitadel";
import PasswordForm from "#/ui/PasswordForm";
import UserAvatar from "#/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";

async function loadSession(loginName: string) {
  const recent = await getMostRecentCookieWithLoginname(loginName);
  console.log("found recent cookie: ", recent);

  return getSession(server, recent.id, recent.token).then(({ session }) => {
    console.log(session);
    return session;
  });
  //   const res = await fetch(
  //     `http://localhost:3000/session?` +
  //       new URLSearchParams({
  //         loginName: loginName,
  //       }),
  //     {
  //       method: "GET",
  //       headers: {
  //         "Content-Type": "application/json",
  //       },
  //     }
  //   );

  //   if (!res.ok) {
  //     throw new Error("Failed to load session");
  //   }

  //   return res.json();
}

export default async function Page({ searchParams }: { searchParams: any }) {
  const { loginName } = searchParams;
  console.log(loginName);

  const sessionFactors = await loadSession(loginName);
  console.log(sessionFactors);

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>{sessionFactors.factors.user.displayName}</h1>
      <p className="ztdl-p mb-6 block">Enter your password.</p>

      <UserAvatar loginName={loginName} showDropdown></UserAvatar>

      <PasswordForm />
    </div>
  );
}
