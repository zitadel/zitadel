'use server';

import { redirect } from 'next/navigation';

/**
 * Server Action to handle authentication redirects
 * This avoids client-side navigation that triggers RSC requests
 */
export async function redirectToLogin(sessionId: string, requestId: string, organization?: string) {
  const params = new URLSearchParams({
    sessionId,
    requestId,
  });

  if (organization) {
    params.append("organization", organization);
  }

  // This server-side redirect doesn't trigger RSC requests
  redirect(`/login?${params}`);
}

export async function redirectToUrl(url: string) {
  // Server-side redirect for any URL
  redirect(url);
}