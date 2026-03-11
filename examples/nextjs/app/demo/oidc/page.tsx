import { getSession } from "@zitadel/nextjs";

export const dynamic = "force-dynamic";

function decodeIdTokenClaims(idToken?: string): Record<string, unknown> | null {
  if (!idToken) {
    return null;
  }

  const [, payload] = idToken.split(".");
  if (!payload) {
    return null;
  }

  try {
    const json = Buffer.from(payload, "base64url").toString("utf8");
    return JSON.parse(json);
  } catch {
    return null;
  }
}

function summarizeIdTokenClaims(claims: Record<string, unknown> | null) {
  if (!claims) {
    return null;
  }

  const keys = [
    "sub",
    "preferred_username",
    "name",
    "email",
    "email_verified",
    "iss",
    "aud",
  ];

  return Object.fromEntries(
    keys.flatMap((key) => (claims[key] === undefined ? [] : [[key, claims[key]]])),
  );
}

function describeFlowEvent(flow?: string, error?: string) {
  if (!flow) {
    return {
      flow: "none",
      description: "No recent redirect event. Use Sign in with ZITADEL to start.",
    };
  }

  const description =
    flow === "signed-in"
      ? "Callback succeeded and the session cookie was created."
      : flow === "signed-out"
        ? "Sign-out completed and redirected back to this page."
        : flow === "callback-error"
          ? "Callback failed before a valid session was created."
          : "OIDC redirect flow event received.";

  return {
    flow,
    description,
    error: error ?? null,
  };
}

export default async function OidcDemoPage(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const session = await getSession();
  const idTokenClaims = decodeIdTokenClaims(session?.idToken);
  const flowEvent = describeFlowEvent(searchParams.flow, searchParams.error);
  const now = Math.floor(Date.now() / 1000);

  const result = session
    ? {
        authenticated: true,
        expiresAt: session.expiresAt,
        expiresAtIso: new Date(session.expiresAt * 1000).toISOString(),
        expiresInSeconds: Math.max(session.expiresAt - now, 0),
        hasAccessToken: Boolean(session.accessToken),
        hasIdToken: Boolean(session.idToken),
        hasRefreshToken: Boolean(session.refreshToken),
        idTokenSummary: summarizeIdTokenClaims(idTokenClaims),
        idTokenClaims,
      }
    : {
        authenticated: false,
        reason: "No active OIDC session cookie found.",
      };

  return (
    <section className="ztdl-lane ztdl-noise">
      <h2>OIDC demo lane</h2>
      <p>Use this page to validate sign in, callback handling, and sign out redirects.</p>
      <ol>
        <li>
          Start with <a href="/api/auth/signin">Sign in with ZITADEL</a>.
        </li>
        <li>After login, you should return to this page with flow: signed-in.</li>
        <li>
          Use <a href="/api/auth/signout">Sign out</a> to clear the session and
          return here with flow: signed-out.
        </li>
      </ol>
      <h3>Latest flow event</h3>
      <pre>{JSON.stringify(flowEvent, null, 2)}</pre>
      <h3>OIDC session result</h3>
      <pre>{JSON.stringify(result, null, 2)}</pre>
    </section>
  );
}
