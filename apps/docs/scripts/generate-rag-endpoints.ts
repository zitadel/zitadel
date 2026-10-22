import { writeFileSync, readFileSync, mkdirSync } from 'fs';
import { join, dirname, basename } from 'path';
import { fileURLToPath } from 'url';
import { glob } from 'glob';

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = join(__dirname, '..');
const OPENAPI_ROOT = join(ROOT_DIR, 'openapi', 'latest');
const CONTENT_API_ROOT = join(ROOT_DIR, 'content', 'reference', 'api');
const OUTPUT_FILE = join(ROOT_DIR, '.data', 'api-endpoints.txt');

async function generateRagContext() {
    console.log(`Scanning OpenAPI files in: ${OPENAPI_ROOT}`);

    // 1. Build an index of all actual generated MDX files on disk
    // Mapping operationId / file basename (e.g., 'zitadel.user.v2.UserService.AddHumanUser') -> URL
    const mdxFiles = await glob('**/*.{md,mdx}', { cwd: CONTENT_API_ROOT });
    const urlMap = new Map<string, string>();

    for (const file of mdxFiles) {
        const fileBase = basename(file, file.endsWith('.mdx') ? '.mdx' : '.md');
        const relUrl = `/docs/reference/api/${file.replace(/\.(md|mdx)$/, '')}`;

        // Index by lowercase filename for case-insensitive matching
        urlMap.set(fileBase.toLowerCase(), relUrl);
    }

    const specs = await glob('**/*.openapi.json', { cwd: OPENAPI_ROOT });
    const endpointBlocks: string[] = [];

    for (const specPath of specs) {
        const fullPath = join(OPENAPI_ROOT, specPath);
        const content = readFileSync(fullPath, 'utf8');
        const doc = JSON.parse(content);

        if (!doc.paths || Object.keys(doc.paths).length === 0) continue;

        for (const [pathUrl, methods] of Object.entries(doc.paths)) {
            for (const [method, details] of Object.entries(methods as Record<string, any>)) {

                let description = details.description || details.summary || 'No description provided.';
                description = description
                    .replace(/\n+/g, ' ')
                    .replace(/\[([^\]]+)\]\([^\)]+\)/g, '$1')
                    .trim();

                if (description.length > 300) {
                    description = description.slice(0, 297) + '...';
                }

                const operationId = details.operationId || 'N/A';

                // 2. Resolve URL by checking the file index built from the actual MDX files
                let docUrl = urlMap.get(operationId.toLowerCase());

                // Fallback search if not found directly by operationId
                if (!docUrl) {
                    const matchedKey = Array.from(urlMap.keys()).find(k => k.includes(operationId.toLowerCase()));
                    if (matchedKey) {
                        docUrl = urlMap.get(matchedKey);
                    }
                }

                // Final fallback if the page doesn't exist yet
                if (!docUrl) {
                    docUrl = `/docs/reference/api`;
                }

                const block = [
                    `API Endpoint: ${method.toUpperCase()} ${pathUrl}`,
                    `Operation ID: ${operationId}`,
                    `Description: ${description}`,
                    `Documentation URL: https://zitadel.com${docUrl}`
                ].join('\n');

                endpointBlocks.push(block);
            }
        }
    }

    mkdirSync(dirname(OUTPUT_FILE), { recursive: true });
    writeFileSync(OUTPUT_FILE, endpointBlocks.join('\n\n'), 'utf8');
    console.log(`Successfully generated ${endpointBlocks.length} endpoints to ${OUTPUT_FILE}`);
}

generateRagContext().catch(err => {
    console.error(err);
    process.exit(1);
});