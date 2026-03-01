import { IdpProcessHandler } from "@/components/idp-process-handler";

/**
 * This page handles the LDAP authentication success callback.
 * Unlike other IdPs that use OAuth redirects, LDAP authentication happens
 * server-side and returns the intent data directly. This page processes
 * that data to create a session and complete the login flow.
 */
export default async function LDAPSuccessPage(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;

  const {
    id,
    token,
    requestId,
    organization,
    link,
    postErrorRedirectUrl,
    linkToSessionId,
    linkFingerprint,
  } = searchParams;

  // Validate required parameters before passing to client component
  const missingParams: string[] = [];
  if (!id) {
    missingParams.push("id");
  }
  if (!token) {
    missingParams.push("token");
  }
  if (missingParams.length > 0) {
    const paramLabel = missingParams.length === 1 ? "parameter" : "parameters";
    throw new Error(
      `Missing required LDAP callback ${paramLabel}: ${missingParams.join(", ")}`
    );
  }

  return (
    <IdpProcessHandler
      provider="ldap"
      id={id}
      token={token}
      requestId={requestId}
      organization={organization}
      link={link}
      sessionId={linkToSessionId}
      linkFingerprint={linkFingerprint}
      postErrorRedirectUrl={postErrorRedirectUrl}
    />
  );
}
