"use client"

import { InstanceSwitcher } from "./instance-switcher"
import { OrganizationSwitcher } from "./organization-switcher"
import { DocumentationNav } from "./documentation-nav"
import { GlobalSearch } from "./global-search"
import { Separator } from "../ui/separator"
import { useDeployment } from "../../context/deployment"

export function Header() {
  const { isCloud } = useDeployment()

  return (
    <header className="sticky top-0 z-50 flex h-14 items-center gap-4 border-b bg-background px-4">
      {/* Left section - Instance Switcher (cloud only) */}
      {isCloud && (
        <>
          <div className="flex items-center gap-2">
            <InstanceSwitcher />
          </div>
          <Separator orientation="vertical" className="h-6" />
        </>
      )}
      
      {/* Organization Switcher */}
      <div className="flex items-center gap-2">
        <OrganizationSwitcher />
      </div>
      
      {/* Center section - Global Search */}
      <GlobalSearch />
      
      {/* Right section - Documentation */}
      <div className="flex items-center gap-2">
        <DocumentationNav />
      </div>
    </header>
  )
}

