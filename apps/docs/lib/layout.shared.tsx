import type { BaseLayoutProps } from 'fumadocs-ui/layouts/shared';
import { BookOpen, Compass, Code2, AppWindow, Server } from 'lucide-react';

export function baseOptions(): BaseLayoutProps {
  return {
    nav: {
      title: 'ZITADEL Docs',
    },
    links: [],
  };
}
