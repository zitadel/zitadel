import { getBrandingSettings, getLoginSettings } from "@/lib/zitadel";
import Alert from "@/ui/Alert";
import DynamicTheme from "@/ui/DynamicTheme";
import VerifyEmailForm from "@/ui/VerifyEmailForm";
import { ExclamationTriangleIcon } from "@heroicons/react/24/outline";

export default async function Page({ searchParams }: { searchParams: any }) {
  const {
    userId,
    loginName,
    sessionId,
    code,
    submit,
    organization,
    authRequestId,
  } = searchParams;

  const branding = await getBrandingSettings(organization);

  const loginSettings = await getLoginSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Verify user</h1>
        <p className="ztdl-p mb-6 block">
          Enter the Code provided in the verification email.
        </p>

        {!userId && (
          <div className="py-4">
            <Alert>
              Could not get the context of the user. Make sure to provide a
              userId as searchParam.
            </Alert>
          </div>
        )}

        {userId ? (
          <VerifyEmailForm
            userId={userId}
            loginName={loginName}
            code={code}
            submit={submit === "true"}
            organization={organization}
            authRequestId={authRequestId}
            sessionId={sessionId}
            loginSettings={loginSettings}
          />
        ) : (
          <div className="w-full flex flex-row items-center justify-center border border-yellow-600/40 dark:border-yellow-500/20 bg-yellow-200/30 text-yellow-600 dark:bg-yellow-700/20 dark:text-yellow-200 rounded-md py-2 scroll-px-40">
            <ExclamationTriangleIcon className="h-5 w-5 mr-2" />
            <span className="text-center text-sm">No userId provided!</span>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}
