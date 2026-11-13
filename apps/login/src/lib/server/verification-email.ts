import "server-only";

import { createLogger } from "@/lib/server/logger";
import { getServiceConfig } from "@/lib/server/service-url";
import { createInviteCode, sendEmailCode } from "@/lib/server/zitadel";
import { headers } from "next/headers";
import { getPublicHostWithProtocol } from "./host";

const logger = createLogger("verify");

export function buildVerificationUrlTemplate(
  hostWithProtocol: string,
  basePath: string,
  isInvite: boolean,
  requestId?: string,
): string {
  let urlTemplate = `${hostWithProtocol}${basePath}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}`;

  if (isInvite) {
    urlTemplate += "&invite=true";
  }

  if (requestId) {
    urlTemplate += `&requestId=${encodeURIComponent(requestId)}`;
  }

  return urlTemplate;
}

type TrySendVerificationCommand = {
  userId: string;
  isInvite: boolean;
  requestId?: string;
};

/**
 * Attempts to send an initial verification email/invite code.
 * Returns `true` if sent successfully, `false` on error (swallowed and logged).
 *
 * The caller decides which user receives the email, so this must never be reachable as a
 * server action. Only server code that resolved the user itself may call it.
 */
export async function trySendVerification(command: TrySendVerificationCommand): Promise<boolean> {
  try {
    const _headers = await headers();
    const { serviceConfig } = getServiceConfig(_headers);
    const hostWithProtocol = await getPublicHostWithProtocol(_headers);

    const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";
    const urlTemplate = buildVerificationUrlTemplate(hostWithProtocol, basePath, command.isInvite, command.requestId);

    if (command.isInvite) {
      await createInviteCode({
        serviceConfig,
        userId: command.userId,
        urlTemplate,
      });
    } else {
      await sendEmailCode({
        serviceConfig,
        userId: command.userId,
        urlTemplate,
      });
    }

    logger.info("Verification email sent successfully", { userId: command.userId, isInvite: command.isInvite });
    return true;
  } catch (err) {
    logger.error("Failed to send verification email", { userId: command.userId, isInvite: command.isInvite, error: err });
    return false;
  }
}
