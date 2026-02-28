import { signIn } from "@zitadel/nextjs/auth/oidc";

export async function GET(request: Request) {
  const url = new URL(request.url);
  const prompt = url.searchParams.get("prompt") ?? undefined;
  const callbackUrl = new URL("/api/auth/callback", request.url).toString();
  await signIn({ callbackUrl, prompt });
}
