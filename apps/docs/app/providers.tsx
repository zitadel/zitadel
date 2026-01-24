'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import AuthRequestProvider from '@/utils/authrequest';
import MixpanelProvider from '@/components/mixpanel-provider';

export function Providers({ children }: any) {
  return (
    <RootProvider>
      <MixpanelProvider>
        <AuthRequestProvider>{children}</AuthRequestProvider>
      </MixpanelProvider>
    </RootProvider>
  );
}
