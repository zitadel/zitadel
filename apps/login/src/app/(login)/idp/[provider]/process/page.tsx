import { IdpProcessHandler } from "@/components/idp-process-handler";

/**
 * This page handles the initial IDP callback with the single-use token.
 * It delegates to a client component which calls the server action.
 * The client component is needed so that cookies can be set properly.
 */
export default async function ProcessPage(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const params = await props.params;
  const searchParams = await props.searchParams;

  const { provider } = params;
  const { id, token, requestId, organization, link, postErrorRedirectUrl, linkToSessionId, linkFingerprint } = searchParams;

  // Validate required parameters before passing to client component
  if (!id || !token) {
    throw new Error("Missing required IDP callback parameters");
  }

  return (
    <IdpProcessHandler
      provider={provider}
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
