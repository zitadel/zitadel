import fs from 'fs';
import path, { join, dirname, resolve } from 'path';
import { spawn, execSync } from 'child_process';
import { fileURLToPath } from 'url';
import semver from 'semver';
import { Readable } from 'stream';

const FALLBACK_VERSION = 'v4.10.0'; // Temporary fallback until > 4.10.0 exists
const FALLBACK_BRANCH = 'main'; // Primary branch to check for 4.10.0 content

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = join(__dirname, '..');
const PROTO_DIR = join(ROOT_DIR, '../../proto');
const CONTENT_DIR = join(ROOT_DIR, 'content');
const PUBLIC_DIR = join(ROOT_DIR, 'public');
const VERSIONS_FILE = join(ROOT_DIR, 'content/versions.json');
const CONTENT_LATEST_DIR = join(ROOT_DIR, 'content');

console.log(`[fetch-docs] __dirname: ${__dirname}`);
console.log(`[fetch-docs] ROOT_DIR: ${ROOT_DIR}`);
console.log(`[fetch-docs] PROTO_DIR: ${PROTO_DIR}`);
console.log(`[fetch-docs] CONTENT_DIR: ${CONTENT_DIR}`);

const REPO = 'zitadel/zitadel';
const CUTOFF = '4.10.0';

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

function getCurrentRef() {
  if (process.env.VERCEL_GIT_COMMIT_REF) return process.env.VERCEL_GIT_COMMIT_REF;
  if (process.env.GITHUB_REF_NAME) return process.env.GITHUB_REF_NAME;
  try {
    return execSync('git branch --show-current').toString().trim();
  } catch (e) {
    return FALLBACK_BRANCH;
  }
}

// Helper to determine the best branch for v4.10.0 content
function getFallbackSource() {
  // Prioritize CI/CD Environment Variables
  if (process.env.VERCEL_GIT_COMMIT_REF) {
    console.log(`[fallback] Detected Vercel Branch: ${process.env.VERCEL_GIT_COMMIT_REF}`);
    return process.env.VERCEL_GIT_COMMIT_REF;
  }
  if (process.env.GITHUB_REF_NAME) {
    console.log(`[fallback] Detected GitHub Action Branch: ${process.env.GITHUB_REF_NAME}`);
    return process.env.GITHUB_REF_NAME;
  }

  // Fallback for local dev if git is available
  try {
    const branch = execSync('git branch --show-current').toString().trim();
    if (branch) return branch;
  } catch (e) {
    // Ignore git errors
  }

  console.log(`[fallback] Defaulting to ${FALLBACK_BRANCH}`);
  return FALLBACK_BRANCH;
}

// sourceRef: Can be a tag (v1.2.3) or a branch (main, fuma-docs)
async function downloadVersion(tag, sourceRef) {
  const currentRef = getCurrentRef();
  const isLocal = sourceRef === currentRef;
  const tempDir = join(ROOT_DIR, `.temp/${tag}`);
  fs.mkdirSync(tempDir, { recursive: true });

  if (isLocal) {
    console.log(`[local] Copying local content for ${tag} (ref: ${sourceRef})...`);
    // Copy apps/docs/content while avoiding recursion
    const localContent = join(ROOT_DIR, 'content');
    const tempContent = join(tempDir, 'apps/docs/content');
    fs.mkdirSync(tempContent, { recursive: true });

    const contentItems = fs.readdirSync(localContent);
    for (const item of contentItems) {
      // Avoid copying versioned folders (starting with v) to prevent recursion if tag is current
      if (item.startsWith('v') && fs.statSync(join(localContent, item)).isDirectory()) continue;
      if (item === 'versions.json') continue;
      fs.cpSync(join(localContent, item), join(tempContent, item), { recursive: true });
    }

    // Copy public directory
    const localPublic = join(ROOT_DIR, 'public');
    const tempPublic = join(tempDir, 'apps/docs/public');
    fs.mkdirSync(tempPublic, { recursive: true });
    fs.cpSync(localPublic, tempPublic, { recursive: true });

    // Copy external files from repo root
    const repoRoot = resolve(ROOT_DIR, '../..');
    const tempCmd = join(tempDir, 'cmd');
    fs.mkdirSync(tempCmd, { recursive: true });
    if (fs.existsSync(join(repoRoot, 'cmd/defaults.yaml'))) {
      fs.cpSync(join(repoRoot, 'cmd/defaults.yaml'), join(tempCmd, 'defaults.yaml'));
    }
    if (fs.existsSync(join(repoRoot, 'cmd/setup/steps.yaml'))) {
      fs.mkdirSync(join(tempCmd, 'setup'), { recursive: true });
      fs.cpSync(join(repoRoot, 'cmd/setup/steps.yaml'), join(tempCmd, 'setup/steps.yaml'));
    }
  } else {
    const isBranch = sourceRef === 'main' || sourceRef === 'master' || !sourceRef.startsWith('v');
    const typeSegment = isBranch ? 'heads' : 'tags';
    const url = `https://github.com/${REPO}/archive/refs/${typeSegment}/${sourceRef}.tar.gz`;

    console.log(`Downloading content for ${tag} (using source: ${sourceRef})...`);

    const res = await fetch(url);
    if (!res.ok) throw new Error(`Failed to download ${url}: ${res.statusText}`);

    // ... rest of tar extraction logic

  const tarArgs = [
    '-xz',
    '-C', tempDir,
    `--strip-components=1`,
    `zitadel-${sourceRef.replace(/\//g, '-')}/apps/docs/content`, // GitHub archive naming usually matches ref name with slashes replaced? No, usually zitadel-ref. For tags it might vary.
    // Actually, GitHub archives top folder is usually `repo-ref`. 
    // For `v4.10.0` -> `zitadel-4.10.0`. For `fuma-docs` -> `zitadel-fuma-docs`.
    // We should probably rely on wildcard or just strict assumption.
    // Let's use a wildcard for the top directory since we strip it anyway?
    // wait, --strip-components=1 removes the top level folder, whatever it is.
    // BUT we are explicitly listing paths INSIDE that folder to extract.
    // `tar` requires separate arguments for paths.
    // If we don't know the exact top folder name, we can't restrict extraction efficiently without wildcards.
    // BUT standard `tar` doesn't always support wildcards in extraction list nicely.
    // Strategy: Extract the whole thing? No, too big.
    // Solution: The top folder name is predictable.
    // Branch: `zitadel-<branch-name>`
    // Tag: `zitadel-<tag-name>` (usually without 'v' if tag has 'v'? No, usually matches tag exactly).
    // Let's refine the top folder guess.
  ];
  
  // GitHub archive folder name logic:
  // Refs/heads/fuma-docs -> zitadel-fuma-docs
  // Refs/tags/v4.10.0 -> zitadel-4.10.0 (often strips 'v'?) OR zitadel-v4.10.0?
  // It's inconsistent. 
  // Better approach: Download and list first? Or extract everything and move?
  // Let's try to extract specific paths but use a wildcard if possible?
  // Actually, preventing the extract of the whole repo is good.
  // Let's guess:
  let topFolder = `zitadel-${sourceRef}`;
  // Tags often strip 'v' in the folder name if using 'archive/vX.Y.Z.tar.gz' vs 'archive/refs/tags/vX.Y.Z.tar.gz'
  // using refs/tags/... generally preserves it or uses the tag name.
  // Let's try to just extract *apps/docs/content* with a wildcard pattern if supported?
  // `*/apps/docs/content`
  
  // Update tarArgs to be safer with wildcards if using GNU tar, but macos bsdtar differs.
  // Safest: Extract all, then filter? Repository is large.
  // Let's stick to the previous code's assumption but make it dynamic.
  // Previous code: `zitadel-${MOCK_REF}/...`
  // If sourceRef is 'v4.10.0', try `zitadel-4.10.0` or `zitadel-v4.10.0`.
  // Let's just assume `zitadel-${sourceRef}` or `zitadel-${sourceRef.replace(/^v/, '')}`
  
  // Hack: We can peak at the first file entries? No.
  // Let's just try to extract `*/apps/docs/content` using --wildcards if on linux (which we are).
  
   const tarArgsWildcard = [
    '-xz',
    '-C', tempDir,
    '--strip-components=1',
    '--wildcards',
    '*/apps/docs/content',
    '*/apps/docs/public',
    '*/cmd/defaults.yaml',
    '*/cmd/setup/steps.yaml'
  ];

    await new Promise((resolve, reject) => {
      const tar = spawn('tar', tarArgsWildcard);
      Readable.fromWeb(res.body).pipe(tar.stdin);
      tar.on('close', (code) => (code === 0 ? resolve() : reject(new Error(`tar exited ${code}`))));
      tar.stderr.on('data', d => {
          const msg = d.toString();
           if (!msg.includes('Not found in archive')) console.error(msg);
      });
    });
  }

  // Name normalization logic for destination
  // If tag is 'v4.10.0', use only 'v4.10' for folder (matches current logic)
  const versionSlug = `v${semver.major(tag)}.${semver.minor(tag)}`;
  
  const contentDest = join(CONTENT_DIR, versionSlug);
  const publicDest = join(PUBLIC_DIR, versionSlug);

  fs.mkdirSync(dirname(contentDest), { recursive: true });
  fs.mkdirSync(dirname(publicDest), { recursive: true });
  
  fs.rmSync(contentDest, { recursive: true, force: true });
  fs.rmSync(publicDest, { recursive: true, force: true });

  // Move content
  if (fs.existsSync(join(tempDir, 'apps/docs/content'))) {
     fs.renameSync(join(tempDir, 'apps/docs/content'), contentDest);
  } else {
     // Fallback: sometimes unzipping structure might differ if wildcards matched differently?
     // Check if there is only one folder in tempDir?
     // With strip-components=1 and wildcards, it "should" land in tempDir/apps/docs/content directly 
     // IF the wildcards matched `topfolder/apps/docs/content`.
     // If `apps/docs/content` is not there, check for `content` root?
     // Let's verify commonly used structure.
     if (!fs.existsSync(join(tempDir, 'apps/docs/content'))) {
        console.warn(`[warn] apps/docs/content not found in archive for ${tag} (ref: ${sourceRef})`);
     }
  }

  // Handle external files (defaults.yaml etc)
  const externalDir = join(contentDest, '_external/cmd');
  fs.mkdirSync(externalDir, { recursive: true });
  
  if (fs.existsSync(join(tempDir, 'cmd/defaults.yaml'))) {
      fs.cpSync(join(tempDir, 'cmd/defaults.yaml'), join(externalDir, 'defaults.yaml'));
  }
   if (fs.existsSync(join(tempDir, 'cmd/setup/steps.yaml'))) {
      fs.mkdirSync(join(externalDir, 'setup'), { recursive: true });
      fs.cpSync(join(tempDir, 'cmd/setup/steps.yaml'), join(externalDir, 'setup/steps.yaml'));
  }

  if (fs.existsSync(join(tempDir, 'apps/docs/public'))) {
    fs.renameSync(join(tempDir, 'apps/docs/public'), publicDest);
  }

  fs.rmSync(tempDir, { recursive: true, force: true });
}



async function downloadFileContent(tagOrBranch, repoPath) {
    const currentRef = getCurrentRef();
    if (tagOrBranch === currentRef) {
        console.log(`[local] Reading local file content for: ${repoPath}`);
        const repoRoot = resolve(ROOT_DIR, '../..');
        const localPath = join(repoRoot, repoPath);
        if (fs.existsSync(localPath)) {
            return fs.readFileSync(localPath, 'utf8');
        }
        return null;
    }

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
             
             const projectRoot = resolve(ROOT_DIR, '../..'); // Repo root (apps/docs/../..)
             
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

             // Handle references to defaults.yaml or setup/steps.yaml specially
             if (originalRelPath.includes('cmd/defaults.yaml') || originalRelPath.includes('cmd/setup/steps.yaml')) {
                 return null; // Let the external file handlers deal with it
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
              const projectRoot = resolve(ROOT_DIR, '../..');
              
              // If it points to cmd/ or similar external
              if (absoluteImportTarget.startsWith(projectRoot) && !absoluteImportTarget.startsWith(CONTENT_LATEST_DIR) && !absoluteImportTarget.startsWith(PUBLIC_DIR)) {
                    const repoRoot = resolve(ROOT_DIR, '../..');
                    
                    // The original imports were relative to docs/content, not apps/docs/content.
                    // We need to account for the extra 'apps/' level.
                    let relativeToRepoRoot;
                    if (absoluteImportTarget.startsWith(join(repoRoot, 'apps'))) {
                        // If it resolved to apps/cmd/defaults.yaml, it should be cmd/defaults.yaml
                        relativeToRepoRoot = absoluteImportTarget.replace(join(repoRoot, 'apps') + '/', '');
                    } else {
                        relativeToRepoRoot = absoluteImportTarget.replace(repoRoot + '/', '');
                    }
                    
                    const localPathInVersion = join(versionDir, '_external', relativeToRepoRoot);
                   
                   // Download if missing
                   if (!fs.existsSync(localPathInVersion)) {
                      console.log(`[fix-imports] Downloading external: ${relativeToRepoRoot}`);
                      const fileContent = await downloadFileContent(tagOrBranch, relativeToRepoRoot); // existing function
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

    if (branch && branch !== 'main' && branch !== 'master') {
        return { label: branch, isUnreleased: true };
    }

    if (branch === 'main' || branch === 'master') {
        return { label: 'ZITADEL Docs', isUnreleased: false };
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
  
  // Conditional Fallback: If no versions found > 4.10.0, inject v4.10.0
  if (others.length === 0) {
      console.log(`[fallback] No versions found strictly > ${CUTOFF}. Injecting ${FALLBACK_VERSION} as fallback.`);
      others.push(FALLBACK_VERSION);
  }

  console.log(`Older versions to fetch: ${others.join(', ') || 'None'}`);

  // Download chosen versions
  // Parallelize download and processing
  await Promise.all(others.map(async (tag) => {
    let sourceRef = tag;
    
    // Explicit logic for version 4.10.0 to use active branch
    if (tag === FALLBACK_VERSION || tag === '4.10.0') {
         sourceRef = getFallbackSource();
    }
    
    const versionSlug = `v${semver.major(tag)}.${semver.minor(tag)}`;
    const contentDest = join(CONTENT_DIR, versionSlug);
    
    // Simple cache check: if directory exists and looks populated, skip
    // We could check for a specific file like meta.json or similar if we wanted to be more robust
    if (fs.existsSync(contentDest)) {
        console.log(`[skip] Version ${versionSlug} already exists. Skipping download.`);
    } else {
        await downloadVersion(tag, sourceRef);
        // Only fix imports if we just downloaded it (or maybe always run it? Safe to rerun)
        // Rerunning fixRelativeImports is relatively cheap compared to download + tar extraction
        // but let's stick to doing it only if we downloaded or if we force it.
        // For now: Always fix imports to be safe, or just on new download. 
        // Let's doing it on new download for speed.
        await fixRelativeImports(contentDest, tag);
    }
  }));

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

  fs.writeFileSync(VERSIONS_FILE, JSON.stringify(versionsJson, null, 2));
  console.log('versions.json generated successfully.');
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
