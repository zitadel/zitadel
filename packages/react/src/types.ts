export interface Instance {
  id: string
  name: string
  domain: string
  createdAt: Date
  status: "active" | "inactive"
  hostingType: "cloud" | "self-hosted"
  region?: string
  version?: string
}

export interface Organization {
  id: string
  name: string
  instanceId: string
  isDefault: boolean
  createdAt: Date
  userCount: number
  projectCount: number
  applicationCount: number
  adminCount: number
  status: "active" | "inactive"
}

export interface User {
  id: string
  email: string
  username: string
  firstName: string
  lastName: string
  displayName: string
  orgId: string
  orgName: string
  status: "active" | "inactive" | "locked" | "pending"
  lastLogin: Date | null
  createdAt: Date
  avatar?: string
  role: "user" | "admin" | "owner"
}

export interface Project {
  id: string
  name: string
  description: string
  orgId: string
  orgName: string
  createdAt: Date
  updatedAt: Date
  status: "active" | "inactive"
  applicationCount: number
}

export interface Application {
  id: string
  name: string
  projectId: string
  projectName: string
  orgId: string
  orgName: string
  type: "web" | "native" | "api" | "user-agent"
  clientId: string
  createdAt: Date
  status: "active" | "inactive"
}

export interface Administrator {
  id: string
  userId: string
  email: string
  displayName: string
  role: "instance-admin" | "iam-admin" | "org-admin"
  grantedAt: Date
}

export interface Session {
  id: string
  userId: string
  userEmail: string
  userName: string
  orgId: string
  orgName: string
  ipAddress: string
  userAgent: string
  authMethod: "password" | "passkey" | "sso" | "mfa"
  createdAt: Date
  lastActivity: Date
  expiresAt: Date
  status: "active" | "expired" | "revoked"
}

export interface Action {
  id: string
  name: string
  script: string
  description: string
  allowedToFail: boolean
  timeout: number
  createdAt: Date
  status: "active" | "inactive"
}

export interface ActivityLogEntry {
  id: string
  action: string
  actionType: "created" | "updated" | "deleted" | "activated" | "deactivated" | "revoked" | "assigned" | "unassigned"
  resourceType: "user" | "project" | "application" | "organization" | "settings" | "role_assignment" | "session"
  resourceId: string
  resourceName: string
  actorId: string
  actorName: string
  orgId: string
  instanceId: string
  timestamp: Date
  details?: string
  severity: "routine" | "important" | "sensitive"
}

export interface RoleAssignment {
  id: string
  userId: string
  userName: string
  userEmail: string
  projectId: string
  projectName: string
  orgId: string
  orgName: string
  roles: string[]
  grantedAt: Date
  grantedBy: string
  grantedById: string
}

export interface AnalyticsData {
  date: string
  apiRequests: number
  activeUsers: number
  newUsers: number
  sessions: number
}
