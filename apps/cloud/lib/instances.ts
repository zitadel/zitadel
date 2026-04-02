import type { TestInstance } from "./service-configs"

/**
 * Instance management — reads configured instances from env vars.
 * Each instance gets a stable ID derived from its URL hostname.
 */

export interface Instance extends TestInstance {
  id: string
}

/** Generate a stable ID from a URL (uses hostname) */
function instanceId(url: string): string {
  try {
    return new URL(url).hostname.replace(/\./g, "-")
  } catch {
    return url.replace(/[^a-z0-9]/gi, "-").toLowerCase()
  }
}

/** Read all configured instances from env */
export function getInstances(): Instance[] {
  const instances: Instance[] = []

  try {
    const raw = process.env.ZITADEL_INSTANCES
    if (raw) {
      const parsed = JSON.parse(raw) as TestInstance[]
      for (const inst of parsed) {
        if (inst.url) {
          instances.push({ ...inst, id: instanceId(inst.url) })
        }
      }
    }
  } catch {}

  // Fallback: single instance from legacy env vars
  if (instances.length === 0) {
    const url = process.env.ZITADEL_INSTANCE_URL
    if (url) {
      instances.push({
        id: instanceId(url),
        name: new URL(url).hostname,
        url,
        pat: process.env.ZITADEL_PAT ?? "",
      })
    }
  }

  return instances
}

/** Get a single instance by ID */
export function getInstance(id: string): Instance | null {
  return getInstances().find((i) => i.id === id) ?? null
}
