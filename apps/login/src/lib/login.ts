import { redirect } from "next/navigation";
import { getLoginSettings } from "./zitadel";

type FinishFlowCommand =
  | {
      sessionId: string;
      authRequestId: string;
    }
  | { loginName: string };

/**
 * for client: redirects user back to OIDC application or to a success page when using authRequestId, check if a default redirect and redirect to it, or just redirect to a success page with the loginName
 * @param command
 * @returns
 */
export async function finishFlow(
  command: FinishFlowCommand & { organization?: string },
) {
  if ("sessionId" in command && "authRequestId" in command) {
    return redirect(
      `/login?` +
        new URLSearchParams({
          sessionId: command.sessionId,
          authRequest: command.authRequestId,
        }),
    );
  }

  const loginSettings = await getLoginSettings(command.organization);
  if (loginSettings?.defaultRedirectUri) {
    return redirect(loginSettings.defaultRedirectUri);
  }

  return redirect(
    `/signedin?` +
      new URLSearchParams({
        loginName: command.loginName,
      }),
  );
}
