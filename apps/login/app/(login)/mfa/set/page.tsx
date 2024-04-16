import { getBrandingSettings, server } from "#/lib/zitadel";
import { Button, ButtonVariants } from "#/ui/Button";
import DynamicTheme from "#/ui/DynamicTheme";
import { TextInput } from "#/ui/Input";
import UserAvatar from "#/ui/UserAvatar";
import { useRouter } from "next/navigation";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, authRequestId, sessionId, organization, code, submit } =
    searchParams;

  const branding = await getBrandingSettings(server, organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Verify 2-Factor</h1>

        <p className="ztdl-p">Choose one of the following second factors.</p>

        <UserAvatar
          showDropdown
          displayName="Max Peintner"
          loginName="max@zitadel.com"
        ></UserAvatar>
        <div className="w-full">
          <TextInput type="password" label="Password" />
        </div>
      </div>
    </DynamicTheme>
  );
}
