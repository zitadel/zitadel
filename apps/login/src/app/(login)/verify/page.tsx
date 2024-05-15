import { getBrandingSettings, server } from "@/lib/zitadel";
import DynamicTheme from "@/ui/DynamicTheme";
import VerifyEmailForm from "@/ui/VerifyEmailForm";
import { ExclamationTriangleIcon } from "@heroicons/react/24/outline";

export default async function Page({ searchParams }: { searchParams: any }) {
  const {
    userId,
    sessionId,
    code,
    submit,
    organization,
    authRequestId,
    loginname,
    passwordset,
  } = searchParams;

  const branding = await getBrandingSettings(server, organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Verify user</h1>
        <p className="ztdl-p mb-6 block">
          Enter the Code provided in the verification email.
        </p>

        {userId ? (
          <VerifyEmailForm
            userId={userId}
            code={code}
            submit={submit === "true"}
            organization={organization}
            authRequestId={authRequestId}
            sessionId={sessionId}
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
