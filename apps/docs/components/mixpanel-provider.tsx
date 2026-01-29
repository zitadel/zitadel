'use client';

import { initMixpanel, mixpanelClient, optInTracking, optOutTracking } from '@/utils/mixpanel';
import { usePathname } from 'next/navigation';
import React, { useEffect } from 'react';

export default function MixpanelProvider({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  useEffect(() => {
    const setup = async () => {
      await initMixpanel();

      const CookieConsent = await import('vanilla-cookieconsent');
      const cookie = CookieConsent.getCookie();

      // @ts-ignore
      if (!cookie || !cookie.services || !cookie.services.analytics) {
        // We only send data when people opt-in, not just if they have not opted out.
        return;
      }
      // @ts-ignore
      const hasMixpanelConsent = cookie?.services?.analytics?.indexOf('mixpanel') !== -1;
      if (hasMixpanelConsent) {
        optInTracking();
      } else {
        optOutTracking();
      }
    };
    setup().catch((error) => {
      console.error('Error setting up Mixpanel:', error);
    });

    const handleConsentChange = async () => {
      const CookieConsent = await import('vanilla-cookieconsent');
      const cookie = CookieConsent.getCookie();
      // @ts-ignore
      const hasMixpanelConsent = cookie?.services?.analytics?.indexOf('mixpanel') !== -1;
      if (hasMixpanelConsent) {
        optInTracking();
      } else {
        optOutTracking();
      }
    };

    window.addEventListener('cc:onChange:mixpanel', handleConsentChange);

    return () => {
      window.removeEventListener('cc:onChange:mixpanel', handleConsentChange);
    };
  }, []);

  useEffect(() => {
    if (pathname) {
      mixpanelClient.track('page_viewed', { path: pathname });
    }
  }, [pathname]);

  return <>{children}</>;
}
