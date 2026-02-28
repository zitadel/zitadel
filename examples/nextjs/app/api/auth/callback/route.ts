import { handleCallback } from "@zitadel/nextjs/auth/oidc";

const OIDC_DEMO_PATH = "/demo/oidc";

export async function GET(request: Request) {
  const callbackUrl = new URL("/api/auth/callback", request.url).toString();
  const postLoginUrl = new URL(OIDC_DEMO_PATH, request.url);

  try {
    await handleCallback(request, { callbackUrl });
    postLoginUrl.searchParams.set("flow", "signed-in");
  } catch (error) {
    postLoginUrl.searchParams.set("flow", "callback-error");
    postLoginUrl.searchParams.set(
      "error",
      (error instanceof Error ? error.message : "Unknown callback error").slice(
        0,
        200,
      ),
    );
  }

  return Response.redirect(postLoginUrl);
}
