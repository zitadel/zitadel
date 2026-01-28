'use client';

import { useSearchParams } from 'next/navigation';
import { Tabs, Tab } from 'fumadocs-ui/components/tabs';
import { ReactNode } from 'react';

const FRAMEWORKS = ['angular', 'react', 'vue'];
const DISPLAY_NAMES = ['Angular', 'React', 'Vue'];

export function FrameworkSelector({ children }: { children: ReactNode }) {
    const searchParams = useSearchParams();
    const frameworkParam = searchParams.get('framework')?.toLowerCase();

    // Find the index based on the URL param, defaulting to 0 (Next.js)
    const activeIndex = frameworkParam && FRAMEWORKS.includes(frameworkParam)
        ? FRAMEWORKS.indexOf(frameworkParam)
        : 0;

    return (
        <Tabs
            items={DISPLAY_NAMES}
            defaultIndex={activeIndex}
            persist // Keeps selection if they navigate away and come back
            groupId="framework-select"
        >
            {children}
        </Tabs>
    );
}