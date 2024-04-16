import { getBrandingSettings, server } from "#/lib/zitadel";
import DynamicTheme from "#/ui/DynamicTheme";

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

        <div></div>
      </div>
    </DynamicTheme>
  );
}
