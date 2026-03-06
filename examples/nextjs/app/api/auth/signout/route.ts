import { signOut } from "@zitadel/nextjs/auth/oidc";

const OIDC_DEMO_PATH = "/demo/oidc";

export async function GET(request: Request) {
  const postLogoutUrl = new URL(OIDC_DEMO_PATH, request.url);
  postLogoutUrl.searchParams.set("flow", "signed-out");
  await signOut({ postLogoutRedirectUri: postLogoutUrl.toString() });
  return Response.redirect(postLogoutUrl);
}
