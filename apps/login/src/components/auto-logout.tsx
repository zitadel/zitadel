"use client";

import { Translated } from "@/components/translated";
import { navigateHard, shouldUseHardNavigation } from "@/lib/client-utils";
import { clearSession } from "@/lib/server/session";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

type Props = {
  sessionId: string;
  postLogoutRedirectUri?: string;
  organization?: string;
};

export function AutoLogout({ sessionId, postLogoutRedirectUri, organization }: Props) {
  const router = useRouter();

  useEffect(() => {
    let isCancelled = false;

    const performAutoLogout = async () => {
      try {
        await clearSession({ sessionId });
        if (isCancelled) {
          return;
        }

        const fallbackParams = new URLSearchParams();
        if (organization) {
          fallbackParams.set("organization", organization);
        }
        const target = postLogoutRedirectUri || `/logout/done?${fallbackParams.toString()}`;

        if (shouldUseHardNavigation(target)) {
          navigateHard(target);
          return;
        }

        router.push(target);
      } catch (error) {
        if (isCancelled) {
          return;
        }
        console.error("Auto-logout failed:", error);
        // Reload to show session selection UI on error
        router.refresh();
      }
    };

    void (async () => {
      await performAutoLogout();
    })();

    return () => {
      isCancelled = true;
    };
  }, [sessionId, postLogoutRedirectUri, organization, router]);

  return (
    <div className="flex flex-col items-center justify-center space-y-4">
      <div className="h-12 w-12 animate-spin rounded-full border-b-2 border-primary"></div>
      <p className="text-sm text-muted-foreground">
        <Translated i18nKey="autoLoggingOut" namespace="logout" />
      </p>
    </div>
  );
}
