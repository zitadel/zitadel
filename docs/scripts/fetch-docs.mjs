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
const PUBLIC_DIR = join(ROOT_DIR, 'public/docs');
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
    const downloadedFiles = new Set();

    for (const file of files) {
        const filePath = join(versionDir, file);
        if (!fs.statSync(filePath).isFile()) continue;
        if (filePath.endsWith('.mdx') || filePath.endsWith('.md')) {
            let content = fs.readFileSync(filePath, 'utf8');
            let changed = false;

            // Regex to find imports with relative paths going outside docs/content
            // Example: import DefaultsYamlSource from "../../../../../cmd/defaults.yaml";
            const importRegex = /import\s+.*\s+from\s+['"](\.\.\/(\.\.\/)+[^'"]+)['"]/g;
            let match;
            
            // We need to collect all matches first because we'll be changing content
            const matches = [];
            while ((match = importRegex.exec(content)) !== null) {
                matches.push({ full: match[0], path: match[1] });
            }

            for (const m of matches) {
                // To find the repo-relative path, we resolve the import as it would be in the original repo
                // Original content was in docs/content/<rest-of-path>
                // Versioned content is in docs/content/<version>/<rest-of-path>
                const versionFolder = path.basename(versionDir);
                const relativePathInContent = filePath.split(join('content', versionFolder))[1];
                if (!relativePathInContent) continue;
                
                // Use absolute paths to avoid confusion
                const originalFilePath = join(CONTENT_LATEST_DIR, relativePathInContent);
                const originalDir = dirname(originalFilePath);
                
                // Resolve the import against the original location
                // e.g. /abs/path/to/docs/content/guides/... + ../../../../../cmd/defaults.yaml
                const absoluteImportTarget = resolve(originalDir, m.path);
                const projectRoot = resolve(ROOT_DIR, '..'); // /abs/path/to/repo

                // Check if it's inside the project but outside docs/content (or strictly outside docs if we want to be safe)
                // The main use case is importing from cmd/ or other repo folders
                if (absoluteImportTarget.startsWith(projectRoot) && !absoluteImportTarget.startsWith(CONTENT_LATEST_DIR)) {
                     const relativeToProjectRoot = absoluteImportTarget.replace(projectRoot + '/', '');
                     
                     // We want to download this file and put it in a local versioned folder
                     const localPathInVersion = join(versionDir, '_external', relativeToProjectRoot);
                     const localDirInVersion = dirname(localPathInVersion);
                     
                     if (!fs.existsSync(localPathInVersion)) {
                         console.log(`[fix-imports] Discovered external import: ${relativeToProjectRoot} in ${file}`);
                         const fileContent = await downloadFileContent(tagOrBranch, relativeToProjectRoot);
                         if (fileContent !== null) {
                             fs.mkdirSync(localDirInVersion, { recursive: true });
                             fs.writeFileSync(localPathInVersion, fileContent);
                             downloadedFiles.add(relativeToProjectRoot);
                         } else {
                             console.warn(`[fix-imports] Failed to download versioned file: ${relativeToProjectRoot}`);
                             continue;
                         }
                     }
 
                     // Update the import path in the MDX file to point to our local copy
                     const newRelativePath = path.relative(dirname(filePath), localPathInVersion);
                     // MDX imports should use forward slashes
                     const normalizedPath = newRelativePath.split(path.sep).join('/');
                     const finalPath = normalizedPath.startsWith('.') ? normalizedPath : './' + normalizedPath;
                     
                     console.log(`[fix-imports] Rewriting import in ${file}: ${m.path} -> ${finalPath}`);
                     const newImport = m.full.replace(m.path, finalPath);
                     content = content.replace(m.full, newImport);
                     changed = true;
                }
            }

            if (changed) {
                fs.writeFileSync(filePath, content);
            }
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
