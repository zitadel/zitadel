import defaultMdxComponents from 'fumadocs-ui/mdx';
import type { MDXComponents } from 'mdx/types';
import { APIPage } from '@/components/api-page';
import { Callout as FumaCallout } from 'fumadocs-ui/components/callout';
import { Tab, Tabs } from 'fumadocs-ui/components/tabs';

const Callout = (props: any) => <FumaCallout {...props} />;

export function useMDXComponents(components?: MDXComponents): MDXComponents {
  return {
    ...defaultMdxComponents,
    Callout,
    Tab,
    Tabs,
    APIPage,
    ...components,
  };
}

export const getMDXComponents = useMDXComponents;
