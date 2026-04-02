"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card"
import { analyticsData } from "../../mock-data"
import { BarChart3, Users, Activity, TrendingUp, TrendingDown } from "lucide-react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../components/ui/tabs"
import { useAppContext } from "../../context/app-context"
import { InstanceSelectorPrompt } from "../../components/instance-selector-prompt"

function StatCard({ 
  title, 
  value, 
  change, 
  changeType,
  icon: Icon,
}: { 
  title: string
  value: string | number
  change: string
  changeType: "up" | "down" | "neutral"
  icon: React.ComponentType<{ className?: string }>
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        <div className={`flex items-center text-xs ${
          changeType === "up" ? "text-foreground" : 
          changeType === "down" ? "text-muted-foreground" : 
          "text-muted-foreground"
        }`}>
          {changeType === "up" && <TrendingUp className="mr-1 h-3 w-3" />}
          {changeType === "down" && <TrendingDown className="mr-1 h-3 w-3" />}
          {change}
        </div>
      </CardContent>
    </Card>
  )
}

function SimpleBarChart({ data, dataKey, color }: { data: typeof analyticsData, dataKey: keyof typeof analyticsData[0], color: string }) {
  const maxValue = Math.max(...data.map(d => Number(d[dataKey])))
  
  return (
    <div className="flex items-end gap-1 h-48">
      {data.map((day, i) => (
        <div key={day.date} className="flex-1 flex flex-col items-center gap-1">
          <div 
            className={`w-full rounded-t ${color}`}
            style={{ height: `${(Number(day[dataKey]) / maxValue) * 100}%`, minHeight: '4px' }}
          />
          {i % 5 === 0 && (
            <span className="text-xs text-muted-foreground">
              {new Date(day.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
            </span>
          )}
        </div>
      ))}
    </div>
  )
}

export default function AnalyticsPage() {
  const { currentInstance } = useAppContext()

  if (!currentInstance) {
    return (
      <InstanceSelectorPrompt 
        title="Continue to Analytics"
        description="Choose an instance to view analytics"
        icon={<BarChart3 className="h-6 w-6 text-muted-foreground" />}
        targetPath="/analytics"
      />
    )
  }

  const last30Days = analyticsData
  const last7Days = analyticsData.slice(-7)
  
  const totalRequests = last30Days.reduce((sum, d) => sum + d.apiRequests, 0)
  const totalActiveUsers = last30Days.reduce((sum, d) => sum + d.activeUsers, 0) / 30
  const totalNewUsers = last30Days.reduce((sum, d) => sum + d.newUsers, 0)
  const totalSessions = last30Days.reduce((sum, d) => sum + d.sessions, 0)
  
  // Calculate changes
  const prev7DaysRequests = last30Days.slice(0, 7).reduce((sum, d) => sum + d.apiRequests, 0)
  const curr7DaysRequests = last7Days.reduce((sum, d) => sum + d.apiRequests, 0)
  const requestsChange = ((curr7DaysRequests - prev7DaysRequests) / prev7DaysRequests * 100).toFixed(1)

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Analytics</h1>
        <p className="text-muted-foreground">
          API usage and user activity metrics
        </p>
      </div>

      {/* Stats Overview */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="API Requests"
          value={totalRequests.toLocaleString()}
          change={`${requestsChange}% from last period`}
          changeType={Number(requestsChange) >= 0 ? "up" : "down"}
          icon={BarChart3}
        />
        <StatCard
          title="Avg. Active Users"
          value={Math.round(totalActiveUsers).toLocaleString()}
          change="Daily average"
          changeType="neutral"
          icon={Users}
        />
        <StatCard
          title="New Users"
          value={totalNewUsers.toLocaleString()}
          change="Last 30 days"
          changeType="up"
          icon={Users}
        />
        <StatCard
          title="Sessions"
          value={totalSessions.toLocaleString()}
          change="Last 30 days"
          changeType="neutral"
          icon={Activity}
        />
      </div>

      {/* Charts */}
      <Tabs defaultValue="requests" className="space-y-4">
        <TabsList>
          <TabsTrigger value="requests">API Requests</TabsTrigger>
          <TabsTrigger value="users">Active Users</TabsTrigger>
          <TabsTrigger value="sessions">Sessions</TabsTrigger>
          <TabsTrigger value="new-users">New Users</TabsTrigger>
        </TabsList>

        <TabsContent value="requests">
          <Card>
            <CardHeader>
              <CardTitle>API Requests</CardTitle>
              <CardDescription>Total API requests over the last 30 days</CardDescription>
            </CardHeader>
            <CardContent>
              <SimpleBarChart data={last30Days} dataKey="apiRequests" color="bg-primary" />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="users">
          <Card>
            <CardHeader>
              <CardTitle>Active Users</CardTitle>
              <CardDescription>Daily active users over the last 30 days</CardDescription>
            </CardHeader>
            <CardContent>
              <SimpleBarChart data={last30Days} dataKey="activeUsers" color="bg-foreground" />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="sessions">
          <Card>
            <CardHeader>
              <CardTitle>Sessions</CardTitle>
              <CardDescription>User sessions over the last 30 days</CardDescription>
            </CardHeader>
            <CardContent>
              <SimpleBarChart data={last30Days} dataKey="sessions" color="bg-foreground/80" />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="new-users">
          <Card>
            <CardHeader>
              <CardTitle>New Users</CardTitle>
              <CardDescription>New user registrations over the last 30 days</CardDescription>
            </CardHeader>
            <CardContent>
              <SimpleBarChart data={last30Days} dataKey="newUsers" color="bg-foreground/60" />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Data Table */}
      <Card>
        <CardHeader>
          <CardTitle>Daily Breakdown</CardTitle>
          <CardDescription>Detailed metrics for the last 7 days</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2 font-medium">Date</th>
                  <th className="text-right py-2 font-medium">API Requests</th>
                  <th className="text-right py-2 font-medium">Active Users</th>
                  <th className="text-right py-2 font-medium">New Users</th>
                  <th className="text-right py-2 font-medium">Sessions</th>
                </tr>
              </thead>
              <tbody>
                {last7Days.reverse().map((day) => (
                  <tr key={day.date} className="border-b last:border-0">
                    <td className="py-2">{new Date(day.date).toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })}</td>
                    <td className="text-right py-2">{day.apiRequests.toLocaleString()}</td>
                    <td className="text-right py-2">{day.activeUsers.toLocaleString()}</td>
                    <td className="text-right py-2">{day.newUsers.toLocaleString()}</td>
                    <td className="text-right py-2">{day.sessions.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
