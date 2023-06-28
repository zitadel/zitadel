import { getLoginSettings, server } from "#/lib/zitadel";
import UsernameForm from "#/ui/UsernameForm";

export default async function Page({
  params,
}: {
  params: { loginname: string };
}) {
  const login = await getLoginSettings(server);

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Welcome back!</h1>
      <p className="ztdl-p">Enter your login data.</p>

      <UsernameForm />
    </div>
  );
}
