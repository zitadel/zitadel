"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card"
import { Button } from "../../components/ui/button"
import { Input } from "../../components/ui/input"
import { Textarea } from "../../components/ui/textarea"
import { useAppContext } from "../../context/app-context"
import { HelpCircle, BookOpen, MessageCircle, Mail, ExternalLink } from "lucide-react"
import { ConsoleLink as Link } from "../../context/link-context"

export default function SupportPage() {
  const { currentInstance } = useAppContext()

  if (!currentInstance) {
    return (
      <div className="flex items-center justify-center h-[50vh]">
        <p className="text-muted-foreground">Please select an instance first.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Get Support</h1>
        <p className="text-muted-foreground">
          Find help and resources for ZITADEL
        </p>
      </div>

      {/* Quick Links */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card className="hover:shadow-md transition-shadow cursor-pointer">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <BookOpen className="h-5 w-5" />
              Documentation
            </CardTitle>
            <CardDescription>
              Browse our comprehensive documentation
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="https://zitadel.com/docs" target="_blank">
              <Button variant="outline" className="w-full">
                View Docs
                <ExternalLink className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card className="hover:shadow-md transition-shadow cursor-pointer">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <MessageCircle className="h-5 w-5" />
              Community
            </CardTitle>
            <CardDescription>
              Join our community on Discord
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="https://discord.gg/zitadel" target="_blank">
              <Button variant="outline" className="w-full">
                Join Discord
                <ExternalLink className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card className="hover:shadow-md transition-shadow cursor-pointer">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <HelpCircle className="h-5 w-5" />
              FAQs
            </CardTitle>
            <CardDescription>
              Frequently asked questions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link href="https://zitadel.com/docs/help" target="_blank">
              <Button variant="outline" className="w-full">
                View FAQs
                <ExternalLink className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>

      {/* Contact Support */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Mail className="h-5 w-5" />
            Contact Support
          </CardTitle>
          <CardDescription>
            Submit a support request and we will get back to you
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <label className="text-sm font-medium">Subject</label>
                <Input placeholder="Brief description of your issue" />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Category</label>
                <select className="w-full h-10 rounded-md border border-input bg-background px-3 text-sm">
                  <option>Technical Issue</option>
                  <option>Billing Question</option>
                  <option>Feature Request</option>
                  <option>Account Access</option>
                  <option>Other</option>
                </select>
              </div>
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Description</label>
              <Textarea 
                placeholder="Please describe your issue in detail..."
                className="min-h-[150px]"
              />
            </div>
            <Button type="submit">Submit Request</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
