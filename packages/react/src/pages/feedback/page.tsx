"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../../components/ui/card"
import { Button } from "../../components/ui/button"
import { Input } from "../../components/ui/input"
import { Textarea } from "../../components/ui/textarea"
import { useAppContext } from "../../context/app-context"
import { MessageSquare, ThumbsUp, ThumbsDown, Star } from "lucide-react"
import { useState } from "react"

export default function FeedbackPage() {
  const { currentInstance } = useAppContext()
  const [rating, setRating] = useState<number | null>(null)
  const [sentiment, setSentiment] = useState<"positive" | "negative" | null>(null)

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
        <h1 className="text-3xl font-bold tracking-tight">Give Feedback</h1>
        <p className="text-muted-foreground">
          Help us improve ZITADEL with your feedback
        </p>
      </div>

      {/* Quick Feedback */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <MessageSquare className="h-5 w-5" />
            Quick Feedback
          </CardTitle>
          <CardDescription>
            How is your experience with ZITADEL?
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Sentiment */}
          <div className="space-y-3">
            <label className="text-sm font-medium">Overall Experience</label>
            <div className="flex gap-4">
              <Button
                variant={sentiment === "positive" ? "default" : "outline"}
                size="lg"
                onClick={() => setSentiment("positive")}
                className="flex-1"
              >
                <ThumbsUp className="mr-2 h-5 w-5" />
                Positive
              </Button>
              <Button
                variant={sentiment === "negative" ? "destructive" : "outline"}
                size="lg"
                onClick={() => setSentiment("negative")}
                className="flex-1"
              >
                <ThumbsDown className="mr-2 h-5 w-5" />
                Needs Improvement
              </Button>
            </div>
          </div>

          {/* Star Rating */}
          <div className="space-y-3">
            <label className="text-sm font-medium">Rate your experience (1-5)</label>
            <div className="flex gap-2">
              {[1, 2, 3, 4, 5].map((star) => (
                <Button
                  key={star}
                  variant="ghost"
                  size="lg"
                  onClick={() => setRating(star)}
                  className="p-2"
                >
                  <Star
                    className={`h-8 w-8 ${
                      rating && star <= rating
                        ? "fill-foreground text-foreground"
                        : "text-muted-foreground"
                    }`}
                  />
                </Button>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Detailed Feedback */}
      <Card>
        <CardHeader>
          <CardTitle>Detailed Feedback</CardTitle>
          <CardDescription>
            Share specific thoughts or suggestions
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Feedback Type</label>
              <select className="w-full h-10 rounded-md border border-input bg-background px-3 text-sm">
                <option>General Feedback</option>
                <option>Feature Suggestion</option>
                <option>UI/UX Improvement</option>
                <option>Documentation</option>
                <option>Performance</option>
              </select>
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Title</label>
              <Input placeholder="Brief summary of your feedback" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Details</label>
              <Textarea 
                placeholder="Tell us more about your feedback..."
                className="min-h-[150px]"
              />
            </div>
            <div className="flex gap-3">
              <Button type="submit">Submit Feedback</Button>
              <Button type="button" variant="outline">Cancel</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
