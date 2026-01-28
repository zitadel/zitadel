import { type BaseLayoutProps } from 'fumadocs-ui/layouts/shared';
import { BookOpen, Compass, Code2, AppWindow, Server } from 'lucide-react';

export const baseOptions: BaseLayoutProps = {
  nav: {
    title: 'ZITADEL Docs',
  },
  links: [
    {
      text: 'Guides',
      url: '/',
      active: 'nested-url',
      icon: <BookOpen />,
    },
    {
      text: 'APIs',
      url: '/apis/introduction',
      active: 'nested-url',
      icon: <Code2 />,
    },
    {
      text: 'Legal',
      url: '/legal/terms-of-service',
      active: 'nested-url',
      icon: <Server />, // Using Server icon as placeholder or pick another one
    },
  ],
  githubUrl: 'https://github.com/zitadel/zitadel',
};
