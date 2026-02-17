'use client';

import * as React from 'react';
import { usePathname, useRouter } from 'next/navigation';
interface DocVersion {
    param: string;
    label: string;
    url: string;
    ref: string;
    refType: string;
    type?: 'external';
}

import rawVersions from '../content/versions.json';
const versions = rawVersions as DocVersion[];

export function VersionSelector() {
    const pathname = usePathname();
    const router = useRouter();

    // Handle basePath issues: usePathname() typically returns path relative to basePath if configured?
    // But purely relying on that is tricky. Let's normalize.
    // If pathname starts with /docs, assume it's the full path, and we strip it for logic,
    // then let router handle re-adding it if needed (or we check behaviors).
    // Actually, Next.js router.push() expects path relative to basePath usually?
    // Let's implement robust segment swapping.

    const normalizePath = (p: string) => {
        if (p.startsWith('/docs')) return p.slice(5) || '/';
        return p;
    };

    const normalizedPath = normalizePath(pathname);

    const currentVersionObj = versions.find((v) => {
        if (!v.param || v.param === 'latest') return false;
        return normalizedPath.startsWith(`/${v.param}`);
    });

    const currentVersion = currentVersionObj ? currentVersionObj.param : 'latest';

    const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const value = e.target.value;
        const version = versions.find((v) => (v.param || v.label) === value);

        if (!version) return;

        if (!version.param) return;

        // Calculate the 'tail' of the path (content path)
        let tail = normalizedPath;
        if (currentVersion !== 'latest') {
            // Strip current version prefix
            // e.g. /v4.10/foo -> /foo
            tail = normalizedPath.replace(new RegExp(`^/${currentVersion}`), '') || '/';
        }

        // Construct new path
        let newPath = tail;
        if (version.param !== 'latest') {
            // Add new version prefix
            // e.g. /foo -> /v4.9/foo
            // Handle root slash
            if (tail === '/') {
                newPath = `/${version.param}`;
            } else {
                newPath = `/${version.param}${tail}`;
            }
        }

        // Next.js router handles basePath automatically if configured correctly.
        // If we push '/v4.10', it becomes '/docs/v4.10'.
        router.push(newPath);
    };

    return (
        <select
            value={currentVersion}
            onChange={handleChange}
            className="p-2 bg-transparent text-sm font-medium rounded-md border border-fd-border hover:bg-fd-accent/50 focus:outline-none focus:ring-2 focus:ring-fd-ring text-fd-foreground w-full mb-4"
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
