import "server-only";

import { ClassifiedConnectError, isClassifiedError } from "@/lib/server/grpc/interceptors/error-classification";
import { Code } from "@connectrpc/connect";
import { ErrorDetailSchema } from "@zitadel/proto/zitadel/message_pb";

/**
 * Whether ZITADEL rejected a call because of the user's input or state (a wrong code, a locked
 * user, a rate limit), so an action should answer with a message instead of failing with a 500.
 *
 * Unauthenticated and PermissionDenied do not count: the login calls the API with its own
 * credentials, so they mean the login itself was rejected, an outage that must stay visible.
 * Callers that send a user's session token handle a rejected token themselves.
 */
export function isRejectedByUser(error: unknown): error is ClassifiedConnectError {
  return (
    isClassifiedError(error) &&
    error.isUserError &&
    error.code !== Code.Unauthenticated &&
    error.code !== Code.PermissionDenied
  );
}

/** Ids of the ZITADEL errors the login answers with a specific message, see {@link zitadelErrorId}. */
export const ZitadelErrorIds: Readonly<
  Record<"userLocked" | "totpReused" | "codeInvalid" | "codeExpired", readonly string[]>
> = {
  /** the user is locked after too many failed TOTP or OTP SMS/email checks */
  userLocked: ["COMMAND-SF3fg", "COMMAND-S6h4R"],
  /** the TOTP code was already used for a successful check */
  totpReused: ["TOTP-Auw0a"],
  /** the code does not match */
  codeInvalid: ["CODE-woT0xc"],
  /** the code expired, or there is no code to check (never sent, or already used) */
  codeExpired: [
    "CODE-QvUQ4P",
    "TWILI-SF3ba",
    "COMMAND-2M9fs",
    "COMMAND-G3t31",
    "COMMAND-SLr02",
    "COMMAND-S34gh",
    "TWILI-JK3ta",
    "TWILI-Ok39a",
  ],
};

/**
 * The id of a ZITADEL error, e.g. `COMMAND-SF3fg`. The API translates the message, so the id is
 * what tells specific rejections apart. It is sent as an error detail and also ends the message.
 */
export function zitadelErrorId(error: ClassifiedConnectError): string | undefined {
  const detailId = error.findDetails(ErrorDetailSchema)[0]?.id;
  if (detailId) {
    return detailId;
  }

  return /\(([A-Za-z]+-[A-Za-z0-9]+)\)$/.exec(error.rawMessage)?.[1];
}

/** Whether the error carries one of the given ZITADEL error ids. */
export function hasZitadelErrorId(error: ClassifiedConnectError, ids: readonly string[]): boolean {
  const id = zitadelErrorId(error);
  return id !== undefined && ids.includes(id);
}
