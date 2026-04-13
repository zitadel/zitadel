import { getPageImage, getPage, source } from '@/lib/source';
import {
  DocsBody,
  DocsPage,
  DocsTitle,
} from 'fumadocs-ui/layouts/docs/page';
import { notFound } from 'next/navigation';
import { getMDXComponents } from '@/mdx-components';
import type { Metadata } from 'next';
import { createRelativeLink } from 'fumadocs-ui/mdx';
import { Callout } from 'fumadocs-ui/components/callout';
import { Tab, Tabs } from 'fumadocs-ui/components/tabs';
import { Feedback } from '@/components/feedback';

export default async function Page(props: any) {
  const params = await props.params;
  const { page, source: pageSource } = getPage(params.slug);
  if (!page) notFound();

  const MDX = page.data.body;

  return (
    <DocsPage toc={page.data.toc} full={page.data.full}>
      <DocsTitle>{page.data.title}</DocsTitle>
      <DocsBody>
        <MDX
          components={getMDXComponents({
            Callout,
            Tab,
            Tabs,
            // this allows you to link to other pages with relative file paths
            a: createRelativeLink(pageSource, page),
          })}
        />
      </DocsBody>
      <Feedback />
    </DocsPage>
  );
}

export const dynamicParams = true;
export const revalidate = 3600;

export async function generateStaticParams() {
  return source.generateParams();
}

export async function generateMetadata(
  props: any,
): Promise<Metadata> {
  const params = await props.params;
  const { page } = getPage(params.slug);
  if (!page) notFound();
  const baseUrl = 'https://zitadel.com/docs';
  const url = params.slug ? `${baseUrl}/${params.slug.join('/')}` : baseUrl;

  let canonicalUrl = url;

  if (params.slug?.[0]?.startsWith('v')) {
    const unversionedSlug = params.slug.slice(1);
    const unversionedPage = source.getPage(unversionedSlug);
    if (unversionedPage) {
      canonicalUrl = `${baseUrl}${unversionedPage.url === '/' ? '' : unversionedPage.url}`;
    }
  }

  let description = page.data.description;
  if (!description) {
    description = `Explore ZITADEL documentation for ${page.data.title}. Learn how to integrate, manage, and secure your applications with our comprehensive identity and access management solutions.`;
  }

  return {
    title: page.data.title,
    description: description.length > 200 ? description.substring(0, 197) + '...' : description,
    alternates: {
      canonical: canonicalUrl,
    },
    openGraph: {
      images: getPageImage(page).url,
    },
  };
}
