'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import AuthRequestProvider from '@/utils/authrequest';

export function Providers({ children }: any) {
  return (
    <RootProvider>
      <AuthRequestProvider>{children}</AuthRequestProvider>
    </RootProvider>
  );
}
