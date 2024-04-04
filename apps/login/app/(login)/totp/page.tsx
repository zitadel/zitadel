import { getBrandingSettings, getLoginSettings, server } from "#/lib/zitadel";
import DynamicTheme from "#/ui/DynamicTheme";
import TOTPForm from "#/ui/TOTPForm";

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
        <p className="ztdl-p">Enter the code from your authenticator app. </p>

        <TOTPForm
          loginName={loginName}
          sessionId={sessionId}
          code={code}
          authRequestId={authRequestId}
          organization={organization}
          submit={submit === "true"}
        />
      </div>
    </DynamicTheme>
  );
}
