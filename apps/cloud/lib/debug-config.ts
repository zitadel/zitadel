"use client"

/**
 * Debug configuration utilities.
 * 
 * Stores test instance configuration in localStorage + cookies.
 * localStorage for the client-side UI, cookies so server actions can read the config.
 * 
 * Only available when NEXT_PUBLIC_VERCEL_ENV === "preview" or NODE_ENV === "development".
 */

export interface DebugInstance {
  id: string
  name: string
  apiUrl: string
  token: string
  createdAt: number
}

export interface DebugConfig {
  enabled: boolean
  mode: "single" | "multi"
  activeInstanceId: string | null
  instances: DebugInstance[]
}

const STORAGE_KEY = "zitadel-debug-config"
const COOKIE_NAME = "zitadel-debug-instance"

const DEFAULT_CONFIG: DebugConfig = {
  enabled: false,
  mode: "single",
  activeInstanceId: null,
  instances: [],
}

export function isDebugAllowed(): boolean {
  if (typeof window === "undefined") return false
  const vercelEnv = process.env.NEXT_PUBLIC_VERCEL_ENV
  const nodeEnv = process.env.NODE_ENV
  return vercelEnv === "preview" || nodeEnv === "development"
}

export function getDebugConfig(): DebugConfig {
  if (typeof window === "undefined") return DEFAULT_CONFIG
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return DEFAULT_CONFIG
    return { ...DEFAULT_CONFIG, ...JSON.parse(raw) }
  } catch {
    return DEFAULT_CONFIG
  }
}

export function saveDebugConfig(config: DebugConfig): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(config))

  // Sync active instance to cookie for server actions
  const active = config.instances.find(i => i.id === config.activeInstanceId)
  if (active && config.enabled) {
    document.cookie = `${COOKIE_NAME}=${encodeURIComponent(JSON.stringify({
      apiUrl: active.apiUrl,
      token: active.token,
    }))}; path=/; samesite=lax`
  } else {
    document.cookie = `${COOKIE_NAME}=; path=/; max-age=0`
  }
}

export function addInstance(name: string, apiUrl: string, token: string): DebugConfig {
  const config = getDebugConfig()
  const instance: DebugInstance = {
    id: crypto.randomUUID(),
    name,
    apiUrl: apiUrl.replace(/\/$/, ""), // strip trailing slash
    token,
    createdAt: Date.now(),
  }
  config.instances.push(instance)
  if (!config.activeInstanceId) {
    config.activeInstanceId = instance.id
  }
  config.enabled = true
  saveDebugConfig(config)
  return config
}

export function removeInstance(id: string): DebugConfig {
  const config = getDebugConfig()
  config.instances = config.instances.filter(i => i.id !== id)
  if (config.activeInstanceId === id) {
    config.activeInstanceId = config.instances[0]?.id ?? null
  }
  if (config.instances.length === 0) {
    config.enabled = false
  }
  saveDebugConfig(config)
  return config
}

export function setActiveInstance(id: string): DebugConfig {
  const config = getDebugConfig()
  config.activeInstanceId = id
  saveDebugConfig(config)
  return config
}

export function setMode(mode: "single" | "multi"): DebugConfig {
  const config = getDebugConfig()
  config.mode = mode
  saveDebugConfig(config)
  return config
}

export function clearDebugConfig(): void {
  localStorage.removeItem(STORAGE_KEY)
  document.cookie = `${COOKIE_NAME}=; path=/; max-age=0`
}
