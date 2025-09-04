'use server';

import { redirect } from 'next/navigation';

// Re-export the main auth flow action from auth-flow.ts
export { completeAuthFlowAction } from './auth-flow';

export async function redirectToUrl(url: string) {
  // Server-side redirect for any URL (kept for other use cases)
  redirect(url);
}