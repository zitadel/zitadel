'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import { ReactNode } from 'react';
import { useParams } from 'next/navigation';

import MixpanelProvider from '@/components/mixpanel-provider';
import PlausibleProvider from '@/components/plausible-provider';
import { getVersionFromSlug } from '@/lib/versions';

export function Providers({ children }: { children: ReactNode }) {
  const params = useParams();
  const slug = params?.slug as string[] | undefined;

  return (
    <RootProvider
      search={{
        options: {
          api: '/docs/api/search',
          // Every search index entry is tagged with its docs version (see
          // app/api/search/route.ts). `defaultTag` applies that filter; `tags`
          // would only add a visible filter pill.
          defaultTag: getVersionFromSlug(slug),
        },
      }}
    >
      <PlausibleProvider />
      <MixpanelProvider>
        {children}
      </MixpanelProvider>
    </RootProvider>
  );
}
