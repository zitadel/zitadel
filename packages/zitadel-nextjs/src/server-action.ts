import { getSession, type SessionData } from "./session.js";

/**
 * Higher-order function that wraps a server action with authentication checks.
 *
 * Reads the session and passes it to the inner handler. Throws an error
 * if no valid session exists, which Next.js will surface as a server action error.
 *
 * @example
 * ```ts
 * "use server";
 * import { protectedAction } from "@zitadel/nextjs";
 *
 * export const updateProfile = protectedAction(
 *   async (session, formData: FormData) => {
 *     // session.accessToken is available here
 *     const name = formData.get("name") as string;
 *     // ... call ZITADEL API with session.accessToken
 *   },
 * );
 * ```
 */
export function protectedAction<TArgs extends unknown[], TResult>(
  action: (session: SessionData, ...args: TArgs) => Promise<TResult>,
  options?: {
    /** Cookie secret. Falls back to ZITADEL_COOKIE_SECRET. */
    cookieSecret?: string;
  },
): (...args: TArgs) => Promise<TResult> {
  return async (...args: TArgs) => {
    const session = await getSession(options?.cookieSecret);
    if (!session) {
      throw new Error("Unauthorized: no valid session");
    }
    return action(session, ...args);
  };
}
