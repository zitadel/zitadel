import './global.css';
import type { Metadata } from 'next';
import { Providers } from './providers';

export const metadata: Metadata = {
  metadataBase: new URL(process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3000'),
  icons: {
    other: [
      {
        rel: 'stylesheet',
        url: '/img/icons/line-awesome/css/line-awesome.min.css',
      },
    ],
  },
};

export default function Layout({ children }: any) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="flex flex-col min-h-screen font-sans bg-fd-background text-fd-foreground">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
