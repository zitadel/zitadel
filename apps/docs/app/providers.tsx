'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import dynamic from "next/dynamic";
import AuthRequestProvider from '@/utils/authrequest';
import MixpanelProvider from '@/components/mixpanel-provider';
import PlausibleProvider from '@/components/plausible-provider';

const SearchDialog = dynamic(() => import("@/components/inkeep-search"));

export function Providers({ children }: any) {
  return (
    <RootProvider
      search={{
        enabled: true,
        SearchDialog: SearchDialog as any,
      }}
    >
      <PlausibleProvider />
      <MixpanelProvider>
        <AuthRequestProvider>{children}</AuthRequestProvider>
      </MixpanelProvider>
    </RootProvider>
  );
}
