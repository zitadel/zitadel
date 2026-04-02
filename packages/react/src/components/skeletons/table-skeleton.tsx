import { Skeleton } from "../ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table"

interface TableSkeletonProps {
  /** Column headers to display */
  columns: string[]
  /** Number of skeleton rows to show */
  rows?: number
  /** Whether the first column has an avatar/icon */
  hasLeadingAvatar?: boolean
}

/**
 * Skeleton loading state for data tables.
 * Matches the table layout with shimmer placeholders.
 */
export function TableSkeleton({
  columns,
  rows = 5,
  hasLeadingAvatar = true,
}: TableSkeletonProps) {
  return (
    <div className="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            {columns.map((col) => (
              <TableHead key={col}>{col}</TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {Array.from({ length: rows }).map((_, i) => (
            <TableRow key={i} className="hover:bg-transparent">
              {columns.map((col, j) => (
                <TableCell key={col}>
                  {j === 0 && hasLeadingAvatar ? (
                    <div className="flex items-center gap-3">
                      <Skeleton className="h-8 w-8 rounded-full flex-shrink-0" />
                      <div className="space-y-1.5 flex-1">
                        <Skeleton className="h-4 w-[140px]" />
                        <Skeleton className="h-3 w-[180px]" />
                      </div>
                    </div>
                  ) : (
                    <Skeleton className="h-4 w-[80px]" />
                  )}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
