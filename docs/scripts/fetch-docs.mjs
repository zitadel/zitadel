import fs from 'fs';
import path, { join, dirname, resolve } from 'path';
import { spawn, execSync } from 'child_process';
import { fileURLToPath } from 'url';
import semver from 'semver';
import { Readable } from 'stream';

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = join(__dirname, '..');
const PROTO_DIR = join(ROOT_DIR, '../proto');
const CONTENT_DIR = join(ROOT_DIR, 'content');
const PUBLIC_DIR = join(ROOT_DIR, 'public');
const VERSIONS_FILE = join(ROOT_DIR, 'content/versions.json');
const CONTENT_LATEST_DIR = join(ROOT_DIR, 'content');

console.log(`[fetch-docs] __dirname: ${__dirname}`);
console.log(`[fetch-docs] ROOT_DIR: ${ROOT_DIR}`);
console.log(`[fetch-docs] PROTO_DIR: ${PROTO_DIR}`);
console.log(`[fetch-docs] CONTENT_DIR: ${CONTENT_DIR}`);

const REPO = 'zitadel/zitadel';
const CUTOFF = '2.0.0';
const ARCHIVE_URL = 'https://archive.zitadel.com';

async function fetchTags() {
  const token = process.env.GITHUB_TOKEN;
  const headers = {
    'User-Agent': 'node-fetch'
  };
  if (token) {
    headers['Authorization'] = `token ${token}`;
  }

  const url = `https://api.github.com/repos/${REPO}/tags?per_page=100`;
  console.log(`Fetching tags from ${url}...`);
  const res = await fetch(url, { headers });
  if (!res.ok) {
    const body = await res.text();
    throw new Error(`Failed to fetch tags: ${res.statusText} - ${body}`);
  }
  const tags = await res.json();
  console.log(`Fetched ${tags.length} tags.`);
  return tags;
}

function filterVersions(tags) {
  console.log(`Filtering tags with cutoff strictly > ${CUTOFF}...`);
  const versions = tags
    .map(t => t.name)
    .filter(v => {
        const valid = semver.valid(v);
        if (!valid) return false;
        // Strict cutoff: Do not fetch or build anything older (including that version)
        return semver.gt(v, CUTOFF);
    })
    .sort((a, b) => semver.rcompare(a, b));

  console.log(`Found ${versions.length} versions matching criteria.`);
  
  const groups = new Map();
  for (const v of versions) {
    const majorMinor = `${semver.major(v)}.${semver.minor(v)}`;
    if (!groups.has(majorMinor)) {
      groups.set(majorMinor, v);
    }
  }

  // User requested Highest + 2 minor versions below it = 3 versions total
  const result = Array.from(groups.values()).slice(0, 3);
  console.log(`Selected versions: ${result.join(', ')}`);
  return result;
}

// Downloads content from the 'fuma-docs' branch but puts it in a version-specific folder
// This is done because older tags do not have valid fumadocs content yet.
async function downloadVersion(tag) {
  // MOCK: Use fuma-docs branch content for validity testing
  const MOCK_REF = 'fuma-docs'; 
  const url = `https://github.com/${REPO}/archive/refs/heads/${MOCK_REF}.tar.gz`;
  
  const tempDir = join(ROOT_DIR, `.temp/${tag}`); // Extract to tag-specific temp to avoid collisions
  fs.mkdirSync(tempDir, { recursive: true });

  console.log(`Downloading content for ${tag} (using source: ${MOCK_REF})...`);

  const res = await fetch(url);
  if (!res.ok) throw new Error(`Failed to download ${url}: ${res.statusText}`);

  const tarArgs = [
    '-xz',
    '-C', tempDir,
    `--strip-components=1`,
    `zitadel-${MOCK_REF}/docs/content`,
    `zitadel-${MOCK_REF}/docs/public`, // Assuming public assets might be needed
    `zitadel-${MOCK_REF}/cmd/defaults.yaml`,
    `zitadel-${MOCK_REF}/cmd/setup/steps.yaml`
  ];

  await new Promise((resolve, reject) => {
    const tar = spawn('tar', tarArgs);
    Readable.fromWeb(res.body).pipe(tar.stdin);
    tar.on('close', (code) => (code === 0 ? resolve() : reject(new Error(`tar exited ${code}`))));
    tar.stderr.on('data', d => {
        const msg = d.toString();
        // Ignore "not found in archive" warnings if they are expected
        if (!msg.includes('Not found in archive')) console.error(msg);
    });
  });

  // Move to final destinations matches the version tag
  const versionSlug = `v${semver.major(tag)}.${semver.minor(tag)}`;
  
  const contentDest = join(CONTENT_DIR, versionSlug);
  const publicDest = join(PUBLIC_DIR, versionSlug);

  fs.mkdirSync(dirname(contentDest), { recursive: true });
  fs.mkdirSync(dirname(publicDest), { recursive: true });
  
  // Clean existing destination to avoid staleness
  fs.rmSync(contentDest, { recursive: true, force: true });
  fs.rmSync(publicDest, { recursive: true, force: true });

  if (fs.existsSync(join(tempDir, 'docs/content'))) {
     fs.renameSync(join(tempDir, 'docs/content'), contentDest);
  } else {
     console.warn(`[warn] docs/content not found in ${MOCK_REF} archive for ${tag}`);
  }

  // Handle external files (defaults.yaml etc)
  // We put them in _external folder inside the version content
  const externalDir = join(contentDest, '_external/cmd');
  fs.mkdirSync(externalDir, { recursive: true });
  
  if (fs.existsSync(join(tempDir, 'cmd/defaults.yaml'))) {
      fs.cpSync(join(tempDir, 'cmd/defaults.yaml'), join(externalDir, 'defaults.yaml'));
  }
   if (fs.existsSync(join(tempDir, 'cmd/setup/steps.yaml'))) {
      // Create setup dir if needed
      fs.mkdirSync(join(externalDir, 'setup'), { recursive: true });
      fs.cpSync(join(tempDir, 'cmd/setup/steps.yaml'), join(externalDir, 'setup/steps.yaml'));
  }

  // Also handling public assets? 
  // If fuma-docs branch has docs/public, we might want to version them or just copy them.
  // For now simple rename if exists
  if (fs.existsSync(join(tempDir, 'docs/public'))) {
    fs.renameSync(join(tempDir, 'docs/public'), publicDest);
  }

  fs.rmSync(tempDir, { recursive: true, force: true });
}



async function downloadFileContent(tagOrBranch, repoPath) {
    const url = `https://raw.githubusercontent.com/${REPO}/${tagOrBranch}/${repoPath}`;
    const res = await fetch(url);
    if (!res.ok) {
        // Fallback for some repo structures if needed, or if file doesn't exist in that version
        return null;
    }
    return await res.text();
}

async function fixRelativeImports(versionDir, tagOrBranch) {
    if (!fs.existsSync(versionDir)) return;
    const files = fs.readdirSync(versionDir, { recursive: true });
    
    // We'll traverse all files to fix links/imports
    for (const file of files) {
        const filePath = join(versionDir, file);
        if (!fs.statSync(filePath).isFile()) continue;
        if (!filePath.endsWith('.mdx') && !filePath.endsWith('.md')) continue;

        let content = fs.readFileSync(filePath, 'utf8');
        let changed = false;

        // Helper to rewrite a single path
        // Returns null if no change needed, or the new string if changed
        const rewritePath = (originalRelPath) => {
             // 1. Resolve where it was originally pointing
             const versionFolder = path.basename(versionDir);
             const relativePathInContent = filePath.split(join('content', versionFolder))[1];
             if (!relativePathInContent) return null; // Should not happen given logic

             // The original file was at CONTENT_LATEST_DIR + relativePathInContent
             const originalFilePath = join(CONTENT_LATEST_DIR, relativePathInContent);
             const originalDir = dirname(originalFilePath);
             
             // Resolve target
             // Note: resolve() handles '..' correctly
             const absoluteTarget = resolve(originalDir, originalRelPath);
             
             const projectRoot = resolve(ROOT_DIR, '..'); // Repo root
             
             // 2. Determine where this target lives now
             // Case A: It's in 'public' -> We moved public assets to 'public/<version>'
             // Check if target starts with ROOT_DIR/public
             if (absoluteTarget.startsWith(PUBLIC_DIR)) {
                 // It refers to a public asset.
                 // The NEW location of this asset is PUBLIC_DIR/<version>/<rest of path>
                 // But wait, PUBLIC_DIR is 'docs/public'.
                 // In downloadVersion, we move 'zitadel-xxx/docs/public' to 'docs/public/<version>'.
                 // So if original path was 'docs/public/img/foo.png', new path is 'docs/public/<version>/img/foo.png'.
                 
                 const relToPublic = absoluteTarget.slice(PUBLIC_DIR.length + 1); // 'img/foo.png'
                 const newTargetAbs = join(PUBLIC_DIR, versionFolder, relToPublic);
                 
                 // Now calculate relative path from the NEW file location to this NEW target
                 // New file is at 'filePath'
                 const newRelPath = path.relative(dirname(filePath), newTargetAbs);
                 return newRelPath.split(path.sep).join('/'); // Normalize to forward slashes
             }
             
             // Case B: It's in 'docs/content' (linking to another doc)
             // If it's a relative link to another doc, we usually want to keep it relative
             // content/<ver>/a.md -> content/<ver>/b.md is same relative relationship as content/a.md -> content/b.md
             // UNLESS it crosses out of content?
             
             // Case C: It's external (e.g. cmd/defaults.yaml)
             // This is what the original code was handling.
             // If it is NOT in docs/content, and NOT in docs/public, we treat it as external file to download.
             if (absoluteTarget.startsWith(projectRoot) && !absoluteTarget.startsWith(CONTENT_LATEST_DIR) && !absoluteTarget.startsWith(PublicOrRootPublic(absoluteTarget))) {
                  // External file logic...
                  // For now, let's keep the original logic for downloading defaults.yaml
                  return null; // We handle this in the regex pass below specially or reuse this logic?
             }

              // Case D: It's a relative link to, say, components/ or other things in docs/
              // Original: docs/content/foo.md -> ../components/Bar
              // New: docs/content/<ver>/foo.md -> ../../components/Bar
              // We just need to recalculate relative path from NEW file to ORIGINAL target
              
              // If it's pointing to something in docs/ (but not content/ or public/), like components/
              if (absoluteTarget.startsWith(ROOT_DIR) && !absoluteTarget.startsWith(CONTENT_LATEST_DIR)) {
                   const newRelPath = path.relative(dirname(filePath), absoluteTarget);
                   return newRelPath.split(path.sep).join('/');
              }
              
              return null;
        };

        // Helper for detecting public dir properly
        // The PUBLIC_DIR variable is 'docs/public'. But absoluteTarget might be resolved via 'src'
        function PublicOrRootPublic(p) {
            return PUBLIC_DIR; 
        }

        // --- Replacements ---
        
        // 1. Imports: import ... from '...'
        // We capture the path: group 2
        const importRegex = /(import\s+.*?\s+from\s+['"])([^'"]+)(['"])/g;
        content = content.replace(importRegex, (match, p1, p2, p3) => {
             if (!p2.startsWith('.')) return match; // Only relative
             const rewritten = rewritePath(p2);
             if (rewritten && rewritten !== p2) {
                 changed = true;
                 return `${p1}${rewritten}${p3}`;
             }
             
             // External file handling (legacy logic preserved/adapted)
             // If rewritePath returned null, maybe it is the "external download" case?
             // Let's implement the external download check here if rewritePath didn't handle it.
             // ... (The original logic specifically looked for `../../../../../cmd` etc)
             // We can incorporate it into rewritePath or do it here.
             
             // Let's check for the cmd/defaults.yaml case explicitly if rewritePath failed
             if (p2.includes('/cmd/')) {
                  // Re-use logic for downloading external files? 
                  // It's cleaner to put it in rewritePath, but I need async. replace does not support async.
                  // This function 'fixRelativeImports' is async. I can collect matches first.
             }
             
             return match;
        });
        
        // 2. Markdown Images: ![alt](src)
        const mdImgRegex = /(!\[.*?\]\()([^\)]+)(\))/g;
        content = content.replace(mdImgRegex, (match, p1, p2, p3) => {
             if (!p2.startsWith('.')) return match;
             const rewritten = rewritePath(p2);
             if (rewritten && rewritten !== p2) {
                 changed = true;
                 return `${p1}${rewritten}${p3}`;
             }
             return match;
        });

        // 3. HTML Attributes: src="..." or href="..."
        // Naive regex, but likely sufficient for MDX
        const htmlAttrRegex = /(src|href)=['"]([^'"]+)['"]/g;
        content = content.replace(htmlAttrRegex, (match, attr, val) => {
             if (!val.startsWith('.')) return match;
             const rewritten = rewritePath(val);
             if (rewritten && rewritten !== val) {
                 changed = true;
                 // reconstruct match
                 const quote = match.includes("'") ? "'" : '"';
                 return `${attr}=${quote}${rewritten}${quote}`;
             }
             return match;
        });
        
        // --- Special handling for the "download external file" case ---
        // The previous regex replacement structure prevents async operations.
        // We have to scan for external files separately or before replacing, OR use a synchronous download (not ideal)
        // OR rely on the existing synchronous logic for rewriting if we pre-download them.
        
        // Let's perform a SCAN for strictly external files (like ../cmd/defaults.yaml) first
        // Reuse original logic for downloading
        const importRegexForDownload = /import\s+.*\s+from\s+['"](\.\.\/(\.\.\/)+[^'"]+)['"]/g;
        let match;
        while ((match = importRegexForDownload.exec(content)) !== null) {
              const relPath = match[1];
              const versionFolder = path.basename(versionDir);
              const relativePathInContent = filePath.split(join('content', versionFolder))[1];
              const originalFilePath = join(CONTENT_LATEST_DIR, relativePathInContent);
              const absoluteImportTarget = resolve(dirname(originalFilePath), relPath);
              const projectRoot = resolve(ROOT_DIR, '..');
              
              // If it points to cmd/ or similar external
              if (absoluteImportTarget.startsWith(projectRoot) && !absoluteImportTarget.startsWith(CONTENT_LATEST_DIR) && !absoluteImportTarget.startsWith(PUBLIC_DIR)) {
                   const relativeToProjectRoot = absoluteImportTarget.replace(projectRoot + '/', '');
                   const localPathInVersion = join(versionDir, '_external', relativeToProjectRoot);
                   
                   // Download if missing
                   if (!fs.existsSync(localPathInVersion)) {
                      console.log(`[fix-imports] Downloading external: ${relativeToProjectRoot}`);
                      const fileContent = await downloadFileContent(tagOrBranch, relativeToProjectRoot); // existing function
                      if (fileContent) {
                          fs.mkdirSync(dirname(localPathInVersion), { recursive: true });
                          fs.writeFileSync(localPathInVersion, fileContent);
                      }
                   }
                   
                   // Rewrite to local path
                   const newRelPath = path.relative(dirname(filePath), localPathInVersion).split(path.sep).join('/');
                   const finalPath = newRelPath.startsWith('.') ? newRelPath : './' + newRelPath;
                   // Perform replacement
                   const newImport = match[0].replace(relPath, finalPath);
                   content = content.replace(match[0], newImport);
                   changed = true;
              }
        }
        
        if (changed) {
            console.log(`[fix-relative] Updated ${file}`);
            fs.writeFileSync(filePath, content);
        }
    }
}

function getLocalVersion() {
    const isVercel = process.env.VERCEL === '1';
    const vercelBranch = process.env.VERCEL_GIT_COMMIT_REF;
    
    let branch = vercelBranch;
    if (!branch) {
        try {
            branch = execSync('git branch --show-current').toString().trim();
        } catch (e) {}
    }

    if (branch && branch !== 'main') {
        return { label: branch, isUnreleased: true };
    }

    try {
        const tag = execSync('git describe --tags --abbrev=0').toString().trim();
        if (semver.valid(tag) && semver.gt(tag, CUTOFF)) {
            return { label: tag, isUnreleased: false };
        }
    } catch (e) {}

    return { label: 'v4.11.0', isUnreleased: true }; 
}

async function run() {
  console.log('Starting version discovery...');
  const tags = await fetchTags();
  const selectedTags = filterVersions(tags);
  
  let localVer = getLocalVersion();
  let others = selectedTags; // In our case, if local is latest, all filtered tags are others

  console.log(`Latest version (Local): ${localVer.label} (Unreleased: ${localVer.isUnreleased})`);
  console.log(`Older versions to fetch: ${others.join(', ') || 'None'}`);

  // Download chosen versions
  for (const tag of others) {
    const versionSlug = `v${semver.major(tag)}.${semver.minor(tag)}`;
    await downloadVersion(tag);
    await fixRelativeImports(join(CONTENT_DIR, versionSlug), tag);
  }

  // Generate versions.json
  const versionsJson = [
    { 
      param: 'latest', 
      label: localVer.isUnreleased ? `${localVer.label} (Unreleased)` : `${localVer.label} (Latest)`, 
      url: '/docs', 
      ref: 'local', 
      refType: 'local' 
    }
  ];

  for (const tag of others) {
    const v = `v${semver.major(tag)}.${semver.minor(tag)}`;
    const versionSlug = `v${semver.major(tag)}${semver.minor(tag)}x`;
    const targetUrl = `https://docs-git-${versionSlug}-zitadel.vercel.app/docs`;
    versionsJson.push({ 
        param: v, 
        label: v, 
        url: `/docs/${v}`, 
        ref: tag, 
        refType: 'tag',
        target: targetUrl
    });
  }

  versionsJson.push({
    label: `Archive (< ${CUTOFF})`,
    url: ARCHIVE_URL,
    type: 'external'
  });

  fs.writeFileSync(VERSIONS_FILE, JSON.stringify(versionsJson, null, 2));
  console.log('versions.json generated successfully.');
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
