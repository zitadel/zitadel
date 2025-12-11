'use client';

import * as React from 'react';
import { usePathname, useRouter } from 'next/navigation';
import versions from '../versions.json';

interface Version {
  version: string;
  label: string;
  type: 'remote' | 'external' | 'local';
  url?: string;
}

export function VersionSelector() {
  const pathname = usePathname();
  const router = useRouter();
  
  // Determine current version from URL
  // /docs/v2/setup -> v2
  // /docs/setup -> current
  const currentVersion = versions.find((v) => 
    pathname.startsWith(`/docs/${v.version}`)
  )?.version || 'current';

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedVersion = e.target.value;
    
    if (selectedVersion === 'current') {
      // Try to keep the slug if possible, stripping the version prefix
      // /docs/v2/setup -> /docs/setup
      const newPath = pathname.replace(/^\/docs\/v[^\/]+/, '/docs');
      router.push(newPath);
      return;
    }

    const versionData = versions.find(v => v.version === selectedVersion);
    if (!versionData) return;

    if (versionData.type === 'external' && versionData.url) {
      window.location.href = versionData.url;
      return;
    }

    // Switch to hydrated version
    // /docs/setup -> /docs/v2/setup
    // /docs/v1/setup -> /docs/v2/setup
    let slug = pathname.replace(/^\/docs/, '');
    if (currentVersion !== 'current') {
        slug = slug.replace(new RegExp(`^/${currentVersion}`), '');
    }
    router.push(`/docs/${selectedVersion}${slug}`);
  };

  return (
    <select 
      value={currentVersion} 
      onChange={handleChange}
      className="p-2 bg-transparent text-sm font-medium rounded-md border border-fd-border hover:bg-fd-accent/50 focus:outline-none focus:ring-2 focus:ring-fd-ring text-fd-foreground"
    >
      <option value="current" className="bg-fd-background text-fd-foreground">
        {process.env.NEXT_PUBLIC_VERCEL_GIT_COMMIT_REF || 'Next'} (Current)
      </option>
      {versions.map((v) => (
        <option key={v.version} value={v.version} className="bg-fd-background text-fd-foreground">
          {v.label}
        </option>
      ))}
    </select>
  );
}
