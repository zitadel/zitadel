import { readFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const TEXT_EXTENSIONS = ['.yaml', '.conf', '.txt', '.json'];
const IMAGE_EXTENSIONS = ['.png', '.jpg', '.jpeg', '.gif', '.svg', '.webp', '.ico', '.bmp'];
const STYLE_EXTENSIONS = ['.css', '.scss', '.sass', '.less'];

export async function load(url, context, nextLoad) {
  // Only handle file URLs
  if (!url.startsWith('file:')) {
    return nextLoad(url, context);
  }

  const u = new URL(url);
  
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
  
  const ext = path.extname(u.pathname).toLowerCase();
  
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
  
  return nextLoad(url, context);
}
