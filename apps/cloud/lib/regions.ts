/**
 * Region configuration — reads available regions from env vars.
 * Follows the same pattern as the website: ZITADEL_PUBLIC_REGION_<id> and ZITADEL_PRIVATE_REGION_<id>.
 * 
 * Copied from website/src/services/region.ts and adapted for the cloud app.
 */

const PUBLIC_REGION_ENV_PREFIX = 'ZITADEL_PUBLIC_REGION_'
const PRIVATE_REGION_ENV_PREFIX = 'ZITADEL_PRIVATE_REGION_'

export type PublicRegion = {
  id: string
  displayName: string
  displayIdx: number
  isoCode: string
  isPreview: boolean
  isDefault: boolean
  priceInCents: number
}

export type PrivateRegion = {
  id: string
  systemApi: {
    target: string
    audience: string
    username: string
    key: string
  }
  stripe: {
    priceId: string
  }
}

export function readPublicRegion(regionId: string): PublicRegion {
  const envKey = PUBLIC_REGION_ENV_PREFIX + regionId
  const envValue = process.env[envKey]
  if (!envValue) throw new Error(`Region not found: ${envKey}`)
  return { ...JSON.parse(envValue), id: regionId }
}

export function readPrivateRegion(regionId: string): PrivateRegion {
  const envKey = PRIVATE_REGION_ENV_PREFIX + regionId
  const envValue = process.env[envKey]
  if (!envValue) throw new Error(`Region not found: ${envKey}`)
  return { ...JSON.parse(envValue), id: regionId }
}

export function availablePublicRegions(): PublicRegion[] {
  return Object.keys(process.env)
    .filter((key) => key.startsWith(PUBLIC_REGION_ENV_PREFIX))
    .map((key) => readPublicRegion(key.replace(PUBLIC_REGION_ENV_PREFIX, '')))
    .sort((a, b) => a.displayIdx - b.displayIdx)
}

export function availablePrivateRegions(): PrivateRegion[] {
  return Object.keys(process.env)
    .filter((key) => key.startsWith(PRIVATE_REGION_ENV_PREFIX))
    .map((key) => readPrivateRegion(key.replace(PRIVATE_REGION_ENV_PREFIX, '')))
}
