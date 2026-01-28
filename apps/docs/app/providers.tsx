'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import dynamic from "next/dynamic";
import AuthRequestProvider from '@/utils/authrequest';
import MixpanelProvider from '@/components/mixpanel-provider';

const SearchDialog = dynamic(() => import("@/components/inkeep-search"));

export function Providers({ children }: any) {
  return (
    <RootProvider
      search={{
        enabled: true,
        SearchDialog: SearchDialog as any,
      }}
    >
      <MixpanelProvider>
        <AuthRequestProvider>{children}</AuthRequestProvider>
      </MixpanelProvider>
    </RootProvider>
  );
}
