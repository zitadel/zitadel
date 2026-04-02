"use server";

/**
 * Discover the current user's ZITADEL roles by calling the v1 auth
 * ListMyMemberships endpoint. This returns IAM/org/project memberships
 * with their role keys (e.g. IAM_OWNER, ORG_OWNER).
 *
 * TODO: Migrate to v2 InternalPermissionService.ListAdministrators
 * once the endpoint is available on target ZITADEL instances.
 */
export async function discoverUserRoles(): Promise<string[]> {
  // Resolve instance config: ZITADEL_INSTANCES first, then ZITADEL_INSTANCE_URL/PAT
  let baseUrl: string | undefined;
  let pat: string | undefined;

  try {
    const raw = process.env.ZITADEL_INSTANCES;
    if (raw) {
      const instances = JSON.parse(raw);
      if (instances.length > 0 && instances[0].url && instances[0].pat) {
        baseUrl = instances[0].url;
        pat = instances[0].pat;
      }
    }
  } catch {}

  if (!baseUrl || !pat) {
    baseUrl = process.env.ZITADEL_INSTANCE_URL;
    pat = process.env.ZITADEL_PAT;
  }

  if (!baseUrl || !pat) {
    console.warn("No ZITADEL instance configured — cannot discover roles");
    return [];
  }

  try {
    const response = await fetch(`${baseUrl}/auth/v1/memberships/me/_search`, {
      method: "POST",
      headers: {
        "Authorization": `Bearer ${pat}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({}),
      // Don't cache role discovery — roles can change
      cache: "no-store",
    });

    if (!response.ok) {
      console.error("Failed to discover roles:", response.status, response.statusText);
      return [];
    }

    const data = await response.json();
    const memberships = data.result ?? [];

    // Collect all unique role keys across all memberships
    const roles = new Set<string>();
    for (const membership of memberships) {
      if (Array.isArray(membership.roles)) {
        for (const role of membership.roles) {
          roles.add(role);
        }
      }
    }

    return Array.from(roles);
  } catch (e) {
    console.error("Failed to discover user roles:", e);
    return [];
  }
}
