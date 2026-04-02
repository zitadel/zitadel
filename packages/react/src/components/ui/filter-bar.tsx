"use client"

import * as React from "react"
import { Search, X } from "lucide-react"
import { cn } from "../../utils"

export interface FilterOption {
  value: string
  label: string
}

export interface FilterDef {
  key: string
  label: string
  /** Fixed options for enum filters. Empty/undefined = free-text filter. */
  options?: FilterOption[]
}

interface FilterBarProps {
  searchPlaceholder?: string
  searchValue: string
  onSearchChange: (value: string) => void
  filters: FilterDef[]
  activeFilters: Record<string, string>
  onFilterChange: (key: string, value: string | null) => void
}

/**
 * GitHub Issues–style search bar.
 *
 * The input is a single raw text field. Filter tokens like `state:active`
 * are parsed on every keystroke and synced to the parent via callbacks.
 * The raw text is always what the user sees — no computed display value.
 */
export function FilterBar({
  searchPlaceholder = "Search...",
  searchValue,
  onSearchChange,
  filters,
  activeFilters,
  onFilterChange,
}: FilterBarProps) {
  const inputRef = React.useRef<HTMLInputElement>(null)
  const containerRef = React.useRef<HTMLDivElement>(null)

  // The raw input is the source of truth for what the user sees.
  // We reconstruct it from parent state on mount only.
  const [rawInput, setRawInput] = React.useState(() => {
    const tokens = Object.entries(activeFilters)
      .map(([key, value]) => {
        const f = filters.find((fd) => fd.key === key)
        const o = f?.options?.find((opt) => opt.value === value)
        return `${f?.label ?? key}:${o?.label ?? value}`
      })
      .join(" ")
    return tokens ? `${tokens} ${searchValue}`.trim() : searchValue
  })

  const [suggestions, setSuggestions] = React.useState<
    { filter: FilterDef; option?: FilterOption; display: string }[]
  >([])
  const [selectedIdx, setSelectedIdx] = React.useState(0)
  const [isOpen, setIsOpen] = React.useState(false)

  // Parse raw input into filter tokens + free text, sync to parent
  const parseAndSync = React.useCallback(
    (raw: string) => {
      const words = raw.split(/\s+/).filter(Boolean)
      const newFilters: Record<string, string> = {}
      const freeWords: string[] = []

      for (const word of words) {
        const colonIdx = word.indexOf(":")
        if (colonIdx > 0 && colonIdx < word.length - 1) {
          const prefix = word.slice(0, colonIdx).toLowerCase()
          const val = word.slice(colonIdx + 1)

          const filterDef = filters.find(
            (f) =>
              f.label.toLowerCase() === prefix ||
              f.key.toLowerCase() === prefix
          )

          if (filterDef) {
            if (filterDef.options && filterDef.options.length > 0) {
              const option = filterDef.options.find(
                (o) =>
                  o.label.toLowerCase() === val.toLowerCase() ||
                  o.value.toLowerCase() === val.toLowerCase()
              )
              if (option) {
                newFilters[filterDef.key] = option.value
                continue
              }
            } else {
              // Free-text filter
              newFilters[filterDef.key] = val
              continue
            }
          }
        }
        freeWords.push(word)
      }

      // Sync filters
      const prevKeys = Object.keys(activeFilters)
      for (const key of prevKeys) {
        if (!(key in newFilters)) {
          onFilterChange(key, null)
        }
      }
      for (const [key, value] of Object.entries(newFilters)) {
        if (activeFilters[key] !== value) {
          onFilterChange(key, value)
        }
      }

      const freeText = freeWords.join(" ")
      if (freeText !== searchValue) {
        onSearchChange(freeText)
      }
    },
    [filters, activeFilters, searchValue, onFilterChange, onSearchChange]
  )

  const handleChange = (value: string) => {
    setRawInput(value)
    parseAndSync(value)
    updateSuggestions(value)
  }

  // Build autocomplete suggestions based on the last word being typed
  const updateSuggestions = (raw: string) => {
    const words = raw.split(/\s+/)
    const lastWord = (words[words.length - 1] ?? "").toLowerCase()

    if (!lastWord) {
      // Show all available filter prefixes
      const used = new Set<string>()
      // Figure out which filters are already used in the raw text
      for (const word of words) {
        const colonIdx = word.indexOf(":")
        if (colonIdx > 0) {
          const prefix = word.slice(0, colonIdx).toLowerCase()
          const f = filters.find(
            (fd) =>
              fd.label.toLowerCase() === prefix ||
              fd.key.toLowerCase() === prefix
          )
          if (f) used.add(f.key)
        }
      }
      setSuggestions(
        filters
          .filter((f) => !used.has(f.key))
          .map((f) => ({
            filter: f,
            display: `${f.label}:`,
          }))
      )
      setSelectedIdx(0)
      return
    }

    const colonIdx = lastWord.indexOf(":")
    if (colonIdx > 0) {
      // User typed prefix:partial — show matching values for enum filters
      const prefix = lastWord.slice(0, colonIdx)
      const partial = lastWord.slice(colonIdx + 1)
      const filterDef = filters.find(
        (f) =>
          f.label.toLowerCase() === prefix || f.key.toLowerCase() === prefix
      )
      if (filterDef?.options && filterDef.options.length > 0) {
        const matches = filterDef.options.filter((o) =>
          o.label.toLowerCase().startsWith(partial)
        )
        setSuggestions(
          matches.map((o) => ({
            filter: filterDef,
            option: o,
            display: `${filterDef.label}:${o.label}`,
          }))
        )
      } else {
        setSuggestions([])
      }
      setSelectedIdx(0)
      return
    }

    // Partial word — suggest matching filter names
    const matching = filters.filter(
      (f) =>
        f.label.toLowerCase().startsWith(lastWord) ||
        f.key.toLowerCase().startsWith(lastWord)
    )
    setSuggestions(
      matching.map((f) => ({
        filter: f,
        display: `${f.label}:`,
      }))
    )
    setSelectedIdx(0)
  }

  const applySuggestion = (s: {
    filter: FilterDef
    option?: FilterOption
    display: string
  }) => {
    // Replace the last word in rawInput with the suggestion
    const words = rawInput.split(/\s+/)
    words[words.length - 1] = s.display
    const newRaw = words.join(" ")

    if (s.option) {
      // Complete token — add a trailing space
      setRawInput(newRaw + " ")
      parseAndSync(newRaw)
      setSuggestions([])
      setIsOpen(false)
    } else {
      // Just prefix, keep cursor there for value typing
      setRawInput(newRaw)
      updateSuggestions(newRaw)
    }
    inputRef.current?.focus()
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (isOpen && suggestions.length > 0) {
      if (e.key === "ArrowDown") {
        e.preventDefault()
        setSelectedIdx((prev) => Math.min(prev + 1, suggestions.length - 1))
        return
      }
      if (e.key === "ArrowUp") {
        e.preventDefault()
        setSelectedIdx((prev) => Math.max(prev - 1, 0))
        return
      }
      if (e.key === "Enter" && suggestions[selectedIdx]) {
        e.preventDefault()
        applySuggestion(suggestions[selectedIdx])
        return
      }
      if (e.key === "Escape") {
        setSuggestions([])
        setIsOpen(false)
        return
      }
    }
  }

  const clearAll = () => {
    setRawInput("")
    onSearchChange("")
    for (const key of Object.keys(activeFilters)) {
      onFilterChange(key, null)
    }
    setSuggestions([])
    inputRef.current?.focus()
  }

  // Click outside closes suggestions
  React.useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(e.target as Node)
      ) {
        setIsOpen(false)
      }
    }
    document.addEventListener("mousedown", handler)
    return () => document.removeEventListener("mousedown", handler)
  }, [])

  return (
    <div ref={containerRef} className="relative max-w-xl">
      <div className="flex items-center rounded-md border bg-background px-3 focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2 focus-within:ring-offset-background">
        <Search className="h-4 w-4 text-muted-foreground shrink-0 mr-2" />
        <input
          ref={inputRef}
          type="text"
          placeholder={searchPlaceholder}
          value={rawInput}
          onChange={(e) => handleChange(e.target.value)}
          onFocus={() => {
            setIsOpen(true)
            updateSuggestions(rawInput)
          }}
          onKeyDown={handleKeyDown}
          className="flex-1 bg-transparent py-2 text-sm outline-none placeholder:text-muted-foreground"
        />
        {rawInput && (
          <button
            onClick={clearAll}
            className="shrink-0 text-muted-foreground hover:text-foreground ml-1"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        )}
      </div>

      {/* Autocomplete dropdown */}
      {isOpen && suggestions.length > 0 && (
        <div className="absolute top-full left-0 right-0 mt-1 rounded-md border bg-popover shadow-md z-50 overflow-hidden">
          <div className="max-h-56 overflow-y-auto py-1">
            {suggestions.map((s, idx) => (
              <button
                key={s.display}
                onClick={() => applySuggestion(s)}
                className={cn(
                  "flex items-center gap-1 w-full px-3 py-1.5 text-sm text-left transition-colors",
                  idx === selectedIdx
                    ? "bg-accent text-accent-foreground"
                    : "hover:bg-muted"
                )}
              >
                <span className="font-mono text-xs">{s.display}</span>
                {s.option && (
                  <span className="text-muted-foreground text-xs">
                    {s.option.label}
                  </span>
                )}
                {!s.option &&
                  s.filter.options &&
                  s.filter.options.length > 0 && (
                    <span className="text-muted-foreground text-xs">
                      (
                      {s.filter.options
                        .map((o) => o.label)
                        .join(", ")}
                      )
                    </span>
                  )}
                {!s.option &&
                  (!s.filter.options || s.filter.options.length === 0) && (
                    <span className="text-muted-foreground text-xs italic">
                      any value
                    </span>
                  )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
