"use client";

import {
  initMixpanel,
  trackPageView,
  hasMixpanelConsent,
  optInTracking,
  optOutTracking,
} from "@/lib/mixpanel";
import { usePathname } from "next/navigation";
import { useEffect, useRef } from "react";

export function MixpanelProvider({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const initializedRef = useRef(false);

  useEffect(() => {
    if (!initializedRef.current) {
      initMixpanel();
      initializedRef.current = true;
    }

    const handleConsentChange = () => {
      if (hasMixpanelConsent()) {
        optInTracking();
      } else {
        optOutTracking();
      }
    };

    window.addEventListener("cc:onChange:mixpanel", handleConsentChange);

    return () => {
      window.removeEventListener("cc:onChange:mixpanel", handleConsentChange);
    };
  }, []);

  useEffect(() => {
    if (initializedRef.current && pathname) {
      trackPageView(pathname);
    }
  }, [pathname]);

  return <>{children}</>;
}
