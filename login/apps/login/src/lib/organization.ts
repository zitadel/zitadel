import { getDefaultOrg } from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";

export async function getEffectiveOrganizationId({
  serviceUrl,
  organization,
}: {
  serviceUrl: string;
  organization?: string;
}): Promise<string | undefined> {
  if (organization) {
    return organization;
  }

  if (process.env.DEFAULT_ORGANIZATION_ID) {
    return process.env.DEFAULT_ORGANIZATION_ID;
  }

  const defaultOrg: Organization | null = await getDefaultOrg({ serviceUrl });
  return defaultOrg?.id;
} 