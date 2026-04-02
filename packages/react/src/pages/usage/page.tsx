"use client"

import { Gauge, Users, Server, Activity, TrendingUp } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card"
import { Progress } from "../../components/ui/progress"
import { useAppContext } from "../../context/app-context"
import { InstanceSelectorPrompt } from "../../components/instance-selector-prompt"

const usageData = {
  activeUsers: {
    current: 487,
    limit: 1000,
    label: "Active Users",
  },
  apiCalls: {
    current: 245000,
    limit: 500000,
    label: "API Calls (Monthly)",
  },
  storage: {
    current: 2.4,
    limit: 10,
    label: "Storage (GB)",
  },
  sessions: {
    current: 1250,
    limit: 5000,
    label: "Active Sessions",
  },
}

export default function UsagePage() {
  const { currentInstance } = useAppContext()

  if (!currentInstance) {
    return (
      <InstanceSelectorPrompt 
        title="Continue to Usage"
        description="Choose an instance to view usage"
        icon={<Gauge className="h-6 w-6 text-muted-foreground" />}
        targetPath="/usage"
      />
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Usage</h1>
        <p className="text-muted-foreground">
          Monitor your resource usage and quotas
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {Object.entries(usageData).map(([key, data]) => {
          const percentage = (data.current / data.limit) * 100
          return (
            <Card key={key}>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">{data.label}</CardTitle>
                <CardDescription>
                  {data.current.toLocaleString()} / {data.limit.toLocaleString()}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Progress value={percentage} className="h-2" />
                <p className="mt-2 text-xs text-muted-foreground">
                  {percentage.toFixed(1)}% used
                </p>
              </CardContent>
            </Card>
          )
        })}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Usage Trends</CardTitle>
          <CardDescription>
            Your resource consumption over the last 30 days
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-48 text-muted-foreground">
            <div className="text-center">
              <TrendingUp className="h-12 w-12 mx-auto mb-2 opacity-50" />
              <p>Usage analytics chart would appear here</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
