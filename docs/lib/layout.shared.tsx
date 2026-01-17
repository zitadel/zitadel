import type { BaseLayoutProps } from 'fumadocs-ui/layouts/shared';
import { BookOpen, Compass, Code2, AppWindow, Server } from 'lucide-react';

export function baseOptions(): BaseLayoutProps {
  return {
    nav: {
      title: 'ZITADEL Docs',
    },
    links: [
      {
        text: 'Guides',
        url: '/docs/guides/overview',
        active: 'nested-url',
        icon: <BookOpen />,
      },
      {
        text: 'Concepts',
        url: '/docs/concepts/principles',
        active: 'nested-url',
        icon: <Compass />,
      },
      {
        text: 'APIs',
        url: '/docs/apis/introduction',
        active: 'nested-url',
        icon: <Code2 />,
      },
      {
        text: 'SDKs',
        url: '/docs/sdk-examples/introduction',
        active: 'nested-url',
        icon: <AppWindow />,
      },
      {
        text: 'Self-Hosting',
        url: '/docs/self-hosting/deploy/overview',
        active: 'nested-url',
        icon: <Server />,
      },
    ],
  };
}
