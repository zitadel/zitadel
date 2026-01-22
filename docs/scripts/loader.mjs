import { readFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const TEXT_EXTENSIONS = ['.yaml', '.yml', '.conf', '.txt', '.json', '.caddyfile', '.go'];
const IMAGE_EXTENSIONS = ['.png', '.jpg', '.jpeg', '.gif', '.svg', '.webp', '.ico', '.bmp'];
const STYLE_EXTENSIONS = ['.css', '.scss', '.sass', '.less'];

export async function load(url, context, nextLoad) {
  // Only handle file URLs
  if (!url.startsWith('file:')) {
    return nextLoad(url, context);
  }

  const u = new URL(url);
  const ext = path.extname(u.pathname).toLowerCase();
  // console.log(`[Loader] Loading: ${url}, ext: ${ext}`);
  
  // Handle ?raw imports
  if (u.searchParams.has('raw')) {
    u.search = '';
    const filePath = fileURLToPath(u);
    const content = await readFile(filePath, 'utf8');
    return {
      format: 'module',
      shortCircuit: true,
      source: `export default ${JSON.stringify(content)};`,
    };
  }
  
  if (TEXT_EXTENSIONS.includes(ext)) {
    u.search = '';
    const filePath = fileURLToPath(u);
    const content = await readFile(filePath, 'utf8');
    return {
      format: 'module',
      shortCircuit: true,
      source: `export default ${JSON.stringify(content)};`,
    };
  }

  if (IMAGE_EXTENSIONS.includes(ext)) {
    const mockImage = {
      src: url,
      height: 100,
      width: 100,
      blurDataURL: 'data:image/png;base64,',
    };
    return {
      format: 'module',
      shortCircuit: true,
      source: `export default ${JSON.stringify(mockImage)};`,
    };
  }

  if (STYLE_EXTENSIONS.includes(ext)) {
    return {
      format: 'module',
      shortCircuit: true,
      source: 'export default {};',
    };
  }
  
  // Try to use nextLoad, but if it fails for YAML, it's likely because it's from root
  try {
    return await nextLoad(url, context);
  } catch (err) {
    if (ext === '.yaml' || ext === '.yml') {
        const filePath = fileURLToPath(u);
        const content = await readFile(filePath, 'utf8');
        return {
          format: 'module',
          shortCircuit: true,
          source: `export default ${JSON.stringify(content)};`,
        };
    }
    throw err;
  }
}
