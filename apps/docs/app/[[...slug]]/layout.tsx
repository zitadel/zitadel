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

  const allPages = [...source.getPages(), ...versionSource.getPages()];
  const labelsMap = new Map<string, string>();
  allPages.forEach(p => {
    if (p.data.sidebar_label) {
      labelsMap.set(p.url, p.data.sidebar_label);
    }
  });

  let tree = source.pageTree;
  if (currentVersion !== 'latest') {
    // Hoist the version folder to root to flatten sidebar
    // versionSource.pageTree -> [ v4.10 folder, v4.9 folder ... ]
    // We want the children of the specific version folder.
    const children = (versionSource.pageTree as any).children || [];

    const versionFolder = children.find((node: any) => {
      if (node.type !== 'folder') return false;
      // Check if the folder's index page belongs to this version
      if (node.index?.url && node.index.url.includes(`/${currentVersion}`)) return true;
      // Fallback: check first child if index missing
      if (node.children?.[0]?.url?.includes(`/${currentVersion}`)) return true;
      return false;
    });

    if (versionFolder) {
      // Hoist children
      const hoistedTree = {
        name: versionFolder.name,
        children: versionFolder.children,
      } as any;

      // Apply custom sidebar structure (from local sidebar-data.ts)
      // We strip the prefix '/v4.10' so the lookup matches generic keys 'guides/...'
      // The logs showed the URLs start with /[version] not /docs/[version]
      const prefix = currentVersion === 'latest' ? '/docs' : `/${currentVersion}`;
      tree = buildCustomTree(hoistedTree, {
        stripPrefix: prefix,
        labels: labelsMap
      });
    } else {
      tree = buildCustomTree(versionSource.pageTree, { labels: labelsMap });
    }
  } else {
    tree = buildCustomTree(source.pageTree, { labels: labelsMap });
  }

  return (
    <DocsLayout
      tree={tree}
      {...baseOptions()}
      sidebar={{
        banner: (
          <VersionSelector />
        ),
      }}
    >
      {props.children}
    </DocsLayout>
  );
}
