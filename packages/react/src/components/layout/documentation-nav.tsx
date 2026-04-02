"use client"

import { Button } from "../ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import { ChevronDown, Book, FileText, Code, ExternalLink, Video, Lightbulb } from "lucide-react"

export function DocumentationNav() {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="gap-2">
          <Book className="h-4 w-4" />
          <span>Documentation</span>
          <ChevronDown className="h-4 w-4 opacity-50" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56" align="end">
        <DropdownMenuLabel>Resources</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem asChild>
          <a href="https://zitadel.com/docs" target="_blank" rel="noopener noreferrer" className="flex items-center">
            <FileText className="mr-2 h-4 w-4" />
            <span className="flex-1">Documentation</span>
            <ExternalLink className="h-3 w-3 opacity-50" />
          </a>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <a href="https://zitadel.com/docs/apis/introduction" target="_blank" rel="noopener noreferrer" className="flex items-center">
            <Code className="mr-2 h-4 w-4" />
            <span className="flex-1">API Reference</span>
            <ExternalLink className="h-3 w-3 opacity-50" />
          </a>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <a href="https://zitadel.com/docs/guides" target="_blank" rel="noopener noreferrer" className="flex items-center">
            <Lightbulb className="mr-2 h-4 w-4" />
            <span className="flex-1">Guides & Tutorials</span>
            <ExternalLink className="h-3 w-3 opacity-50" />
          </a>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <a href="https://www.youtube.com/@zaborostraea7674" target="_blank" rel="noopener noreferrer" className="flex items-center">
            <Video className="mr-2 h-4 w-4" />
            <span className="flex-1">Video Tutorials</span>
            <ExternalLink className="h-3 w-3 opacity-50" />
          </a>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuLabel>Quick Links</DropdownMenuLabel>
        <DropdownMenuItem asChild>
          <a href="https://github.com/zitadel/zitadel" target="_blank" rel="noopener noreferrer" className="flex items-center">
            <Code className="mr-2 h-4 w-4" />
            <span className="flex-1">GitHub Repository</span>
            <ExternalLink className="h-3 w-3 opacity-50" />
          </a>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
