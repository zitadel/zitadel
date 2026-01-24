import { source, versionSource } from '@/lib/source';
import { DocsLayout } from 'fumadocs-ui/layouts/docs';
import { baseOptions } from '@/lib/layout.shared';
import { buildCustomTree } from '@/lib/custom-tree';
import versions from '@/content/versions.json';
import { VersionSelector } from '@/components/version-selector';

export default async function Layout(props: { children: React.ReactNode; params: Promise<{ slug?: string[] }> }) {
  const params = await props.params;
  const slug = params.slug || [];

  // Determine current version and tail path
  // If first slug segment matches a version param (e.g. v4.10), that's the version.
  // Otherwise it's latest.
  const versionParam = slug[0] && versions.find(v => v.param === slug[0]) ? slug[0] : undefined;
  const currentVersion = versionParam || 'latest';

  let tree = source.pageTree;
  if (currentVersion !== 'latest') {
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
