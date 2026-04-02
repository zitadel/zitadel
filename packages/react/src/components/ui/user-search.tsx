"use client"

import * as React from "react"
import { Search, User, Building2, Loader2, X } from "lucide-react"
import { Input } from "./input"
import { cn } from "../../utils"
import { searchUsers, type UserSearchResult } from "../../api/search-users"

interface UserSearchProps {
  /** Called when a user is selected */
  onSelect: (user: UserSearchResult) => void
  /** Placeholder text */
  placeholder?: string
  /** Additional CSS classes */
  className?: string
  /** Maximum results to show */
  limit?: number
  /** If set, disables the search */
  disabled?: boolean
}

/**
 * Reusable user search component backed by the ListUsers v2 API.
 * Searches across username, display name, and email with debounce.
 * Shows username and organization in results.
 */
export function UserSearch({
  onSelect,
  placeholder = "Search users by name, email, or username...",
  className,
  limit = 8,
  disabled = false,
}: UserSearchProps) {
  const [query, setQuery] = React.useState("")
  const [results, setResults] = React.useState<UserSearchResult[]>([])
  const [isLoading, setIsLoading] = React.useState(false)
  const [isOpen, setIsOpen] = React.useState(false)
  const [selectedIndex, setSelectedIndex] = React.useState(-1)
  const containerRef = React.useRef<HTMLDivElement>(null)
  const inputRef = React.useRef<HTMLInputElement>(null)
  const debounceRef = React.useRef<ReturnType<typeof setTimeout>>(undefined)

  // Debounced search
  React.useEffect(() => {
    if (!query.trim()) {
      setResults([])
      setIsOpen(false)
      return
    }

    setIsLoading(true)
    clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(async () => {
      try {
        const users = await searchUsers(query.trim(), limit)
        setResults(users)
        setIsOpen(users.length > 0)
        setSelectedIndex(-1)
      } catch (e) {
        console.error("User search failed:", e)
        setResults([])
      } finally {
        setIsLoading(false)
      }
    }, 300)

    return () => clearTimeout(debounceRef.current)
  }, [query, limit])

  // Close on outside click
  React.useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener("mousedown", handleClick)
    return () => document.removeEventListener("mousedown", handleClick)
  }, [])

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!isOpen || results.length === 0) return

    switch (e.key) {
      case "ArrowDown":
        e.preventDefault()
        setSelectedIndex((i) => (i < results.length - 1 ? i + 1 : 0))
        break
      case "ArrowUp":
        e.preventDefault()
        setSelectedIndex((i) => (i > 0 ? i - 1 : results.length - 1))
        break
      case "Enter":
        e.preventDefault()
        if (selectedIndex >= 0 && selectedIndex < results.length) {
          handleSelect(results[selectedIndex])
        }
        break
      case "Escape":
        setIsOpen(false)
        break
    }
  }

  const handleSelect = (user: UserSearchResult) => {
    onSelect(user)
    setQuery("")
    setResults([])
    setIsOpen(false)
  }

  return (
    <div ref={containerRef} className={cn("relative", className)}>
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          ref={inputRef}
          placeholder={placeholder}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => results.length > 0 && setIsOpen(true)}
          onKeyDown={handleKeyDown}
          className="pl-9 pr-9"
          disabled={disabled}
        />
        {isLoading && (
          <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 animate-spin text-muted-foreground" />
        )}
        {!isLoading && query && (
          <button
            onClick={() => { setQuery(""); setResults([]); setIsOpen(false) }}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        )}
      </div>

      {/* Results dropdown */}
      {isOpen && results.length > 0 && (
        <div className="absolute z-50 mt-1 w-full rounded-md border bg-popover shadow-md animate-in fade-in-0 zoom-in-95">
          <ul className="max-h-[280px] overflow-auto p-1" role="listbox">
            {results.map((user, index) => (
              <li
                key={user.userId}
                role="option"
                aria-selected={index === selectedIndex}
                className={cn(
                  "flex items-center gap-3 rounded-sm px-3 py-2 cursor-pointer text-sm transition-colors",
                  index === selectedIndex
                    ? "bg-accent text-accent-foreground"
                    : "hover:bg-accent/50"
                )}
                onClick={() => handleSelect(user)}
                onMouseEnter={() => setSelectedIndex(index)}
              >
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 shrink-0">
                  <User className="h-4 w-4 text-primary" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="font-medium truncate">{user.displayName}</span>
                    <span className="text-xs text-muted-foreground truncate">
                      @{user.username}
                    </span>
                  </div>
                  <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                    {user.email && <span className="truncate">{user.email}</span>}
                    {user.email && user.organizationId && <span>·</span>}
                    {user.organizationId && (
                      <span className="flex items-center gap-1 shrink-0">
                        <Building2 className="h-3 w-3" />
                        {user.organizationId}
                      </span>
                    )}
                  </div>
                </div>
              </li>
            ))}
          </ul>
          {results.length === limit && (
            <div className="border-t px-3 py-1.5 text-center text-xs text-muted-foreground">
              Showing first {limit} results — refine your search
            </div>
          )}
        </div>
      )}

      {/* No results */}
      {isOpen && query.trim() && results.length === 0 && !isLoading && (
        <div className="absolute z-50 mt-1 w-full rounded-md border bg-popover shadow-md p-4 text-center">
          <p className="text-sm text-muted-foreground">No users found</p>
        </div>
      )}
    </div>
  )
}
