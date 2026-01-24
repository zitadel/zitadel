import type * as PageTree from 'fumadocs-core/page-tree';
import {
    guidesSidebar,
    apisSidebar,
    legalSidebar,
    type SidebarItem, // Import the shared type
} from './sidebar-data';

// --- Logic ---

export function buildCustomTree(originalTree: PageTree.Root, options?: { stripPrefix?: string, suppressWarnings?: boolean }): PageTree.Root {
    const start = performance.now();
    const pageLookup = new Map<string, PageTree.Item>();
    const folderLookup = new Map<string, PageTree.Folder>();

    /**
     * 1. Normalizer
     */
    function normalize(path: string): string {
        let p = path.toLowerCase();
        // Ensure leading slash for consistent prefix checking
        if (!p.startsWith('/')) p = '/' + p;

        const prefix = options?.stripPrefix ?
            (options.stripPrefix.startsWith('/') ? options.stripPrefix : '/' + options.stripPrefix).toLowerCase()
            : undefined;

        if (prefix) {
            if (p.startsWith(prefix)) {
                p = p.substring(prefix.length);
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
    function collect(node: PageTree.Node, path: string = '') {
        if (node.type === 'page') {
            pageLookup.set(node.url, node);
            pageLookup.set(normalize(node.url), node);
        }

        if (node.type === 'folder') {
            const nodeName = typeof node.name === 'string' ? node.name : '';
            const currentPath = path ? `${path}/${nodeName}` : nodeName;

            // Use URL of index or children to infer the real disk-like path
            const sampleUrl = node.index?.url || node.children.find(c => c.type === 'page')?.url;
            if (sampleUrl) {
                const folderUrl = node.index ? sampleUrl : sampleUrl.split('/').slice(0, -1).join('/');
                folderLookup.set(normalize(folderUrl), node);
                folderLookup.set(folderUrl.replace(/^\/|docs\//, ''), node);
            }

            if (node.index) {
                pageLookup.set(node.index.url, node.index);
                pageLookup.set(normalize(node.index.url), node.index);
            }

            node.children.forEach(c => collect(c, currentPath));
        }
    }
    originalTree.children.forEach(c => collect(c));

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
        }

        // Suffix Scan Fallback
        for (const [lookupKey, node] of pageLookup.entries()) {
            if (lookupKey.endsWith('/' + key)) return node;
        }

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
            const folder = findFolder(item);
            if (folder) {
                const nodes: PageTree.Node[] = [];
                const seenUrls = new Set<string>();

                // Add index first
                if (folder.index && folder.index.url) {
                    nodes.push(folder.index);
                    seenUrls.add(folder.index.url);
                }

                // Add children that are not already present
                folder.children.forEach(child => {
                    if (child.type === 'page' && child.url) {
                        if (!seenUrls.has(child.url)) {
                            nodes.push(child);
                            seenUrls.add(child.url);
                        }
                    } else {
                        // For folders, we rely on their distinct identity or recursive structure.
                        // Assuming folder children don't duplicate the parent index URL.
                        nodes.push(child);
                    }
                });

                // Sort unique nodes by name
                return nodes.sort((a, b) => {
                    const nameA = typeof a.name === 'string' ? a.name : '';
                    const nameB = typeof b.name === 'string' ? b.name : '';
                    return nameA.localeCompare(nameB);
                });
            }

            const page = findPage(item);
            if (page) return page;

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
                    return [...folder.children].sort((a, b) => {
                        const nameA = typeof a.name === 'string' ? a.name : '';
                        const nameB = typeof b.name === 'string' ? b.name : '';
                        return nameA.localeCompare(nameB);
                    });
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

    const end = performance.now();
    const duration = end - start;
    if (duration > 100) {
        console.log(`[CustomTree] Build time: ${duration.toFixed(2)}ms`);
    }

    return {
        name: originalTree.name,
        children: newChildren
    };
}