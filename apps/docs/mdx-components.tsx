import defaultMdxComponents from 'fumadocs-ui/mdx';
import type { MDXComponents } from 'mdx/types';
import type { GeneratedPageProps } from 'fumadocs-openapi';
import { OpenAPIPage } from '@/components/api-page';
import { TerminologyUpdate } from '@/components/terminology-update';
import { Callout } from 'fumadocs-ui/components/callout';
import { Tab, Tabs } from 'fumadocs-ui/components/tabs';
import { Step, Steps } from 'fumadocs-ui/components/steps';
import Admonition from '@/components/docusaurus/admonition';
import { GithubCodeBlock } from '@/components/github-code-block';
import { PreventLockout } from '@/components/prevent-lockout';

export function useMDXComponents(components?: MDXComponents): MDXComponents {
  return {
    ...defaultMdxComponents,
    Callout,
    Admonition,
    Tab,
    Tabs,
    Step,
    Steps,
    // Rendered by the generated API reference pages. Documents aren't preloaded,
    // `OpenAPIPage` loads them by ID (see lib/openapi.ts).
    OpenAPIPage: (props: GeneratedPageProps) => (
      <OpenAPIPage {...props} preloaded={{ docs: {} }} />
    ),
    TerminologyUpdate,
    GithubCodeBlock,
    PreventLockout,
    ...components,
  };
}

export const getMDXComponents = useMDXComponents;
