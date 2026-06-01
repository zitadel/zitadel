import { getLLMText, source } from '@/lib/source';

export const revalidate = false;
export const dynamic = 'force-static';

export async function GET() {
  const scan = source.getPages().map(getLLMText);
  const scanned = await Promise.all(scan);

  return new Response(scanned.join('\n\n'));
}
