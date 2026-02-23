'use client';

import { ReactNode } from 'react';
import { RootProvider } from 'fumadocs-ui/provider/next';
import dynamic from "next/dynamic";

import MixpanelProvider from '@/components/mixpanel-provider';
import PlausibleProvider from '@/components/plausible-provider';

const SearchDialog = dynamic(() => import("@/components/inkeep-search"));

export function Providers({ children }: { children: ReactNode }) {
  return (
    <RootProvider
      search={{
        enabled: true,
        SearchDialog: SearchDialog as any,
      }}
    >
      <PlausibleProvider />
      <MixpanelProvider>
        {children}
      </MixpanelProvider>
    </RootProvider>
  );
}
