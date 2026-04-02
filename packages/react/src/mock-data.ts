import type {
  Instance,
  Organization,
  User,
  Project,
  Application,
  Administrator,
  Session,
  Action,
  ActivityLogEntry,
  RoleAssignment,
  AnalyticsData,
} from "./types"

// Seeded random number generator for deterministic data (mulberry32 algorithm)
// Uses only integer operations for consistent server/client results
function seededRandom(seed: number): number {
  let t = (seed + 0x6D2B79F5) | 0
  t = Math.imul(t ^ (t >>> 15), (t | 1)) | 0
  t = (t ^ (t + Math.imul(t ^ (t >>> 7), (t | 61)))) | 0
  return ((t ^ (t >>> 14)) >>> 0) / 4294967296
}

// Helper to generate deterministic dates based on index
function deterministicDate(start: Date, end: Date, seed: number): Date {
  const range = end.getTime() - start.getTime()
  return new Date(start.getTime() + seededRandom(seed) * range)
}

// Helper to pick element deterministically
function deterministicElement<T>(arr: T[], seed: number): T {
  return arr[Math.floor(seededRandom(seed) * arr.length)]
}

// Instance name prefixes and suffixes for realistic names
const instancePrefixes = ["Production", "Staging", "Development", "QA", "Demo", "Partner", "Customer", "Internal", "External", "Test", "UAT", "Sandbox", "Training", "Beta", "Alpha", "Integration", "Preview", "Release", "Canary", "Edge"]
const instanceSuffixes = ["", "US", "EU", "APAC", "Primary", "Secondary", "Backup", "Main", "Legacy", "New", "v2", "Core", "Hub", "Gateway", "Portal", "Platform", "Service", "System", "Cluster", "Node"]
const regions = ["EU (Frankfurt)", "US (Virginia)", "US (Oregon)", "Asia (Singapore)", "Asia (Tokyo)", "Australia (Sydney)", "EU (London)", "EU (Ireland)", "On-premise", "On-premise (K8s)"]
const versions = ["2.45.0", "2.44.2", "2.46.0-rc1", "2.43.1", "2.40.0", "2.45.1", "2.44.0", "2.42.3"]
const instanceStatuses: Instance["status"][] = ["active", "active", "active", "active", "active", "active", "active", "active", "inactive"]

// Generate 100 instances
export const instances: Instance[] = Array.from({ length: 100 }, (_, i) => {
  const prefixIndex = i % instancePrefixes.length
  const suffixIndex = Math.floor(i / instancePrefixes.length) % instanceSuffixes.length
  const prefix = instancePrefixes[prefixIndex]
  const suffix = instanceSuffixes[suffixIndex]
  const name = suffix ? `${prefix} ${suffix}` : prefix
  const isCloud = i % 3 !== 2 // ~67% cloud, ~33% self-hosted
  const region = isCloud ? regions[i % 8] : regions[8 + (i % 2)]
  const domainSuffix = isCloud ? "zitadel.cloud" : "internal"
  const domainPrefix = name.toLowerCase().replace(/\s+/g, "-").replace(/[()]/g, "")
  
  return {
    id: `inst-${i + 1}`,
    name,
    domain: `${domainPrefix}.${domainSuffix}`,
    createdAt: deterministicDate(new Date("2022-01-01"), new Date("2025-01-01"), i * 7),
    status: deterministicElement(instanceStatuses, i * 3),
    hostingType: isCloud ? "cloud" : "self-hosted",
    region,
    version: deterministicElement(versions, i * 5),
  }
})

// Generate 500 organizations with unique names - distributed across instances
const orgPrefixes = ["Acme", "Global", "Tech", "Digital", "Cloud", "Smart", "Innovate", "Next", "Prime", "Alpha", "Beta", "Gamma", "Delta", "Omega", "Sigma", "Phoenix", "Atlas", "Titan", "Nova", "Apex"]
const orgSuffixes = ["Corp", "Inc", "Ltd", "Solutions", "Systems", "Group", "Labs", "Industries", "Ventures", "Partners"]
const orgModifiers = ["", "North", "South", "East", "West", "Central"]

export const organizations: Organization[] = Array.from({ length: 500 }, (_, i) => {
  // Distribute orgs across instances - each instance gets ~5 organizations
  const instanceIndex = Math.floor(i / 5) % instances.length
  const instanceId = instances[instanceIndex].id
  const prefixIndex = i % orgPrefixes.length
  const suffixIndex = Math.floor(i / orgPrefixes.length) % orgSuffixes.length
  const modifierIndex = Math.floor(i / (orgPrefixes.length * orgSuffixes.length)) % orgModifiers.length
  const modifier = orgModifiers[modifierIndex]
  const baseName = `${orgPrefixes[prefixIndex]} ${orgSuffixes[suffixIndex]}`
  const name = modifier ? `${baseName} ${modifier}` : baseName
  
  // First org for each instance is the default
  const isDefault = i % 5 === 0
  
  return {
    id: `org-${i + 1}`,
    name,
    instanceId,
    isDefault,
    createdAt: deterministicDate(new Date("2023-01-01"), new Date("2025-01-01"), i * 11),
    userCount: Math.floor(seededRandom(i * 13) * 50) + 5,
    projectCount: Math.floor(seededRandom(i * 17) * 10) + 1,
    applicationCount: Math.floor(seededRandom(i * 19) * 15) + 1,
    adminCount: Math.floor(seededRandom(i * 23) * 3) + 1,
    status: deterministicElement(["active", "active", "active", "active", "active", "active", "active", "inactive"], i * 29) as Organization["status"],
  }
})

// Generate ~500 users across organizations
const firstNames = ["James", "Emma", "Oliver", "Sophia", "Liam", "Isabella", "Noah", "Mia", "Ethan", "Charlotte", "Lucas", "Amelia", "Mason", "Harper", "Logan", "Evelyn", "Alexander", "Abigail", "Sebastian", "Emily", "Michael", "Sarah", "David", "Anna", "Daniel", "Lisa", "Matthew", "Jessica", "Andrew", "Jennifer"]
const lastNames = ["Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson"]
const statuses: User["status"][] = ["active", "active", "active", "active", "inactive", "locked", "pending"]
const roles: User["role"][] = ["user", "user", "user", "admin", "owner"]

export const users: User[] = Array.from({ length: 500 }, (_, i) => {
  const firstName = firstNames[i % 30]
  const lastName = lastNames[Math.floor(i / 30) % 30]
  const org = organizations[i % 100]
  return {
    id: `user-${i + 1}`,
    email: `${firstName.toLowerCase()}.${lastName.toLowerCase()}${i}@${org.name.toLowerCase().replace(/\s/g, "")}.com`,
    firstName,
    lastName,
    displayName: `${firstName} ${lastName}`,
    username: `${firstName.toLowerCase()}.${lastName.toLowerCase()}${i}`,
    orgId: org.id,
    orgName: org.name,
    status: deterministicElement(statuses, i * 19),
    lastLogin: seededRandom(i * 23) > 0.2 ? deterministicDate(new Date("2024-10-01"), new Date("2025-03-01"), i * 29) : null,
    createdAt: deterministicDate(new Date("2023-01-01"), new Date("2025-01-01"), i * 31),
    role: deterministicElement(roles, i * 37),
  }
})

// Generate ~50 projects
const projectNames = ["Customer Portal", "Admin Dashboard", "Mobile App", "API Gateway", "Analytics Platform", "Billing System", "Notification Service", "User Management", "Content Management", "E-Commerce Platform", "Booking System", "Inventory Management", "HR Portal", "CRM System", "Support Desk", "Reporting Tool", "Integration Hub", "Data Pipeline", "Auth Service", "Marketplace"]

export const projects: Project[] = Array.from({ length: 50 }, (_, i) => {
  const org = organizations[i % 100]
  return {
    id: `proj-${i + 1}`,
    name: `${projectNames[i % 20]} ${Math.floor(i / 20) > 0 ? Math.floor(i / 20) + 1 : ""}`.trim(),
    description: `Project for ${org.name} - ${projectNames[i % 20]}`,
    orgId: org.id,
    orgName: org.name,
    createdAt: deterministicDate(new Date("2023-06-01"), new Date("2025-01-01"), i * 41),
    updatedAt: deterministicDate(new Date("2024-06-01"), new Date("2025-03-01"), i * 43),
    status: seededRandom(i * 47) > 0.1 ? "active" : "inactive",
    applicationCount: Math.floor(seededRandom(i * 53) * 5) + 1,
  }
})

// Generate ~100 applications
const appTypes: Application["type"][] = ["web", "native", "api", "user-agent"]

export const applications: Application[] = Array.from({ length: 100 }, (_, i) => {
  const project = projects[i % 50]
  const type = appTypes[i % 4]
  const typeName = type === "web" ? "Web App" : type === "native" ? "Mobile App" : type === "api" ? "API Client" : "Browser Extension"
  // Generate deterministic client ID
  const clientIdChars = "abcdefghijklmnopqrstuvwxyz0123456789"
  let clientId = "client_"
  for (let j = 0; j < 13; j++) {
    clientId += clientIdChars[Math.floor(seededRandom(i * 59 + j) * clientIdChars.length)]
  }
  return {
    id: `app-${i + 1}`,
    name: `${project.name} ${typeName}`,
    projectId: project.id,
    projectName: project.name,
    orgId: project.orgId,
    orgName: project.orgName,
    type,
    clientId,
    createdAt: deterministicDate(new Date("2023-08-01"), new Date("2025-01-01"), i * 61),
    status: seededRandom(i * 67) > 0.1 ? "active" : "inactive",
  }
})

// Generate administrators
export const administrators: Administrator[] = [
  { id: "admin-1", userId: "user-1", email: users[0].email, displayName: users[0].displayName, role: "iam-admin", grantedAt: new Date("2023-01-15") },
  { id: "admin-2", userId: "user-2", email: users[1].email, displayName: users[1].displayName, role: "instance-admin", grantedAt: new Date("2023-02-20") },
  { id: "admin-3", userId: "user-5", email: users[4].email, displayName: users[4].displayName, role: "instance-admin", grantedAt: new Date("2023-06-10") },
  { id: "admin-4", userId: "user-10", email: users[9].email, displayName: users[9].displayName, role: "org-admin", grantedAt: new Date("2024-01-05") },
  { id: "admin-5", userId: "user-15", email: users[14].email, displayName: users[14].displayName, role: "org-admin", grantedAt: new Date("2024-03-15") },
]

// Generate sessions
const userAgents = [
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0",
  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1",
  "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0) Mobile Safari/605.1",
  "Mozilla/5.0 (Linux; Android 14) Chrome/120.0",
]
const authMethods: Session["authMethod"][] = ["password", "passkey", "sso", "mfa"]

export const sessions: Session[] = Array.from({ length: 200 }, (_, i) => {
  const user = users[i % 500]
  const createdAt = deterministicDate(new Date("2025-01-01"), new Date("2025-03-01"), i * 71)
  const lastActivity = deterministicDate(createdAt, new Date("2025-03-15"), i * 73)
  return {
    id: `session-${i + 1}`,
    userId: user.id,
    userEmail: user.email,
    userName: user.displayName,
    orgId: user.orgId,
    orgName: user.orgName,
    ipAddress: `${Math.floor(seededRandom(i * 73) * 255)}.${Math.floor(seededRandom(i * 79) * 255)}.${Math.floor(seededRandom(i * 83) * 255)}.${Math.floor(seededRandom(i * 89) * 255)}`,
    userAgent: deterministicElement(userAgents, i * 97),
    authMethod: deterministicElement(authMethods, i * 99),
    createdAt,
    lastActivity,
    expiresAt: new Date(createdAt.getTime() + 24 * 60 * 60 * 1000),
    status: deterministicElement(["active", "active", "active", "expired", "revoked"], i * 101),
  }
})

// Generate actions
export const actions: Action[] = [
  { id: "action-1", name: "Pre-Auth Check", script: "function preAuth(ctx) { return true; }", description: "Check user before authentication", allowedToFail: false, timeout: 5000, createdAt: new Date("2024-01-15"), status: "active" },
  { id: "action-2", name: "Post-Login Hook", script: "function postLogin(ctx) { log(ctx.user); }", description: "Execute after successful login", allowedToFail: true, timeout: 3000, createdAt: new Date("2024-02-20"), status: "active" },
  { id: "action-3", name: "Token Enrichment", script: "function enrichToken(ctx) { ctx.claims.custom = true; }", description: "Add custom claims to tokens", allowedToFail: false, timeout: 2000, createdAt: new Date("2024-03-10"), status: "active" },
  { id: "action-4", name: "User Provisioning", script: "function provision(ctx) { createExternal(ctx.user); }", description: "Provision user to external systems", allowedToFail: true, timeout: 10000, createdAt: new Date("2024-04-05"), status: "inactive" },
  { id: "action-5", name: "Audit Logger", script: "function audit(ctx) { sendToSIEM(ctx); }", description: "Send events to SIEM", allowedToFail: true, timeout: 5000, createdAt: new Date("2024-05-15"), status: "active" },
]

// Generate activity log entries - fixed circular dependency
const actionTypes: ActivityLogEntry["actionType"][] = ["created", "updated", "deleted", "activated", "deactivated", "revoked", "assigned", "unassigned"]
const resourceTypesForLog: ActivityLogEntry["resourceType"][] = ["user", "project", "application", "organization", "settings", "role_assignment", "session"]

// Determine severity based on action and resource type
function getSeverity(actionType: string, resourceType: string): ActivityLogEntry["severity"] {
  // Sensitive: role escalations, deletions of important resources, security events
  if (actionType === "deleted" && ["organization", "project", "user"].includes(resourceType)) return "sensitive"
  if (actionType === "revoked") return "sensitive"
  if (resourceType === "role_assignment" && ["assigned", "deleted"].includes(actionType)) return "sensitive"
  if (resourceType === "session" && actionType === "revoked") return "sensitive"
  // Important: creation of resources, activations/deactivations
  if (actionType === "created" && ["user", "project", "application", "organization"].includes(resourceType)) return "important"
  if (["activated", "deactivated"].includes(actionType)) return "important"
  // Routine: updates, settings changes
  return "routine"
}

export const activityLog: ActivityLogEntry[] = Array.from({ length: 200 }, (_, i) => {
  const resourceType = deterministicElement(resourceTypesForLog, i * 103)
  const actionType = deterministicElement(actionTypes, i * 107)
  const actor = users[Math.floor(seededRandom(i * 109) * 50)]
  const org = organizations[Math.floor(seededRandom(i * 113) * 100)]
  const instance = instances.find(inst => inst.id === org.instanceId) || instances[0]
  
  let resourceName = ""
  let resourceId = ""
  
  switch (resourceType) {
    case "user":
      const user = users[Math.floor(seededRandom(i * 127) * 500)]
      resourceName = user.displayName
      resourceId = user.id
      break
    case "project":
      const project = projects[Math.floor(seededRandom(i * 131) * 50)]
      resourceName = project.name
      resourceId = project.id
      break
    case "application":
      const app = applications[Math.floor(seededRandom(i * 137) * 100)]
      resourceName = app.name
      resourceId = app.id
      break
    case "organization":
      resourceName = org.name
      resourceId = org.id
      break
    case "settings":
      resourceName = "Instance Settings"
      resourceId = "settings"
      break
    case "role_assignment":
      // Generate synthetic role assignment name without referencing the array
      const assignedUser = users[Math.floor(seededRandom(i * 141) * 500)]
      const assignedProject = projects[Math.floor(seededRandom(i * 143) * 50)]
      resourceName = `${assignedUser.displayName} → ${assignedProject.name}`
      resourceId = `role-${Math.floor(seededRandom(i * 145) * 150) + 1}`
      break
    case "session":
      // Generate synthetic session name without referencing the array
      const sessionUser = users[Math.floor(seededRandom(i * 147) * 500)]
      resourceName = `Session for ${sessionUser.displayName}`
      resourceId = `session-${Math.floor(seededRandom(i * 149) * 200) + 1}`
      break
  }
  
  return {
    id: `log-${i + 1}`,
    action: `${resourceType}.${actionType}`,
    actionType,
    resourceType,
    resourceId,
    resourceName,
    actorId: actor.id,
    actorName: actor.displayName,
    orgId: org.id,
    instanceId: instance.id,
    timestamp: deterministicDate(new Date("2024-12-01"), new Date("2025-03-15"), i * 139),
    severity: getSeverity(actionType, resourceType),
  }
}).sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime())

// Generate role assignments
export const roleAssignments: RoleAssignment[] = Array.from({ length: 150 }, (_, i) => {
  const user = users[i % 500]
  const project = projects[Math.floor(seededRandom(i * 149) * 50)]
  const org = organizations.find(o => o.id === project.orgId) || organizations[0]
  const grantor = users[Math.floor(seededRandom(i * 167) * 10)]
  const possibleRoles = ["viewer", "editor", "admin", "owner", "developer", "tester", "analyst"]
  const roleCount = Math.floor(seededRandom(i * 151) * 3) + 1
  const assignedRoles: string[] = []
  for (let j = 0; j < roleCount; j++) {
    const role = deterministicElement(possibleRoles, i * 157 + j)
    if (!assignedRoles.includes(role)) {
      assignedRoles.push(role)
    }
  }
  
  return {
    id: `role-${i + 1}`,
    userId: user.id,
    userName: user.displayName,
    userEmail: user.email,
    projectId: project.id,
    projectName: project.name,
    orgId: org.id,
    orgName: org.name,
    roles: assignedRoles,
    grantedAt: deterministicDate(new Date("2024-01-01"), new Date("2025-03-01"), i * 163),
    grantedBy: grantor.displayName,
    grantedById: grantor.id,
  }
})

// Generate analytics data (last 30 days) - use fixed base date for determinism
const baseDate = new Date("2025-03-17")
export const analyticsData: AnalyticsData[] = Array.from({ length: 30 }, (_, i) => {
  const date = new Date(baseDate)
  date.setDate(date.getDate() - (29 - i))
  return {
    date: date.toISOString().split("T")[0],
    apiRequests: Math.floor(seededRandom(i * 173) * 50000) + 10000,
    activeUsers: Math.floor(seededRandom(i * 179) * 300) + 100,
    newUsers: Math.floor(seededRandom(i * 181) * 20) + 5,
    sessions: Math.floor(seededRandom(i * 191) * 500) + 200,
  }
})

// Helper functions to filter data
export function getOrganizationsByInstance(instanceId: string): Organization[] {
  return organizations.filter(org => org.instanceId === instanceId)
}

export function getUsersByOrganization(orgId: string): User[] {
  return users.filter(user => user.orgId === orgId)
}

export function getProjectsByOrganization(orgId: string): Project[] {
  return projects.filter(project => project.orgId === orgId)
}

export function getApplicationsByOrganization(orgId: string): Application[] {
  return applications.filter(app => app.orgId === orgId)
}

export function getApplicationsByProject(projectId: string): Application[] {
  return applications.filter(app => app.projectId === projectId)
}

export function getSessionsByOrganization(orgId: string): Session[] {
  return sessions.filter(session => session.orgId === orgId)
}

export function getActivityLogByOrganization(orgId: string): ActivityLogEntry[] {
  return activityLog.filter(entry => entry.orgId === orgId)
}

export function getActivityLogByInstance(instanceId: string): ActivityLogEntry[] {
  return activityLog.filter(entry => entry.instanceId === instanceId)
}

export function getRoleAssignmentsByOrganization(orgId: string): RoleAssignment[] {
  const orgProjects = getProjectsByOrganization(orgId)
  const projectIds = new Set(orgProjects.map(p => p.id))
  return roleAssignments.filter(ra => projectIds.has(ra.projectId))
}
