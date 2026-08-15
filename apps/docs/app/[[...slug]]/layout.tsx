import { getAllDocPages, source, versionSource } from '@/lib/source';
import { ReactNode } from 'react';
import { DocsLayout } from 'fumadocs-ui/layouts/docs';
import { baseOptions } from '@/lib/layout.shared';
import { buildCustomTree } from '@/lib/custom-tree';
import rawVersions from '@/content/versions.json';
import { VersionSelector } from '@/components/version-selector';

type DocVersion = { param: string };
type PageTree = typeof source.pageTree;

const versions = rawVersions as DocVersion[];
const versionParams = new Set(versions.map((version) => version.param));

const labelsMap = new Map<string, string>();
for (const page of getAllDocPages()) {
  if (page.data.sidebar_label) {
    labelsMap.set(page.url, page.data.sidebar_label);
  }
}

const latestTree = buildCustomTree(source.pageTree, { labels: labelsMap });
const versionTreeCache = new Map<string, PageTree>();

function getVersionTree(currentVersion: string): PageTree {
  const cached = versionTreeCache.get(currentVersion);
  if (cached) return cached;

  const children = (versionSource.pageTree as any).children || [];
  const versionFolder = children.find((node: any) => {
    if (node.type !== 'folder') return false;
    if (node.index?.url && node.index.url.includes(`/${currentVersion}`)) return true;
    return Boolean(node.children?.[0]?.url?.includes(`/${currentVersion}`));
  });

  const tree = versionFolder
    ? buildCustomTree(
        {
          name: versionFolder.name,
          children: versionFolder.children,
        } as any,
        {
          stripPrefix: `/${currentVersion}`,
          labels: labelsMap,
        },
      )
    : buildCustomTree(versionSource.pageTree, { labels: labelsMap });

  versionTreeCache.set(currentVersion, tree);
  return tree;
}

export default async function Layout(props: { children: ReactNode; params: Promise<{ slug?: string[] }> }) {
  const params = await props.params;
  const slug = params.slug || [];

  const currentVersion = slug[0] && versionParams.has(slug[0]) ? slug[0] : 'latest';
  const tree = currentVersion === 'latest' ? latestTree : getVersionTree(currentVersion);

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
