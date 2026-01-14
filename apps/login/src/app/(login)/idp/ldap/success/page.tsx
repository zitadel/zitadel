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
  if (!id || !token) {
    throw new Error("Missing required LDAP callback parameters");
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
