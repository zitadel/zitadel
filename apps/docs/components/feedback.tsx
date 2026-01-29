'use client';

import { useState } from 'react';
import { mixpanelClient } from '@/utils/mixpanel';
import { usePathname } from 'next/navigation';
import { ThumbsUp, ThumbsDown } from 'lucide-react';
import { cn } from '@/utils/cn';

export function Feedback() {
  const pathname = usePathname();
  const [voted, setVoted] = useState<boolean | null>(null);

  const handleVote = (useful: boolean) => {
    setVoted(useful);
    mixpanelClient.track('helpfulness_vote', {
      path: pathname,
      useful: useful,
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
