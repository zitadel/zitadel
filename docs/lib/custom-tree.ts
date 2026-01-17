import type * as PageTree from 'fumadocs-core/page-tree';
import { guidesSidebar, apisSidebar } from './sidebar-data';

type SidebarItem = {
    type?: 'category' | 'link' | 'doc';
    label?: string;
    items?: readonly any[]; 
    href?: string;
    id?: string;
    collapsed?: boolean;
} | string;

export function buildCustomTree(originalTree: PageTree.Root): PageTree.Root {
    const allPages = new Map<string, PageTree.Item>();
    const allFolders = new Map<string, PageTree.Folder>();

    function collect(node: PageTree.Node) {
        if (node.type === 'page') {
            allPages.set(node.url, node);
        }
        if (node.type === 'folder') {
            // Index the folder by its "path" (approximated from children or name?)
            // Fumadocs folder doesn't always have a URL.
            // But we can approximate path if it has an index page?
            if (node.index) {
                // If index page is /docs/foo/index, folder is /docs/foo
                 const url = node.index.url;
                 if(url) {
                    const folderPath = url.replace(/\/$/, '').replace(/\/index$/, '');
                    allFolders.set(folderPath, node);
                 }
            } else {
                 // Try to guess from first child? Or maybe name?
                 // This is tricky. 
                 // But for "reference/api/user", we know the path structure.
                 // Let's rely on finding children matching the path.
            }
            
            // Also we can traverse children to find match
            node.children.forEach(collect);
        }
    }
    originalTree.children.forEach(collect);
    
    // Helper to find folder by path scanning
    function findFolder(pathHint: string): PageTree.Folder | undefined {
        if (!pathHint) return undefined;
        let clean = pathHint;
        if (clean.startsWith('/')) clean = clean.substring(1);
        if (clean.startsWith('docs/')) clean = clean.substring(5);
        clean = clean.replace(/\/index$/, '').replace(/\/$/, '');
        
        const targetUrl = `/docs/${clean}`;
        // Normalize for fuzzy comparison
        const targetUrlNorm = targetUrl.replace(/_/g, '-').toLowerCase();

        let found: PageTree.Folder | undefined;
        
        function scan(node: PageTree.Node) {
            if (found) return;
            if (node.type === 'folder') {
                // Check index page
                if (node.index) {
                    const idxUrl = node.index.url.replace(/\/index$/, '');
                     if (idxUrl === targetUrl) { found = node; return; }
                     // Fuzzy match index
                     if (idxUrl.replace(/_/g, '-').toLowerCase() === targetUrlNorm) { found = node;  return; }
                }
                
                // If checking children pages for common prefix might be too aggressive/expensive here?
                // But let's check exact children match if no index
                const firstPage = node.children.find((c: PageTree.Node) => c.type === 'page') as PageTree.Item;
                if (firstPage) {
                    const pageDir = firstPage.url.substring(0, firstPage.url.lastIndexOf('/'));
                     if (pageDir === targetUrl) { found = node; return; }
                     if (pageDir.replace(/_/g, '-').toLowerCase() === targetUrlNorm) { found = node; return; }
                }
                
                node.children.forEach(scan);
            }
        }
        originalTree.children.forEach(scan);
        return found;
    }

    function findPage(path: string): PageTree.Item | undefined {
        if (!path) return undefined;
        let clean = path;
        if (clean.startsWith('/')) clean = clean.substring(1);
        if (clean.startsWith('docs/')) clean = clean.substring(5);
        clean = clean.replace(/\/index$/, '').replace(/\/$/, '');

        // 1. Try exact lookup first (fastest)
        const exactUrl = `/docs/${clean}`;
        if (allPages.has(exactUrl)) return allPages.get(exactUrl);

        // 2. Fuzzy match: normalize separators
        // We match if the *end* of the URL segments matches the path segments
        const cleanSegments = clean.split('/').filter(Boolean);
        
        for (const [url, node] of allPages.entries()) {
             const urlPath = url.startsWith('/docs/') ? url.substring(6) : url;
             const urlSegments = urlPath.split('/').filter(Boolean);
             
             // Check if cleanSegments match end of urlSegments
             if (urlSegments.length < cleanSegments.length) continue;
             
             const offset = urlSegments.length - cleanSegments.length;
             let match = true;
             for (let i = 0; i < cleanSegments.length; i++) {
                 const s1 = cleanSegments[i].replace(/_/g, '-').toLowerCase();
                 const s2 = urlSegments[i + offset].replace(/_/g, '-').toLowerCase();
                 if (s1 !== s2) {
                     match = false;
                     break;
                 }
             }
             
             if (match) return node;
        }
        
        return undefined;
    }

    function buildNode(item: SidebarItem): PageTree.Node | null {
        if (typeof item === 'string') {
            // Could be page or folder
            const page = findPage(item);
            if (page) return page;
            
            const folder = findFolder(item);
            if (folder) return folder;
            
            console.warn(`[Sidebar] Missing item for path: ${item}`);
            return null;
        }
        
        if (item.type === 'link') {
            return {
                type: 'page',
                name: item.label || 'Link',
                url: item.href || '#',
                external: true
            } as PageTree.Item;
        }

        if (item.type === 'doc') {
             const node = findPage(item.id || '');
             if (!node) return null;
             return {
                 ...node,
                 name: item.label || node.name
             };
        }

        if (item.type === 'category' || !item.type) {
            const children: PageTree.Node[] = [];
            if (item.items) {
                for (const child of item.items) {
                    const built = buildNode(child);
                    if (built) children.push(built);
                }
            }
            
            return {
                type: 'folder',
                name: item.label || 'Category',
                children: children,
                defaultOpen: !item.collapsed 
            } as PageTree.Folder;
        }
        
        return null; 
    }

    const guidesFolder = buildNode({
        type: 'category',
        label: 'Guides',
        items: guidesSidebar as any
    });

    const apisFolder = buildNode({
         type: 'category',
         label: 'APIs',
         items: apisSidebar as any
    });
    
    // We can also have legal, etc. if needed.
    const legalFolder = findFolder('legal');
    
    const newChildren: PageTree.Node[] = [];
    const indexNode = findPage('index');
    if (indexNode) newChildren.push(indexNode);
    
    if (guidesFolder) newChildren.push(guidesFolder);
    if (apisFolder) newChildren.push(apisFolder);
    
    // Self-hosting is in Guides usually? In new nav "Deploy & Operate" is under GuidesSidebar.
    // Concepts is in GuidesSidebar.
    
    // If legal exists, add it
    // if (legalFolder) newChildren.push(legalFolder);
    // Actually Legal was a link in sidebar?
    // "type": "link", "label": "Rate Limits (Cloud)", "href": "/legal/policies/rate-limit-policy"
    
    // If there are top-level pages not covered, we might lose them.
    // But user requested "make sure ONLY this navigation structure is used".
    
    return {
        name: originalTree.name,
        children: newChildren
    };
}
