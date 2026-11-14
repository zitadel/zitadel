import "server-only";

import { getDomainSettings, getOrganizationDomains, ServiceConfig } from "@/lib/server/zitadel";

/**
 * Removes an `@<domain>` suffix from a login name, ignoring case, for the first
 * of the domains it ends with, but only where it doubles an e-mail-style username:
 * what remains must itself contain an `@`.
 *
 * A login name with a single `@` is left alone. Without reading the user it
 * cannot be told whether `theodora@studio.example.com` is the username `theodora`
 * with the organization's domain appended or the username
 * `theodora@studio.example.com` as typed, and cutting the second one would be wrong.
 */
export function stripLoginNameSuffix(loginName: string, domains: string[]): string {
  const lowerLoginName = loginName.toLowerCase();
  const domain = domains.find((candidate) => lowerLoginName.endsWith(`@${candidate.toLowerCase()}`));
  if (!domain) {
    return loginName;
  }

  const username = loginName.slice(0, -(domain.length + 1));
  return username.includes("@") ? username : loginName;
}

/**
 * The login name as it should be shown to the user. With "user login must be
 * domain" a login name is `<username>@<organization domain>`, one per verified
 * domain of the organization; when the branding hides that suffix and the
 * organization's domain settings really append it, a suffix that doubles an
 * e-mail-style username is removed.
 *
 * Only the organization is read, never the user, so the result depends on the
 * login name and the organization alone: known and unknown users look the same
 * and "ignore unknown usernames" keeps its effect. Falls back to the login name
 * when the organization's settings cannot be read.
 */
export async function getShownLoginName({
  serviceConfig,
  loginName,
  organizationId,
  hideLoginNameSuffix,
}: {
  serviceConfig: ServiceConfig;
  loginName?: string;
  organizationId?: string;
  hideLoginNameSuffix?: boolean;
}): Promise<string | undefined> {
  if (!hideLoginNameSuffix || !loginName || !organizationId) {
    return loginName;
  }

  const [domainSettings, domains] = await Promise.all([
    getDomainSettings({ serviceConfig, organization: organizationId }).catch(() => undefined),
    getOrganizationDomains({ serviceConfig, organizationId }).catch(() => []),
  ]);

  if (!domainSettings?.loginNameIncludesDomain) {
    return loginName;
  }

  return stripLoginNameSuffix(loginName, domains);
}
