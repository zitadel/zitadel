import { getOrganizationPrimaryDomain, ServiceConfig } from "@/lib/zitadel";

/**
 * Removes the `@<domain>` suffix from a login name, ignoring case. A login name
 * that is nothing but the suffix, or does not end with it, is returned unchanged.
 */
export function stripLoginNameSuffix(loginName: string, domain?: string): string {
  if (!domain) {
    return loginName;
  }

  const suffix = `@${domain}`;
  if (loginName.length > suffix.length && loginName.toLowerCase().endsWith(suffix.toLowerCase())) {
    return loginName.slice(0, -suffix.length);
  }

  return loginName;
}

/**
 * The login name as it should be shown to the user. With "user login must be
 * domain" a login name is `<username>@<organization primary domain>`; when the
 * branding hides that suffix, the suffix is removed. Only the organization is
 * read, never the user, so the result is the same for known and unknown users
 * and does not weaken "ignore unknown usernames". Falls back to the login name
 * when the organization cannot be read.
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

  const domain = await getOrganizationPrimaryDomain({ serviceConfig, organizationId }).catch(() => undefined);

  return stripLoginNameSuffix(loginName, domain);
}
