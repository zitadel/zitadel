import fs from 'fs';
import path from 'path';

const CONTENT_DIR = 'content';

function walk(dir, callback) {
  fs.readdirSync(dir).forEach((file) => {
    const filepath = path.join(dir, file);
    const stats = fs.statSync(filepath);
    if (stats.isDirectory()) {
      walk(filepath, callback);
    } else if (stats.isFile() && (filepath.endsWith('.mdx') || filepath.endsWith('.md'))) {
      callback(filepath);
    }
  });
}

walk(CONTENT_DIR, (filepath) => {
  let content = fs.readFileSync(filepath, 'utf8');
  let changed = false;

  const segments = filepath.split(path.sep);
  const depth = segments.length - 1; 
  const correctPrefix = '../'.repeat(depth) + 'public/';

  // Match any sequence of ../ followed by public/
  const relativeLinkRegex = /(\.\.\/)+public\//g;
  
  if (relativeLinkRegex.test(content)) {
     const newContent = content.replace(relativeLinkRegex, correctPrefix);
     if (newContent !== content) {
         content = newContent;
         changed = true;
     }
  }

  // Also catch any remaining absolute ones if they exist
  // Regex to match Markdown images: ![alt](/docs/path...)
  // and HTML images: src="/docs/path..."
  // Capturing group 1: '![...](', 2: path
  // OR src=", 2: path
  
  // Replaces /docs/img/ -> correctPrefix + img/
  // Replaces /docs/vX.Y/ -> correctPrefix + vX.Y/

  const PUBLIC_PREFIX = '/docs/img/';
  const VERSION_PREFIXES = ['/docs/v4.10/', '/docs/v4.9/', '/docs/v4.8/'];
  const TARGET_PATHS = [PUBLIC_PREFIX, ...VERSION_PREFIXES];

  // Pattern: (!\[.*?\]\(|src=["'])(/docs/(?:img|v\d+\.\d+)/.*?)(\)|["'])
  const pattern = /(!\[.*?\]\(|src=["'])\/docs\/((?:img|v\d+\.\d+)\/.*?)(\)|["'])/g;

  if (pattern.test(content)) {
      const newContent = content.replace(pattern, (match, prefix, pathPart, suffix) => {
          // prefix e.g. '![alt](' or 'src="'
          // pathPart e.g. 'img/foo.png' or 'v4.10/foo.png'
          // suffix e.g. ')' or '"'
          return `${prefix}${correctPrefix}${pathPart}${suffix}`;
      });
      
      if (newContent !== content) {
          content = newContent;
          changed = true;
      }
  }

  if (changed) {
    fs.writeFileSync(filepath, content);
    console.log(`Updated: ${filepath}`);
  }
});
