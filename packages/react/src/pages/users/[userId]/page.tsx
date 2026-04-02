import { toJson, create } from "@zitadel/client"
import {
  GetUserByIDRequestSchema,
  GetUserByIDResponseSchema,
  UserService,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb"
import { createClient } from "@connectrpc/connect"
import { getTransport } from "../../../api/transport"
import { listUserSessions } from "../../../api/sessions"
import { listAuthMethodTypes, listPasskeys, listAuthFactors } from "../../../api/user-security"
import { listUserMetadata } from "../../../api/user-metadata"
import { UserDetailClient } from "./user-detail-client"

interface UserDetailPageProps {
  params: Promise<{ userId: string }>
}

/**
 * User detail page — server component that fetches a user by ID,
 * sessions, auth methods, passkeys, and metadata in parallel.
 */
export default async function UserDetailPage({ params }: UserDetailPageProps) {
  const { userId } = await params
  let user: any = null
  let initialSessions: any[] = []
  let totalSessions = 0
  let authMethods: string[] = []
  let authFactors: any[] = []
  let passkeys: any[] = []
  let metadata: any[] = []
  let error: string | null = null

  try {
    const transport = getTransport()
    const userClient = createClient(UserService, transport)
    const userRequest = create(GetUserByIDRequestSchema, { userId })

    const [userResponse, sessionResult, authMethodsResult, authFactorsResult, passkeysResult, metadataResult] = await Promise.all([
      userClient.getUserByID(userRequest),
      listUserSessions(userId, 10, 0).catch((e) => {
        console.error("Failed to load sessions:", e)
        return { sessions: [], totalResult: 0 }
      }),
      listAuthMethodTypes(userId).catch((e) => {
        console.error("Failed to load auth methods:", e)
        return [] as string[]
      }),
      listAuthFactors(userId).catch((e) => {
        console.error("Failed to load auth factors:", e)
        return [] as any[]
      }),
      listPasskeys(userId).catch((e) => {
        console.error("Failed to load passkeys:", e)
        return [] as any[]
      }),
      listUserMetadata(userId).catch((e) => {
        console.error("Failed to load metadata:", e)
        return [] as any[]
      }),
    ])

    const userJson = toJson(GetUserByIDResponseSchema, userResponse) as any
    user = userJson.user ?? null
    initialSessions = sessionResult.sessions
    totalSessions = sessionResult.totalResult
    authMethods = authMethodsResult
    authFactors = authFactorsResult
    passkeys = passkeysResult
    metadata = metadataResult
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load user"
    console.error("Failed to load user:", e)
  }

  return (
    <UserDetailClient
      user={user}
      userId={userId}
      initialSessions={initialSessions}
      totalSessions={totalSessions}
      initialAuthMethods={authMethods}
      initialAuthFactors={authFactors}
      initialPasskeys={passkeys}
      initialMetadata={metadata}
      error={error}
    />
  )
}
