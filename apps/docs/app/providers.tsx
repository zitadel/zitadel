'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import { ReactNode } from 'react';
import { useParams } from 'next/navigation';

import MixpanelProvider from '@/components/mixpanel-provider';
import PlausibleProvider from '@/components/plausible-provider';

export function Providers({ children }: { children: ReactNode }) {
  const params = useParams();
  
  // Extract version from Next.js dynamic route params, fallback to 'latest'
  const slug = params?.slug as string[] | undefined;
  const currentVersion = slug?.[0]?.startsWith('v') ? slug[0] : 'latest';

  return (
    <RootProvider
      search={{
        options: {
          api: '/docs/api/search',
          // Scope the search results silently to the current version
          tags: [
            {
              name: currentVersion === 'latest' ? 'Latest' : currentVersion,
              value: currentVersion,
            }
          ]
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