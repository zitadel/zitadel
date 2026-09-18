import { User } from "@zitadel/proto/zitadel/user/v2/user_pb";

/**
 * Whether `loginName` is one of the user's login names. A user holds one per
 * verified domain, and searchUsers matches them case-insensitively, so checking
 * only the preferred one rejects a valid name the search itself just matched.
 */
export function isLoginNameOfUser(user: Pick<User, "preferredLoginName" | "loginNames">, loginName: string) {
  const wanted = loginName.toLowerCase();
  return [user.preferredLoginName, ...(user.loginNames ?? [])].some((name) => name?.toLowerCase() === wanted);
}
