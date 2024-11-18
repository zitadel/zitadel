import { redirect } from "next/navigation";

type FinishFlowCommand =
  | {
      sessionId: string;
      authRequestId: string;
    }
  | { loginName: string };

/**
 * on client: redirects user back to OIDC application or to a success page
 * @param command
 * @returns
 */
export function finishFlow(
  command: FinishFlowCommand & { organization?: string },
) {
  return "sessionId" in command && "authRequestId" in command
    ? redirect(
        `/login?` +
          new URLSearchParams({
            sessionId: command.sessionId,
            authRequest: command.authRequestId,
          }),
      )
    : redirect(
        `/signedin?` +
          new URLSearchParams({
            loginName: command.loginName,
          }),
      );
}
