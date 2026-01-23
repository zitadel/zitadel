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
  // /v2/setup -> v2
  // /setup -> current
  const currentVersion = versions.find((v) =>
    pathname.startsWith(`/${v.version}`)
  )?.version || 'current';

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedVersion = e.target.value;

    if (selectedVersion === 'current') {
      // Try to keep the slug if possible, stripping the version prefix
      // /v2/setup -> /setup
      const newPath = pathname.replace(/^\/v[^\/]+/, '');
      // Ensure we don't end up with empty string, default to /
      router.push(newPath || '/');
      return;
    }

    const versionData = versions.find(v => v.version === selectedVersion);
    if (!versionData) return;

    if (versionData.type === 'external' && versionData.url) {
      window.location.href = versionData.url;
      return;
    }

    // Switch to hydrated version
    // /setup -> /v2/setup
    // /v1/setup -> /v2/setup
    let slug = pathname;
    if (currentVersion !== 'current') {
      slug = slug.replace(new RegExp(`^/${currentVersion}`), '');
    }
    // ensure slug starts with /
    if (!slug.startsWith('/')) slug = '/' + slug;

    router.push(`/${selectedVersion}${slug}`);
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
