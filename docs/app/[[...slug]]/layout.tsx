import { source, versionSource } from '@/lib/source';
import { DocsLayout } from 'fumadocs-ui/layouts/docs';
import { baseOptions } from '@/lib/layout.shared';
import { VersionSelector } from '@/components/version-selector';
import { buildCustomTree } from '@/lib/custom-tree';

export default async function Layout(props: { children: React.ReactNode; params: Promise<{ slug?: string[] }> }) {
  const params = await props.params;
  const slug = params.slug || [];
  
  let tree = source.pageTree;
  if (slug.length > 0 && slug[0].startsWith('v')) {
     tree = versionSource.pageTree;
  } else {
     tree = buildCustomTree(source.pageTree);
  }

  return (
    <DocsLayout 
      tree={tree} 
      {...baseOptions()}
      sidebar={{
        banner: <VersionSelector />,
      }}
    >
      {props.children}
    </DocsLayout>
  );
}
