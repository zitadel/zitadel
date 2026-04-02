"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@zitadel/react/components/ui/card"
import { Button } from "@zitadel/react/components/ui/button"
import { Badge } from "@zitadel/react/components/ui/badge"
import { Progress } from "@zitadel/react/components/ui/progress"
import { CreditCard, Download, TrendingUp, Users, Zap, Database } from "lucide-react"

export default function BillingPage() {
  return (
    <div className="space-y-6 max-w-4xl mx-auto p-8">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Billing & Usage</h1>
        <p className="text-muted-foreground">
          Manage your subscription and monitor resource usage
        </p>
      </div>

      {/* Current Plan */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <CreditCard className="h-5 w-5" />
                Current Plan
              </CardTitle>
              <CardDescription>Your active subscription plan</CardDescription>
            </div>
            <Badge className="text-sm">Pro</Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex items-baseline gap-1">
            <span className="text-4xl font-bold">$299</span>
            <span className="text-muted-foreground">/month</span>
          </div>
          <p className="text-sm text-muted-foreground mt-2">
            Billed monthly. Next billing date: April 1, 2026
          </p>
          <div className="flex gap-3 mt-4">
            <Button>Upgrade Plan</Button>
            <Button variant="outline">Manage Subscription</Button>
          </div>
        </CardContent>
      </Card>

      {/* Usage Overview */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Monthly Active Users</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">2,450</div>
            <Progress value={49} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              49% of 5,000 limit
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">API Requests</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">1.2M</div>
            <Progress value={60} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              60% of 2M limit
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Storage Used</CardTitle>
            <Database className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">4.2 GB</div>
            <Progress value={42} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              42% of 10 GB limit
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Auth Events</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">45,231</div>
            <Progress value={23} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-2">
              23% of 200K limit
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Invoices */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Invoices</CardTitle>
          <CardDescription>Download your past invoices</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {[
              { date: "March 1, 2026", amount: "$299.00", status: "Paid" },
              { date: "February 1, 2026", amount: "$299.00", status: "Paid" },
              { date: "January 1, 2026", amount: "$299.00", status: "Paid" },
              { date: "December 1, 2025", amount: "$249.00", status: "Paid" },
            ].map((invoice, i) => (
              <div key={i} className="flex items-center justify-between py-2 border-b last:border-0">
                <div>
                  <p className="font-medium">{invoice.date}</p>
                  <p className="text-sm text-muted-foreground">{invoice.amount}</p>
                </div>
                <div className="flex items-center gap-3">
                  <Badge variant="secondary">{invoice.status}</Badge>
                  <Button variant="ghost" size="sm">
                    <Download className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
