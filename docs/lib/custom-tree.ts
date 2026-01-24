import type * as PageTree from 'fumadocs-core/page-tree';
import {
    guidesSidebar,
    apisSidebar,
    legalSidebar,
    type SidebarItem, // Import the shared type
} from './sidebar-data';

// --- Logic ---

export function buildCustomTree(originalTree: PageTree.Root, options?: { stripPrefix?: string, suppressWarnings?: boolean }): PageTree.Root {
    const pageLookup = new Map<string, PageTree.Item>();
    const folderLookup = new Map<string, PageTree.Folder>();

    /**
     * 1. Normalizer
     */
    function normalize(path: string): string {
        let p = path.toLowerCase();

        if (options?.stripPrefix) {
            if (p.startsWith(options.stripPrefix)) {
                p = p.substring(options.stripPrefix.length);
            }
        } else {
            p = p.replace(/^\/|docs\/|\/$/g, '');
        }

        return p
            .replace(/^\//, '')
            .replace(/\/index$/, '')
            .replace(/_/g, '-');
    }

    /**
     * 2. Collector
     */
    function collect(node: PageTree.Node) {
        if (node.type === 'page') {
            pageLookup.set(node.url, node);
            pageLookup.set(normalize(node.url), node);
        }

        if (node.type === 'folder') {
            if (node.index) {
                pageLookup.set(node.index.url, node.index);
                pageLookup.set(normalize(node.index.url), node.index);
            }

            const childUrl = node.index?.url || node.children.find(c => c.type === 'page')?.url;
            if (childUrl) {
                const rawPath = childUrl.split('/').slice(0, -1).join('/');
                folderLookup.set(normalize(rawPath), node);

                const cleanDir = rawPath.replace(/^\/|docs\//, '');
                folderLookup.set(cleanDir, node);
            }
            node.children.forEach(collect);
        }
    }
    originalTree.children.forEach(collect);

    /**
     * 3. Optimized Page Finder
     */
    function findPage(path: string): PageTree.Item | undefined {
        if (!path) return undefined;

        const key = normalize(path);
        if (pageLookup.has(key)) return pageLookup.get(key);

        // Deduped Match
        const segments = key.split('/');
        if (segments.length >= 2 && segments[segments.length - 1] === segments[segments.length - 2]) {
            const dedupedKey = segments.slice(0, -1).join('/');
            if (pageLookup.has(dedupedKey)) return pageLookup.get(dedupedKey);
            if (path.includes('traefik')) console.log('[CustomTree] Dedup failed for:', dedupedKey, 'Orig:', key);
        }

        // Suffix Scan Fallback
        for (const [lookupKey, node] of pageLookup.entries()) {
            if (lookupKey.endsWith('/' + key)) return node;
        }

        if (path.includes('traefik')) console.log('[CustomTree] Failed to find completely:', key, 'Prefix:', options?.stripPrefix);

        return undefined;
    }

    /**
     * 4. Folder Finder
     */
    function findFolder(dirName: string): PageTree.Folder | undefined {
        if (!dirName) return undefined;

        if (folderLookup.has(dirName)) return folderLookup.get(dirName);
        const normKey = normalize(dirName);
        if (folderLookup.has(normKey)) return folderLookup.get(normKey);

        for (const [key, folder] of folderLookup.entries()) {
            if (key.includes(normKey)) return folder;
        }
        return undefined;
    }

    /**
     * 5. Recursive Builder
     */
    function buildNode(item: SidebarItem): PageTree.Node | PageTree.Node[] | null {
        // String Shorthand
        if (typeof item === 'string') {
            const page = findPage(item);
            if (page) return page;

            const folder = findFolder(item);
            if (folder) {
                // Safe sort for strict TS
                return [...folder.children].sort((a, b) =>
                    String(a.name ?? '').localeCompare(String(b.name ?? ''))
                );
            }

            if (!options?.suppressWarnings) {
                console.warn(`[Sidebar] Missing item: ${item}`);
            }
            return null;
        }

        // Links
        if (item.type === 'link') {
            const isExternal = item.href && (item.href.startsWith('http') || item.href.startsWith('//'));
            if (isExternal) {
                return {
                    type: 'page',
                    name: item.label || 'Link',
                    url: item.href || '#',
                    external: true,
                    icon: item.icon
                } as PageTree.Item;
            }
            if (item.href) {
                const page = findPage(item.href);
                if (page) return { ...page, name: item.label || page.name, icon: item.icon };
            }
            return null;
        }

        // Docs
        if (item.type === 'doc') {
            const node = findPage(item.id || '');
            if (!node) return null;
            return { ...node, name: item.label || node.name, icon: item.icon };
        }

        // Autogenerated
        if (item.type === 'autogenerated') {
            if (item.dirName) {
                const folder = findFolder(item.dirName);
                if (folder) {
                    // Safe sort for strict TS
                    return [...folder.children].sort((a, b) =>
                        String(a.name ?? '').localeCompare(String(b.name ?? ''))
                    );
                }
            }
            return null;
        }

        // Categories
        if (item.type === 'category' || !item.type) {
            const children: PageTree.Node[] = [];
            if (item.items) {
                for (const child of item.items) {
                    const built = buildNode(child);
                    if (built) {
                        if (Array.isArray(built)) children.push(...built);
                        else children.push(built);
                    }
                }
            }

            let indexPage: PageTree.Item | undefined;
            if (item.link?.type === 'generated-index' && item.link.slug) {
                indexPage = findPage(item.link.slug);
            } else if (item.link?.type === 'doc' && item.link.id) {
                indexPage = findPage(item.link.id);
            }

            return {
                type: 'folder',
                name: item.label || 'Category',
                children: children,
                defaultOpen: item.collapsed === false,
                index: indexPage,
                icon: item.icon
            } as PageTree.Folder;
        }

        return null;
    }

    // 6. Construction
    const newChildren: PageTree.Node[] = [];

    const indexNode = findPage('index');
    if (indexNode) newChildren.push(indexNode);

    guidesSidebar.forEach(item => {
        const node = buildNode(item);
        if (node && !Array.isArray(node)) newChildren.push(node);
    });

    const apisFolder = buildNode({
        type: 'category',
        label: 'APIs',
        items: apisSidebar
    });
    if (apisFolder && !Array.isArray(apisFolder)) newChildren.push(apisFolder);

    legalSidebar.forEach(item => {
        const node = buildNode(item);
        if (node && !Array.isArray(node)) newChildren.push(node);
    });

    return {
        name: originalTree.name,
        children: newChildren
    };
}