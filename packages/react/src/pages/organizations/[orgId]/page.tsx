import { fetchOrganization } from "../../../api/organizations"
import { fetchUsers } from "../../../api/fetch-users"
import { countOrgProjects, countOrgApplications } from "../../../api/org-resources"
import { listOrgDomains, listOrgMetadata, type OrgDomain, type OrgMetadataEntry } from "../../../api/org-settings"
import { OrgDetailClient } from "./org-detail-client"

interface Params {
  orgId: string
}

/**
 * Organization detail page — server component.
 */
export default async function OrganizationDetailPage({ params }: { params: Promise<Params> }) {
  const { orgId } = await params
  let organization: any = null
  let users: any[] = []
  let projectCount = 0
  let applicationCount = 0
  let domains: OrgDomain[] = []
  let metadata: OrgMetadataEntry[] = []
  let error: string | null = null

  try {
    const [org, usersResult, projCount, appCount, domainList, metadataList] = await Promise.all([
      fetchOrganization(orgId),
      fetchUsers(orgId),
      countOrgProjects(orgId).catch(() => 0),
      countOrgApplications(orgId).catch(() => 0),
      listOrgDomains(orgId).catch(() => [] as OrgDomain[]),
      listOrgMetadata(orgId).catch(() => [] as OrgMetadataEntry[]),
    ])
    organization = org
    users = usersResult.users ?? []
    projectCount = projCount
    applicationCount = appCount
    domains = domainList
    metadata = metadataList
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load organization"
    console.error("Failed to load organization:", e)
  }

  return (
    <OrgDetailClient
      organization={organization}
      orgId={orgId}
      users={users}
      projectCount={projectCount}
      applicationCount={applicationCount}
      initialDomains={domains}
      initialMetadata={metadata}
      error={error}
    />
  )
}
