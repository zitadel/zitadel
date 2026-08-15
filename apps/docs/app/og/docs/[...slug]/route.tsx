import { ImageResponse } from 'next/og';
import { notFound } from 'next/navigation';
import { getAllDocPages, getPage, getPageImage } from '@/lib/source';
import { generate as DefaultImage } from 'fumadocs-ui/og';

export const revalidate = false;
export const dynamicParams = false;
export const dynamic = 'force-static';

export async function GET(_request: Request, context: any) {
  const { slug } = await context.params;
  if (slug[slug.length - 1] !== 'image.png') notFound();

  const { page } = getPage(slug.slice(0, -1));
  if (!page) notFound();

  return new ImageResponse(
    <DefaultImage
      title={page.data.title}
      description={page.data.description}
      site="ZITADEL Docs"
    />,
    {
      width: 1200,
      height: 630,
    },
  );
}

export function generateStaticParams() {
  return getAllDocPages().map((page) => ({
    slug: getPageImage(page).segments,
  }));
}
