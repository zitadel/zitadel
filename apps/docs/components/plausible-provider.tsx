'use client';

import Script from 'next/script';

export default function PlausibleProvider() {
    const domain = process.env.NEXT_PUBLIC_PLAUSIBLE_DOMAIN || 'zitadel.com';

    return (
        <Script
            defer
            data-domain={domain}
            data-api="/docs/pl/api/event"
            src="/docs/pl/js/script.js"
            strategy="afterInteractive"
        />
    );
}
