import { readFile } from 'node:fs/promises';
import { fileURLToPath, pathToFileURL } from 'node:url';
import path from 'node:path';

const TEXT_EXTENSIONS = ['.yaml', '.yml', '.conf', '.txt', '.json', '.caddyfile', '.go'];
const IMAGE_EXTENSIONS = ['.png', '.jpg', '.jpeg', '.gif', '.svg', '.webp', '.ico', '.bmp'];
const STYLE_EXTENSIONS = ['.css', '.scss', '.sass', '.less'];

export async function resolve(specifier, context, nextResolve) {
  const ext = path.extname(specifier.split('?')[0]).toLowerCase();

  try {
    return await nextResolve(specifier, context);
  } catch (err) {
    if (TEXT_EXTENSIONS.includes(ext) || IMAGE_EXTENSIONS.includes(ext) || STYLE_EXTENSIONS.includes(ext)) {
      // If it's an absolute path or looks like one, and it failed, we still want to load it
      // This happens on Vercel if files are missing on disk but referenced in MDX
      let url = specifier;
      if (!specifier.startsWith('file:')) {
        url = pathToFileURL(path.resolve(specifier)).href;
      }
      return {
        url,
        shortCircuit: true,
      };
    }
    throw err;
  }
}

export async function load(url, context, nextLoad) {
  // Only handle file URLs
  if (!url.startsWith('file:')) {
    return nextLoad(url, context);
  }

  const u = new URL(url);
  const ext = path.extname(u.pathname).toLowerCase();
  
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
