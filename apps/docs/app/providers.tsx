'use client';

import { RootProvider } from 'fumadocs-ui/provider/next';
import { ReactNode } from 'react';
import MixpanelProvider from '@/components/mixpanel-provider';
import PlausibleProvider from '@/components/plausible-provider';

// Import our new custom dialog
import CustomSearchDialog from '@/components/custom-search-dialog';

export function Providers({ children }: { children: ReactNode }) {
  return (
    <RootProvider
      search={{
        // Replace Fumadocs' default Orama Dialog with our Custom Debounced Dialog
        SearchDialog: CustomSearchDialog,
      }}
    >
      <PlausibleProvider />
      <MixpanelProvider>
        {children}
      </MixpanelProvider>
    </RootProvider>
  );
}