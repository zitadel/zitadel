import "server-only";

import { isClassifiedError } from "@/lib/server/grpc/interceptors/error-classification";
import { Code } from "@zitadel/client";
import { getTranslations } from "next-intl/server";

// The ZITADEL error IDs returned when a user may authenticate but not access the application,
// mapped to their key in the "error" translation namespace.
const AUTHORIZATION_ERRORS = {
  "Errors.User.GrantRequired": "grantRequired",
  "Errors.User.ProjectRequired": "projectRequired",
} as const;

/**
 * Returns a translated message for errors where the user is not allowed to access the
 * application (missing user grant or project grant), or undefined for any other error.
 */
export async function getAuthorizationErrorMessage(error: unknown): Promise<string | undefined> {
  if (!isClassifiedError(error) || error.code !== Code.PermissionDenied) {
    return undefined;
  }

  const match = Object.entries(AUTHORIZATION_ERRORS).find(([errorId]) => error.rawMessage.includes(errorId));
  if (!match) {
    return undefined;
  }

  const t = await getTranslations("error");
  return t(match[1]);
}
