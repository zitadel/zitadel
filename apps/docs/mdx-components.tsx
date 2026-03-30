import defaultMdxComponents from 'fumadocs-ui/mdx';
import type { MDXComponents } from 'mdx/types';
import { APIPage } from '@/components/api-page';
import { TerminologyUpdate } from '@/components/terminology-update';
import { Callout } from 'fumadocs-ui/components/callout';
import { Tab, Tabs } from 'fumadocs-ui/components/tabs';
import { Step, Steps } from 'fumadocs-ui/components/steps';
import Admonition from '@/components/docusaurus/admonition';

export function useMDXComponents(components?: MDXComponents): MDXComponents {
  return {
    ...defaultMdxComponents,
    Callout,
    Admonition,
    Tab,
    Tabs,
    Step,
    Steps,
    APIPage,
    TerminologyUpdate,
    ...components,
  };
}

export const getMDXComponents = useMDXComponents;
