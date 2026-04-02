"use client"

import * as React from "react"
import { ConsoleLink as Link } from "../../context/link-context"
import {
  ChevronUp,
  Settings,
  Shield,
  CreditCard,
  Gauge,
  Headphones,
  MessageSquare,
  LogOut,
  Zap,
} from "lucide-react"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import { Avatar, AvatarFallback } from "../ui/avatar"
import { administrators } from "../../mock-data"

const accountMenuItems = [
  {
    title: "Upgrade to Pro",
    href: "/billing",
    icon: Zap,
    highlight: true,
  },
  {
    title: "Team Members",
    href: "/administrators",
    icon: Shield,
    badge: administrators.length.toString(),
  },
  {
    title: "Usage",
    href: "/usage",
    icon: Gauge,
  },
  {
    title: "Billing",
    href: "/billing",
    icon: CreditCard,
  },
  {
    title: "Settings",
    href: "/account-settings",
    icon: Settings,
  },
]

const supportMenuItems = [
  {
    title: "Get Support",
    href: "/support",
    icon: Headphones,
  },
  {
    title: "Give Feedback",
    href: "/feedback",
    icon: MessageSquare,
  },
]

export function AccountDropdown() {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left hover:bg-muted/50 transition-colors">
          <Avatar className="h-8 w-8">
            <AvatarFallback className="bg-muted text-foreground text-sm font-medium">
              AD
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium truncate">Admin User</p>
            <p className="text-xs text-muted-foreground truncate">admin@zitadel.cloud</p>
          </div>
          <ChevronUp className="h-4 w-4 text-muted-foreground shrink-0" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent 
        align="start" 
        side="top" 
        className="w-64"
        sideOffset={8}
      >
        {/* User Info Header */}
        <div className="px-3 py-3">
          <div className="flex items-center gap-3">
            <Avatar className="h-10 w-10">
              <AvatarFallback className="bg-muted text-foreground font-medium">
                AD
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium">Admin User</p>
              <p className="text-xs text-muted-foreground truncate">admin@zitadel.cloud</p>
            </div>
          </div>
        </div>
        
        <DropdownMenuSeparator />
        
        {/* Account Menu Items */}
        {accountMenuItems.map((item) => (
          <DropdownMenuItem key={item.title} asChild>
            <Link href={item.href} className={`flex items-center gap-3 cursor-pointer ${item.highlight ? "text-primary" : ""}`}>
              <item.icon className={`h-4 w-4 ${item.highlight ? "text-primary" : "text-muted-foreground"}`} />
              <span>{item.title}</span>
              {item.badge && (
                <span className="ml-auto text-xs text-muted-foreground">{item.badge}</span>
              )}
            </Link>
          </DropdownMenuItem>
        ))}
        
        <DropdownMenuSeparator />
        
        {/* Support Items */}
        {supportMenuItems.map((item) => (
          <DropdownMenuItem key={item.href} asChild>
            <Link href={item.href} className="flex items-center gap-3 cursor-pointer">
              <item.icon className="h-4 w-4 text-muted-foreground" />
              <span>{item.title}</span>
            </Link>
          </DropdownMenuItem>
        ))}
        
        <DropdownMenuSeparator />
        
        {/* Sign Out */}
        <DropdownMenuItem className="flex items-center gap-3 cursor-pointer text-muted-foreground">
          <LogOut className="h-4 w-4" />
          <span>Sign out</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
