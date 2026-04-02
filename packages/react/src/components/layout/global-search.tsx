"use client"

import * as React from "react"
import { useConsoleRouter as useRouter } from "../../hooks/use-console-router"
import {
  Search,
  Building2,
  Users,
  FolderKanban,
  AppWindow,
  KeyRound,
  Zap,
  User,
  Plus,
  Shield,
  Palette,
  Link2,
  Key,
  BookOpen,
  FileText,
  MessageCircle,
  Github,
  LifeBuoy,
  Activity,
  History,
  Fingerprint,
  Globe,
  Settings,
  ExternalLink,
  Loader2,
} from "lucide-react"
import { Button } from "../ui/button"
import { Badge } from "../ui/badge"
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "../ui/command"
import { searchAll, type SearchResult } from "../../api/search-resources"

// Resource type icons
const typeIcons: Record<string, React.ReactNode> = {
  user: <User className="mr-2 h-4 w-4" />,
  org: <Building2 className="mr-2 h-4 w-4" />,
  project: <FolderKanban className="mr-2 h-4 w-4" />,
  app: <AppWindow className="mr-2 h-4 w-4" />,
}

// Scope prefixes the user can type to filter to a single resource type
const scopePrefixes: { prefix: string; scope: "user" | "org" | "project" | "app"; label: string }[] = [
  { prefix: "user:", scope: "user", label: "Users" },
  { prefix: "org:", scope: "org", label: "Organizations" },
  { prefix: "project:", scope: "project", label: "Projects" },
  { prefix: "app:", scope: "app", label: "Applications" },
]

export function GlobalSearch() {
  const router = useRouter()
  const [open, setOpen] = React.useState(false)
  const [rawQuery, setRawQuery] = React.useState("")
  const [results, setResults] = React.useState<{
    users: SearchResult[]
    organizations: SearchResult[]
    projects: SearchResult[]
    applications: SearchResult[]
  }>({ users: [], organizations: [], projects: [], applications: [] })
  const [isSearching, setIsSearching] = React.useState(false)
  const debounceRef = React.useRef<ReturnType<typeof setTimeout> | null>(null)

  // Parse scope prefix from raw query
  const { scope, query } = React.useMemo(() => {
    const lower = rawQuery.toLowerCase()
    for (const sp of scopePrefixes) {
      if (lower.startsWith(sp.prefix)) {
        return { scope: sp.scope, query: rawQuery.slice(sp.prefix.length).trim() }
      }
    }
    return { scope: undefined as "user" | "org" | "project" | "app" | undefined, query: rawQuery.trim() }
  }, [rawQuery])

  // Debounced search
  React.useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current)

    if (!query && !scope) {
      setResults({ users: [], organizations: [], projects: [], applications: [] })
      setIsSearching(false)
      return
    }

    setIsSearching(true)
    debounceRef.current = setTimeout(async () => {
      try {
        const res = await searchAll(query, scope)
        setResults(res)
      } catch {
        setResults({ users: [], organizations: [], projects: [], applications: [] })
      } finally {
        setIsSearching(false)
      }
    }, 300)

    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current)
    }
  }, [query, scope])

  // Keyboard shortcut to open search
  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        setOpen((o) => !o)
      }
    }
    document.addEventListener("keydown", down)
    return () => document.removeEventListener("keydown", down)
  }, [])

  // Reset on close
  React.useEffect(() => {
    if (!open) {
      setRawQuery("")
      setResults({ users: [], organizations: [], projects: [], applications: [] })
    }
  }, [open])

  const handleSelect = (callback: () => void) => {
    setOpen(false)
    callback()
  }

  const navigateToResult = (result: SearchResult) => {
    switch (result.type) {
      case "user":
        router.push(`/users/${result.id}`)
        break
      case "org":
        router.push(`/organizations/${result.id}`)
        break
      case "project":
        router.push(`/projects/${result.id}`)
        break
      case "app":
        router.push(`/applications/${result.id}`)
        break
    }
  }

  const hasResults =
    results.users.length > 0 ||
    results.organizations.length > 0 ||
    results.projects.length > 0 ||
    results.applications.length > 0

  const showStaticActions = !query && !scope

  return (
    <>
      <Button
        variant="outline"
        className="relative h-9 w-full max-w-sm justify-start text-sm text-muted-foreground sm:pr-12"
        onClick={() => setOpen(true)}
      >
        <Search className="mr-2 h-4 w-4" />
        <span className="hidden lg:inline-flex">Search or type a command...</span>
        <span className="inline-flex lg:hidden">Search...</span>
        <kbd className="pointer-events-none absolute right-1.5 top-1.5 hidden h-6 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Button>

      <CommandDialog open={open} onOpenChange={setOpen} shouldFilter={false}>
        <CommandInput
          placeholder={scope ? `Search ${scope}s...` : "Search users, orgs, projects, apps... (try user: org: project: app:)"}
          value={rawQuery}
          onValueChange={setRawQuery}
        />
        <CommandList className="max-h-[400px]">
          {/* Loading indicator */}
          {isSearching && (
            <div className="flex items-center justify-center py-6 text-sm text-muted-foreground">
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Searching...
            </div>
          )}

          {/* No results */}
          {!isSearching && query && !hasResults && (
            <CommandEmpty>No results found for &ldquo;{query}&rdquo;</CommandEmpty>
          )}

          {/* Scope hints when empty */}
          {showStaticActions && (
            <>
              <CommandGroup heading="FILTER BY TYPE">
                {scopePrefixes.map((sp) => (
                  <CommandItem
                    key={sp.prefix}
                    onSelect={() => setRawQuery(sp.prefix)}
                    value={sp.prefix}
                  >
                    {typeIcons[sp.scope]}
                    <span>{sp.prefix}</span>
                    <span className="ml-2 text-muted-foreground text-xs">Search {sp.label.toLowerCase()}</span>
                    <Badge variant="outline" className="ml-auto text-[10px]">filter</Badge>
                  </CommandItem>
                ))}
              </CommandGroup>
              <CommandSeparator />
            </>
          )}

          {/* Search Results */}
          {!isSearching && results.users.length > 0 && (
            <>
              <CommandGroup heading="USERS">
                {results.users.map((r) => (
                  <CommandItem
                    key={`user-${r.id}`}
                    onSelect={() => handleSelect(() => navigateToResult(r))}
                    value={`user-${r.id}-${r.name}`}
                  >
                    <User className="mr-2 h-4 w-4" />
                    <span>{r.name}</span>
                    {r.description && (
                      <span className="ml-auto text-xs text-muted-foreground">{r.description}</span>
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>
              <CommandSeparator />
            </>
          )}

          {!isSearching && results.organizations.length > 0 && (
            <>
              <CommandGroup heading="ORGANIZATIONS">
                {results.organizations.map((r) => (
                  <CommandItem
                    key={`org-${r.id}`}
                    onSelect={() => handleSelect(() => navigateToResult(r))}
                    value={`org-${r.id}-${r.name}`}
                  >
                    <Building2 className="mr-2 h-4 w-4" />
                    <span>{r.name}</span>
                    {r.description && (
                      <span className="ml-auto text-xs text-muted-foreground">{r.description}</span>
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>
              <CommandSeparator />
            </>
          )}

          {!isSearching && results.projects.length > 0 && (
            <>
              <CommandGroup heading="PROJECTS">
                {results.projects.map((r) => (
                  <CommandItem
                    key={`project-${r.id}`}
                    onSelect={() => handleSelect(() => navigateToResult(r))}
                    value={`project-${r.id}-${r.name}`}
                  >
                    <FolderKanban className="mr-2 h-4 w-4" />
                    <span>{r.name}</span>
                    {r.description && (
                      <span className="ml-auto text-xs text-muted-foreground">{r.description}</span>
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>
              <CommandSeparator />
            </>
          )}

          {!isSearching && results.applications.length > 0 && (
            <>
              <CommandGroup heading="APPLICATIONS">
                {results.applications.map((r) => (
                  <CommandItem
                    key={`app-${r.id}`}
                    onSelect={() => handleSelect(() => navigateToResult(r))}
                    value={`app-${r.id}-${r.name}`}
                  >
                    <AppWindow className="mr-2 h-4 w-4" />
                    <span>{r.name}</span>
                    {r.description && (
                      <Badge variant="secondary" className="ml-auto text-[10px]">{r.description}</Badge>
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>
              <CommandSeparator />
            </>
          )}

          {/* Static actions — always available */}
          {showStaticActions && (
            <>
              <CommandGroup heading="NAVIGATION">
                <CommandItem onSelect={() => handleSelect(() => router.push("/users"))}>
                  <Users className="mr-2 h-4 w-4" />
                  <span>Users</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/organizations"))}>
                  <Building2 className="mr-2 h-4 w-4" />
                  <span>Organizations</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/projects"))}>
                  <FolderKanban className="mr-2 h-4 w-4" />
                  <span>Projects</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/applications"))}>
                  <AppWindow className="mr-2 h-4 w-4" />
                  <span>Applications</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/sessions"))}>
                  <KeyRound className="mr-2 h-4 w-4" />
                  <span>Sessions</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/actions"))}>
                  <Zap className="mr-2 h-4 w-4" />
                  <span>Actions</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/settings"))}>
                  <Settings className="mr-2 h-4 w-4" />
                  <span>Settings</span>
                </CommandItem>
              </CommandGroup>

              <CommandSeparator />

              <CommandGroup heading="QUICK ACTIONS">
                <CommandItem onSelect={() => handleSelect(() => router.push("/users?action=create"))}>
                  <Plus className="mr-2 h-4 w-4" />
                  <span>Create User</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/applications?action=create"))}>
                  <Plus className="mr-2 h-4 w-4" />
                  <span>Create Application</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/settings/idps"))}>
                  <Globe className="mr-2 h-4 w-4" />
                  <span>Identity Providers</span>
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => router.push("/settings/branding"))}>
                  <Palette className="mr-2 h-4 w-4" />
                  <span>Customize Login UI</span>
                </CommandItem>
              </CommandGroup>

              <CommandSeparator />

              <CommandGroup heading="HELP">
                <CommandItem onSelect={() => handleSelect(() => window.open("https://zitadel.com/docs", "_blank"))}>
                  <BookOpen className="mr-2 h-4 w-4" />
                  <span>Documentation</span>
                  <ExternalLink className="ml-auto h-3 w-3 text-muted-foreground" />
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => window.open("https://zitadel.com/docs/apis/introduction", "_blank"))}>
                  <FileText className="mr-2 h-4 w-4" />
                  <span>API Reference</span>
                  <ExternalLink className="ml-auto h-3 w-3 text-muted-foreground" />
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => window.open("https://zitadel.com/chat", "_blank"))}>
                  <MessageCircle className="mr-2 h-4 w-4" />
                  <span>Discord Community</span>
                  <ExternalLink className="ml-auto h-3 w-3 text-muted-foreground" />
                </CommandItem>
                <CommandItem onSelect={() => handleSelect(() => window.open("https://github.com/zitadel/zitadel", "_blank"))}>
                  <Github className="mr-2 h-4 w-4" />
                  <span>GitHub</span>
                  <ExternalLink className="ml-auto h-3 w-3 text-muted-foreground" />
                </CommandItem>
              </CommandGroup>
            </>
          )}
        </CommandList>
      </CommandDialog>
    </>
  )
}
