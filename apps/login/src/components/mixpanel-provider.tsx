"use client";

import { initMixpanel, trackPageView } from "@/lib/mixpanel";
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
  }, []);

  useEffect(() => {
    if (initializedRef.current && pathname) {
      trackPageView(pathname);
    }
  }, [pathname]);

  return <>{children}</>;
}
