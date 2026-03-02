'use client';

import { useState, FormEvent } from 'react';
import { mixpanelClient } from '@/utils/mixpanel';
import { usePathname } from 'next/navigation';
import { ThumbsUp, ThumbsDown } from 'lucide-react';
import { cn } from '@/utils/cn';

export function Feedback() {
  const MAX_CHARS = 500;
  const pathname = usePathname();
  const [voted, setVoted] = useState<boolean | null>(null);
  const [showReasonInput, setShowReasonInput] = useState(false);
  const [reason, setReason] = useState('');

  const handleVote = (useful: boolean) => {
    if (!useful) {
      setShowReasonInput(true);
      return;
    }

    setVoted(true);
    mixpanelClient.track('helpfulness_vote', {
      path: pathname,
      useful: true,
    });
  };

  const handleSubmitNegativeFeedback = (e: FormEvent) => {
    e.preventDefault();

    if (!reason.trim()) return;

    setVoted(false);
    setShowReasonInput(false);

    mixpanelClient.track('helpfulness_vote', {
      path: pathname,
      useful: false,
      reason: reason.trim(),
    });
  };

  if (voted !== null) {
    return (
      <div className="flex flex-col gap-2 border-t pt-8 mt-8">
        <p className="text-sm text-fd-muted-foreground">
          Thanks for your feedback!
        </p>
      </div>
    );
  }

  if (showReasonInput) {
    return (
      <div className="flex flex-col gap-3 max-w-md pt-8 mt-8 border-t">
        <label
          htmlFor="feedback-reason"
          className="text-sm font-medium text-fd-foreground"
        >
          How can we improve this page?
        </label>
        <form onSubmit={handleSubmitNegativeFeedback} className="flex flex-col gap-2">
          <textarea
            id="feedback-reason"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="What was missing or confusing?"
            required
            maxLength={MAX_CHARS}
            className="w-full p-2 text-sm bg-transparent border rounded-md resize-y border-fd-border text-fd-foreground focus:outline-none focus:ring-2 focus:ring-fd-accent min-h-[80px]"
          />
          <div className="flex items-center justify-between mt-1">
            <span className={cn(
              "text-xs",
              reason.length >= MAX_CHARS ? "text-red-500 font-medium" : "text-fd-muted-foreground"
            )}>
              {reason.length} / {MAX_CHARS}
            </span>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => {
                  setShowReasonInput(false);
                  setReason('');
                }}
                className="px-4 py-2 text-sm font-medium transition-colors border rounded-md border-fd-border hover:bg-fd-accent/10"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={!reason.trim()}
                className="px-4 py-2 text-sm font-medium transition-colors rounded-md bg-fd-accent text-fd-accent-foreground disabled:opacity-50 hover:opacity-90"
              >
                Submit
              </button>
            </div>
          </div>
        </form>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-2 border-t pt-8 mt-8">
      <p className="text-sm text-fd-muted-foreground">Was this page helpful?</p>
      <div className="flex gap-2">
        <button
          onClick={() => handleVote(true)}
          className={cn(
            "p-2 rounded-md transition-colors border hover:bg-fd-accent hover:text-fd-accent-foreground border-fd-border text-fd-muted-foreground"
          )}
          aria-label="Helpful"
        >
          <ThumbsUp className="size-4" />
        </button>
        <button
          onClick={() => handleVote(false)}
          className={cn(
            "p-2 rounded-md transition-colors border hover:bg-fd-accent hover:text-fd-accent-foreground border-fd-border text-fd-muted-foreground"
          )}
          aria-label="Not helpful"
        >
          <ThumbsDown className="size-4" />
        </button>
      </div>
    </div>
  );
}