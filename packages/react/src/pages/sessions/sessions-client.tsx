"use client"

import { useState, useTransition } from "react"
import { ConsoleLink as Link } from "../../context/link-context"
import { Badge } from "../../components/ui/badge"
import { Button } from "../../components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card"
import {
  KeyRound,
  Monitor,
  Globe,
  Clock,
  Fingerprint,
  ChevronLeft,
  ChevronRight,
  Loader2,
  User,
} from "lucide-react"
import { fetchAllSessions } from "../../api/all-sessions"

interface SessionsClientProps {
  initialSessions: any[]
  totalSessions: number
  error: string | null
}

const PAGE_SIZE = 20

function formatDate(dateStr?: string) {
  if (!dateStr) return "—"
  return new Date(dateStr).toLocaleString()
}

function getSessionFactors(session: any) {
  const factors = session.factors ?? {}
  const items: string[] = []
  if (factors.user) items.push("User")
  if (factors.password) items.push("Password")
  if (factors.webAuthN) items.push("WebAuthn")
  if (factors.totp) items.push("TOTP")
  if (factors.otpSms) items.push("OTP SMS")
  if (factors.otpEmail) items.push("OTP Email")
  if (factors.intent) items.push("Intent")
  return items
}

function getSessionUser(session: any) {
  const user = session.factors?.user?.user ?? {}
  return {
    id: session.factors?.user?.id ?? user.id ?? "",
    displayName: user.displayName ?? user.loginName ?? "Unknown",
    loginName: user.loginName ?? "",
  }
}

export function SessionsClient({ initialSessions, totalSessions, error }: SessionsClientProps) {
  const [sessions, setSessions] = useState(initialSessions)
  const [total, setTotal] = useState(totalSessions)
  const [page, setPage] = useState(0)
  const [isLoading, startTransition] = useTransition()

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

  function goToPage(newPage: number) {
    startTransition(async () => {
      try {
        const result = await fetchAllSessions(PAGE_SIZE, newPage * PAGE_SIZE)
        setSessions(result.sessions)
        setTotal(result.totalResult)
        setPage(newPage)
      } catch (e) {
        console.error("Failed to load sessions page:", e)
      }
    })
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Sessions</h1>
          <p className="text-sm text-muted-foreground">Active authentication sessions</p>
        </div>
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-6 text-center">
          <p className="text-sm font-medium text-destructive">Failed to load sessions</p>
          <p className="text-xs text-muted-foreground mt-1">{error}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">
          Sessions {isLoading && <Loader2 className="inline h-5 w-5 animate-spin ml-2" />}
        </h1>
        <p className="text-sm text-muted-foreground">
          {total} authentication session{total !== 1 ? "s" : ""}
        </p>
      </div>

      {/* Session list */}
      <div className="space-y-3">
        {sessions.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <KeyRound className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <p className="text-muted-foreground">No sessions found</p>
            </CardContent>
          </Card>
        ) : (
          sessions.map((session: any) => {
            const factors = getSessionFactors(session)
            const userInfo = getSessionUser(session)
            const userAgent = session.userAgent ?? {}
            const headerName = userAgent.header?.name ?? ""

            return (
              <Card key={session.id}>
                <CardContent className="p-4">
                  <div className="flex items-start justify-between">
                    <div className="space-y-2 flex-1 min-w-0">
                      {/* User info */}
                      <div className="flex items-center gap-2">
                        <User className="h-4 w-4 text-muted-foreground shrink-0" />
                        {userInfo.id ? (
                          <Link
                            href={`/users/${userInfo.id}`}
                            className="font-medium text-sm hover:underline truncate"
                          >
                            {userInfo.displayName}
                          </Link>
                        ) : (
                          <span className="font-medium text-sm truncate">
                            {userInfo.displayName}
                          </span>
                        )}
                        {userInfo.loginName && (
                          <span className="text-xs text-muted-foreground truncate">
                            ({userInfo.loginName})
                          </span>
                        )}
                      </div>

                      {/* Session ID */}
                      <div className="flex items-center gap-2">
                        <Monitor className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                        <span className="text-xs font-mono text-muted-foreground truncate">
                          {session.id}
                        </span>
                      </div>

                      {/* User agent */}
                      {headerName && (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Globe className="h-3.5 w-3.5 shrink-0" />
                          <span className="truncate">{headerName}</span>
                        </div>
                      )}

                      {/* Dates */}
                      <div className="flex items-center gap-4 text-xs text-muted-foreground">
                        <span className="flex items-center gap-1">
                          <Clock className="h-3 w-3" />
                          Created {formatDate(session.creationDate)}
                        </span>
                        {session.expirationDate && (
                          <span className="flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            Expires {formatDate(session.expirationDate)}
                          </span>
                        )}
                      </div>

                      {/* Auth factors */}
                      {factors.length > 0 && (
                        <div className="flex items-center gap-2 pt-1">
                          <Fingerprint className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                          <div className="flex flex-wrap gap-1">
                            {factors.map((f) => (
                              <Badge key={f} variant="secondary" className="text-xs">
                                {f}
                              </Badge>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>
            )
          })
        )}
      </div>

      {/* Pagination */}
      {total > PAGE_SIZE && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            Page {page + 1} of {totalPages} · {total} sessions
          </p>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => goToPage(page - 1)}
              disabled={page === 0 || isLoading}
            >
              <ChevronLeft className="h-4 w-4 mr-1" />
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => goToPage(page + 1)}
              disabled={page >= totalPages - 1 || isLoading}
            >
              Next
              <ChevronRight className="h-4 w-4 ml-1" />
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
