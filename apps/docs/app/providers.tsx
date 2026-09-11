'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import { ReactNode } from 'react';

import MixpanelProvider from '@/components/mixpanel-provider';
import PlausibleProvider from '@/components/plausible-provider';

export function Providers({ children }: { children: ReactNode }) {
  return (
    <RootProvider
      search={{
        options: {
          api: '/docs/api/search',
        }
      }}
    >
      <PlausibleProvider />
      <MixpanelProvider>
        {children}
      </MixpanelProvider>
    </RootProvider>
  );
}