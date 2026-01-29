'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import dynamic from "next/dynamic";
import { ThemeProvider } from 'next-themes';

import MixpanelProvider from '@/components/mixpanel-provider';
import PlausibleProvider from '@/components/plausible-provider';

const SearchDialog = dynamic(() => import("@/components/inkeep-search"));

export function Providers({ children }: { children: React.ReactNode }) {
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
