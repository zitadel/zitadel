"use client"

import { ChevronLeft, ChevronRight } from "lucide-react"
import { Button } from "./button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./select"

interface TablePaginationProps {
  /** Current page index (0-based) */
  page: number
  /** Current page size */
  pageSize: number
  /** Total number of results from the API */
  totalResult: number
  /** Called when the page changes */
  onPageChange: (page: number) => void
  /** Called when the page size changes */
  onPageSizeChange: (pageSize: number) => void
  /** Available page size options */
  pageSizeOptions?: number[]
}

/**
 * Pagination footer for data tables.
 * Shows row range, page size selector, and prev/next buttons.
 */
export function TablePagination({
  page,
  pageSize,
  totalResult,
  onPageChange,
  onPageSizeChange,
  pageSizeOptions = [10, 20, 50],
}: TablePaginationProps) {
  const totalPages = Math.max(1, Math.ceil(totalResult / pageSize))
  const start = page * pageSize + 1
  const end = Math.min((page + 1) * pageSize, totalResult)
  const hasPrev = page > 0
  const hasNext = page < totalPages - 1

  return (
    <div className="flex items-center justify-between px-4 py-3 border-t">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <span>Rows per page</span>
        <Select
          value={String(pageSize)}
          onValueChange={(v) => {
            onPageSizeChange(Number(v))
            onPageChange(0) // Reset to first page
          }}
        >
          <SelectTrigger className="h-8 w-[70px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {pageSizeOptions.map((size) => (
              <SelectItem key={size} value={String(size)}>
                {size}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="flex items-center gap-4">
        <span className="text-sm text-muted-foreground">
          {totalResult > 0
            ? `${start}–${end} of ${totalResult}`
            : "No results"}
        </span>
        <div className="flex items-center gap-1">
          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8"
            disabled={!hasPrev}
            onClick={() => onPageChange(page - 1)}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="icon"
            className="h-8 w-8"
            disabled={!hasNext}
            onClick={() => onPageChange(page + 1)}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}
