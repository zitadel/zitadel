'use client';

import * as React from 'react';
import { usePathname, useRouter } from 'next/navigation';
import versions from '../content/versions.json';

interface Version {
  param?: string;
  label: string;
  url: string;
  ref?: string;
  refType?: string;
  type?: string;
}

export function VersionSelector() {
  const pathname = usePathname();
  const router = useRouter();

  // Determine current version from URL
  // /docs -> latest
  // /docs/v4.10 -> v4.10
  // Default to latest if no match found

  // Create a map for easy lookup
  // We want to match longest prefix first, so "v4.10" matches before "latest" (which is root /docs)
  // Actually, simplified:
  // If path starts with /docs/vX.Y, then it is vX.Y
  // If path is /docs or /docs/..., and not /docs/vX.Y, it is likely latest (assuming /docs prefix)

  // Check versions with 'param' defined.
  const currentVersionObj = versions.find((v) => {
    if (!v.param) return false;
    if (v.param === 'latest') {
      // 'latest' usually maps to root /docs (without version prefix)
      // So if it DOES NOT start with any OTHER version prefix, it is likely latest.
      // However, let's look for explicit matches first.
      return false;
    }
    return pathname.startsWith(`/docs/${v.param}`);
  });

  const currentVersion = currentVersionObj ? currentVersionObj.param : 'latest';

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedParam = e.target.value;
    const selectedVersion = versions.find(v => v.param === selectedParam || (v.type === 'external' && v.label === selectedParam)); // Fallback for external

    if (!selectedVersion) {
      // Handle external link selected via label if param is missing?
      // Our versions.json has explicit external entries without param sometimes?
      // Actually the Archive entry has no param.
      // Let's iterate versions to find match by URL or Label if param missing
      const external = versions.find(v => v.label === selectedParam); // value in option might be label if param missing
      if (external && external.type === 'external') {
        window.location.href = external.url;
      }
      return;
    }

    if (selectedVersion.type === 'external') {
      window.location.href = selectedVersion.url;
      return;
    }

    // Switching logic
    if (selectedParam === 'latest') {
      // Switch to /docs/...
      // If current is /docs/v4.10/foo -> /docs/foo
      const newPath = pathname.replace(new RegExp(`^/docs/${currentVersion}`), '/docs');
      router.push(newPath);
    } else {
      // Switch to /docs/v4.x/...
      if (currentVersion === 'latest') {
        // /docs/foo -> /docs/v4.x/foo
        const newPath = pathname.replace(/^\/docs/, `/docs/${selectedParam}`);
        router.push(newPath);
      } else {
        // /docs/vOld/foo -> /docs/vNew/foo
        const newPath = pathname.replace(new RegExp(`^/docs/${currentVersion}`), `/docs/${selectedParam}`);
        router.push(newPath);
      }
    }
  };

  return (
    <select
      value={currentVersion}
      onChange={handleChange}
      className="p-2 bg-transparent text-sm font-medium rounded-md border border-fd-border hover:bg-fd-accent/50 focus:outline-none focus:ring-2 focus:ring-fd-ring text-fd-foreground"
    >
      {versions.map((v) => (
        <option
          key={v.param || v.label}
          value={v.param || v.label}
          className="bg-fd-background text-fd-foreground"
        >
          {v.label}
        </option>
      ))}
    </select>
  );
}
