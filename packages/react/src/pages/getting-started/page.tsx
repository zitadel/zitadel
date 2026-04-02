"use client"

import { ConsoleLink as Link } from "../../context/link-context"
import { 
  Info, 
  BookOpen,
  Code2,
  Terminal,
  FileCode,
  MessageSquare,
  Rocket,
  ExternalLink,
  ArrowRight,
  LogIn,
  FolderKanban,
  AppWindow
} from "lucide-react"
import { Button } from "../../components/ui/button"
import { Card, CardContent } from "../../components/ui/card"

// Tech stack icons as simple colored circles with letters
const techStackItems = [
  { name: "React", color: "bg-sky-500", letter: "R" },
  { name: "Next.js", color: "bg-foreground", letter: "N" },
  { name: "Angular", color: "bg-red-500", letter: "A" },
  { name: "Vue", color: "bg-emerald-500", letter: "V" },
  { name: "Go", color: "bg-cyan-500", letter: "Go" },
  { name: "Python", color: "bg-yellow-500", letter: "Py" },
  { name: "Java", color: "bg-orange-500", letter: "J" },
  { name: "iOS", color: "bg-zinc-800", letter: "iOS" },
]

const nextSteps = [
  {
    icon: LogIn,
    title: "Log in to your application",
    description: "Integrate your application with Zitadel for authentication and test it by logging in with your admin user.",
    action: { label: "Log in", href: "#" },
    color: "bg-amber-700",
  },
  {
    icon: FolderKanban,
    title: "Create a project",
    description: "Add a project and define its roles and role assignments.",
    action: { label: "Create project", href: "/projects" },
    color: "bg-emerald-600",
  },
  {
    icon: AppWindow,
    title: "Register your application",
    description: "Register your web, native, api or saml application and setup an authentication flow.",
    action: { label: "Register application", href: "/applications" },
    color: "bg-violet-600",
  },
]

const developerTools = [
  {
    icon: MessageSquare,
    title: "ZITADEL Community",
    description: "Discord & GitHub for help and support",
    href: "https://zitadel.com/docs/support/troubleshooting",
  },
  {
    icon: Code2,
    title: "API Reference",
    description: "Explore ZITADEL's APIs",
    href: "https://zitadel.com/docs/apis/introduction",
  },
  {
    icon: BookOpen,
    title: "Documentation",
    description: "Comprehensive guides and tutorials",
    href: "https://zitadel.com/docs",
  },
  {
    icon: FileCode,
    title: "Example Projects",
    description: "Pre-built authentication samples",
    href: "https://zitadel.com/docs/examples",
  },
  {
    icon: Terminal,
    title: "ZITADEL CLI",
    description: "Manage configuration from terminal",
    href: "https://zitadel.com/docs/guides/manage/terraform",
  },
  {
    icon: Rocket,
    title: "Quick Start Guides",
    description: "Get up and running fast",
    href: "https://zitadel.com/docs/guides/start/quickstart",
  },
]

export default function GettingStartedPage() {
  return (
    <div className="space-y-8 p-6 max-w-5xl">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-semibold">Getting Started</h1>
      </div>

      {/* Info Banner */}
      <div className="flex items-center gap-3 px-4 py-3 bg-primary/5 border border-primary/20 rounded-lg">
        <Info className="h-4 w-4 text-primary flex-shrink-0" />
        <p className="text-sm">
          New to ZITADEL? Try our onboarding guide to get started.{" "}
          <Link href="https://zitadel.com/docs/guides/start/quickstart" className="text-primary hover:underline font-medium">
            Start the guide.
          </Link>
        </p>
      </div>

      {/* Start Building Section */}
      <div className="space-y-4">
        <h2 className="text-lg font-medium">Start Building</h2>
        
        <Card className="overflow-hidden">
          <CardContent className="p-0">
            <div className="flex flex-col md:flex-row">
              {/* Left content */}
              <div className="flex-1 p-6 space-y-4">
                <h3 className="text-lg font-medium">Integrate ZITADEL into your application</h3>
                <p className="text-sm text-muted-foreground">
                  Integrate ZITADEL into your application or use one of our samples to get started in minutes.
                </p>
                <div className="flex items-center gap-3 pt-2">
                  <Button asChild size="sm">
                    <Link href="/applications">
                      Create Application
                    </Link>
                  </Button>
                  <Link 
                    href="https://zitadel.com/docs/guides/integrate" 
                    className="text-sm text-primary hover:underline"
                  >
                    Learn More
                  </Link>
                </div>
              </div>
              
              {/* Right side - Tech stack icons */}
              <div className="flex-shrink-0 p-6 bg-muted/30 flex items-center justify-center">
                <div className="grid grid-cols-4 gap-3">
                  {techStackItems.map((tech) => (
                    <div
                      key={tech.name}
                      className={`h-10 w-10 rounded-lg ${tech.color} text-white flex items-center justify-center text-xs font-semibold`}
                      title={tech.name}
                    >
                      {tech.letter}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Next Steps Section */}
      <div className="space-y-4">
        <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Your Next Steps</h2>
        
        <div className="grid sm:grid-cols-3 gap-4">
          {nextSteps.map((step) => (
            <Card key={step.title} className="hover:border-primary/50 transition-colors">
              <CardContent className="p-5 space-y-4">
                <div className="flex items-start gap-4">
                  <div className={`flex h-12 w-12 items-center justify-center rounded-full ${step.color} text-white flex-shrink-0`}>
                    <step.icon className="h-5 w-5" />
                  </div>
                  <div className="space-y-2 flex-1">
                    <h3 className="font-medium">{step.title}</h3>
                    <p className="text-sm text-muted-foreground leading-relaxed">{step.description}</p>
                  </div>
                </div>
                <Link 
                  href={step.action.href} 
                  className="text-sm text-primary hover:text-primary/80 flex items-center gap-1 font-medium"
                >
                  {step.action.label}
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>

      {/* Developer Tools Section */}
      <div className="space-y-4">
        <h2 className="text-lg font-medium">Developer Tools</h2>
        
        <div className="grid sm:grid-cols-2 gap-3">
          {developerTools.map((tool) => (
            <Link
              key={tool.title}
              href={tool.href}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 p-4 rounded-lg border hover:bg-muted/50 transition-colors group"
            >
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted text-muted-foreground flex-shrink-0">
                <tool.icon className="h-5 w-5" />
              </div>
              <div className="flex-1 min-w-0">
                <h3 className="font-medium text-sm group-hover:text-primary transition-colors">{tool.title}</h3>
                <p className="text-sm text-muted-foreground">{tool.description}</p>
              </div>
              <ExternalLink className="h-4 w-4 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0" />
            </Link>
          ))}
        </div>
      </div>

      {/* Dismiss link */}
      <div className="text-center pt-4">
        <Link href="/" className="text-sm text-muted-foreground hover:text-foreground hover:underline">
          I&apos;m done with this setup. Hide this
        </Link>
      </div>
    </div>
  )
}
