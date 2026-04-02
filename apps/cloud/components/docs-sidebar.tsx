import Link from "next/link"
import { listDocSections, type SidebarNode } from "@/lib/docs-content"

/**
 * Docs sidebar — shows the documentation navigation tree.
 * Fetched at request time with ISR.
 */
export async function DocsSidebar({ currentSlug }: { currentSlug?: string[] }) {
  const sections = await listDocSections()
  const currentPath = currentSlug?.join("/") ?? ""

  return (
    <aside className="w-64 border-r shrink-0 overflow-y-auto p-4 hidden md:block">
      <div className="mb-4">
        <Link href="/docs" className="font-semibold text-sm hover:text-foreground">
          ZITADEL Docs
        </Link>
      </div>
      <nav className="space-y-1">
        {sections.map((section) => (
          <SidebarItem
            key={section.href}
            node={section}
            currentPath={currentPath}
          />
        ))}
      </nav>
      <div className="mt-6 pt-4 border-t">
        <Link href="/" className="text-xs text-muted-foreground hover:text-foreground">
          ← Cloud Home
        </Link>
      </div>
    </aside>
  )
}

function SidebarItem({ node, currentPath }: { node: SidebarNode; currentPath: string }) {
  const isActive = currentPath.startsWith(node.slug.join("/"))

  return (
    <div>
      <Link
        href={node.href}
        className={`block px-2 py-1.5 rounded-md text-sm transition-colors ${
          isActive
            ? "bg-accent text-foreground font-medium"
            : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
        }`}
      >
        {node.title}
      </Link>
      {node.children && isActive && (
        <div className="ml-3 mt-1 space-y-0.5 border-l pl-2">
          {node.children.map((child) => (
            <SidebarItem key={child.href} node={child} currentPath={currentPath} />
          ))}
        </div>
      )}
    </div>
  )
}
