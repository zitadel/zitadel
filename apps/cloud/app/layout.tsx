import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import './globals.css'
import { DebugBanner } from '@/components/debug-banner'

const _geist = Geist({ subsets: ["latin"] })
const _geistMono = Geist_Mono({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: 'ZITADEL Cloud',
  description: 'ZITADEL Cloud — manage instances, billing, and more',
}

/**
 * Root layout — bare shell. No sidebar, no providers.
 * 
 * The console sidebar and providers are scoped to /console/* via console/layout.tsx.
 * Docs and debug pages get their own layouts without the console chrome.
 */
export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="font-sans antialiased" suppressHydrationWarning>
        <DebugBanner />
        {children}
      </body>
    </html>
  )
}
