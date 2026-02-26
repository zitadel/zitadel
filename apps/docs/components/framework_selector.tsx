'use client';

import { useSearchParams } from 'next/navigation';
import { Tabs } from 'fumadocs-ui/components/tabs';
import { ReactNode, Suspense, useMemo } from 'react';

const FRAMEWORKS = ['angular', 'dotnet', 'astro', 'django', 'expressjs', 'fastapi', 'fastify', 'flask', 'flutter', 'go', 'hono', 'laravel', 'nestjs', 'nextjs', 'nuxt', 'qwik', 'react', 'solidstart', 'spring', 'symfony', 'svelte', 'vue'];
const DISPLAY_NAMES = ['Angular', 'ASP.NET Core', 'Astro', 'Django', 'Express.js', 'FastAPI', 'Fastify', 'Flask', 'Flutter', 'Go', 'Hono', 'Laravel', 'NestJS', 'Next.js', 'Nuxt', 'Qwik', 'React', 'SolidStart', 'Spring', 'Symfony', 'Svelte', 'Vue.js'];

function FrameworkSelectorInner({ children }: { children: ReactNode }) {
    const searchParams = useSearchParams();
    
    const activeIndex = useMemo(() => {
        const query = searchParams.get('framework')?.toLowerCase();
        if (!query) return 0;

        // Strip dots/spaces to match 'nextjs' or 'vue' reliably
        const sanitizedQuery = query.replace(/[\s.]/g, '');
        const index = FRAMEWORKS.findIndex(f => 
            f === sanitizedQuery || f.replace(/[\s.]/g, '') === sanitizedQuery
        );
        sessionStorage.setItem('framework-select', sanitizedQuery);
        return index === -1 ? 0 : index;
    }, [searchParams]);

    return (
        <Tabs
            // Using a unique key based on the index forces the component 
            // to ignore any cached local storage state and use the new defaultIndex.
            key={`framework-group-${activeIndex}`} 
            items={DISPLAY_NAMES}
            defaultIndex={activeIndex}
            groupId="framework-select"
        >
            {children}
        </Tabs>
    );
}

export function FrameworkSelector(props: { children: ReactNode }) {
    return (
        <Suspense fallback={<div className="h-40 w-full animate-pulse bg-fd-muted rounded-lg" />}>
            <FrameworkSelectorInner {...props} />
        </Suspense>
    );
}
