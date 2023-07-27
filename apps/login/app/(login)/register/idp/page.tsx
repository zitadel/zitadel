import { getLegalAndSupportSettings, server } from "#/lib/zitadel";
import { SignInWithIDP } from "#/ui/SignInWithIDP";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const legal = await getLegalAndSupportSettings(server);

  console.log(server);

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Register</h1>
      <p className="ztdl-p">Create your ZITADEL account.</p>

      {legal && <SignInWithIDP server={server}></SignInWithIDP>}
    </div>
  );
}
