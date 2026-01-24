import fs from 'fs';
import { join, dirname } from 'path';
import { spawn, execSync } from 'child_process';
import { fileURLToPath } from 'url';
import semver from 'semver';
import { Readable } from 'stream';

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = join(__dirname, '..');
const PROTO_DIR = join(ROOT_DIR, '../proto');
const CONTENT_DIR = join(ROOT_DIR, 'content/versions');
const PUBLIC_DIR = join(ROOT_DIR, 'public/docs');
const VERSIONS_FILE = join(ROOT_DIR, 'content/versions.json');

const REPO = 'zitadel/zitadel';
const CUTOFF = '4.10.0';
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

  const result = Array.from(groups.values()).slice(0, 4);
  console.log(`Selected versions: ${result.join(', ')}`);
  return result;
}

async function downloadVersion(tag) {
  const url = `https://github.com/${REPO}/archive/refs/tags/${tag}.tar.gz`;
  const tempDir = join(ROOT_DIR, `.temp/${tag}`);
  fs.mkdirSync(tempDir, { recursive: true });

  const res = await fetch(url);
  if (!res.ok) throw new Error(`Failed to download ${url}: ${res.statusText}`);

  const tarArgs = [
    '-xz',
    '-C', tempDir,
    `--strip-components=1`,
    `zitadel-${tag.replace(/^v/, '')}/docs/content`,
    `zitadel-${tag.replace(/^v/, '')}/docs/public`,
    `zitadel-${tag.replace(/^v/, '')}/proto`
  ];

  await new Promise((resolve, reject) => {
    const tar = spawn('tar', tarArgs);
    Readable.fromWeb(res.body).pipe(tar.stdin);
    tar.on('close', (code) => (code === 0 ? resolve() : reject(new Error(`tar exited ${code}`))));
    tar.stderr.on('data', d => console.error(d.toString()));
  });

  // Move to final destinations
  const versionSlug = `v${semver.major(tag)}.${semver.minor(tag)}`;
  
  const contentDest = join(CONTENT_DIR, versionSlug);
  const protoDest = join(PROTO_DIR, versionSlug);
  const publicDest = join(PUBLIC_DIR, versionSlug);

  fs.rmSync(contentDest, { recursive: true, force: true });
  fs.rmSync(protoDest, { recursive: true, force: true });
  fs.rmSync(publicDest, { recursive: true, force: true });

  fs.mkdirSync(dirname(contentDest), { recursive: true });
  fs.mkdirSync(dirname(protoDest), { recursive: true });
  fs.mkdirSync(dirname(publicDest), { recursive: true });

  if (fs.existsSync(join(tempDir, 'docs/content'))) {
     fs.renameSync(join(tempDir, 'docs/content'), contentDest);
  }
  if (fs.existsSync(join(tempDir, 'proto'))) {
     fs.renameSync(join(tempDir, 'proto'), protoDest);
  }
  if (fs.existsSync(join(tempDir, 'docs/public'))) {
    fs.renameSync(join(tempDir, 'docs/public'), publicDest);
  }

  fs.rmSync(tempDir, { recursive: true, force: true });
}

async function downloadTestVersion() {
    const url = `https://github.com/${REPO}/archive/refs/heads/fuma-docs.tar.gz`;
    const tag = 'fuma-docs';
    const tempDir = join(ROOT_DIR, `.temp/${tag}`);
    fs.mkdirSync(tempDir, { recursive: true });

    console.log(`Downloading test version from ${url}...`);
    const res = await fetch(url);
    if (!res.ok) throw new Error(`Failed to download ${url}: ${res.statusText}`);

    const tarArgs = [
      '-xz',
      '-C', tempDir,
      `--strip-components=1`,
      `zitadel-${tag}/docs/content`,
      `zitadel-${tag}/proto`
    ];

    await new Promise((resolve, reject) => {
      const tar = spawn('tar', tarArgs);
      Readable.fromWeb(res.body).pipe(tar.stdin);
      tar.on('close', (code) => (code === 0 ? resolve() : reject(new Error(`tar exited ${code}`))));
      tar.stderr.on('data', d => console.error(d.toString()));
    });

    const versionSlug = 'vTest';
    const contentDest = join(CONTENT_DIR, versionSlug);
    const protoDest = join(PROTO_DIR, versionSlug);

    fs.rmSync(contentDest, { recursive: true, force: true });
    fs.rmSync(protoDest, { recursive: true, force: true });

    if (fs.existsSync(join(tempDir, 'docs/content'))) {
       fs.mkdirSync(dirname(contentDest), { recursive: true });
       fs.renameSync(join(tempDir, 'docs/content'), contentDest);
    } else {
       console.warn(`Warning: docs/content not found in ${tag} archive`);
    }
    if (fs.existsSync(join(tempDir, 'proto'))) {
       fs.mkdirSync(dirname(protoDest), { recursive: true });
       fs.renameSync(join(tempDir, 'proto'), protoDest);
    } else {
       console.warn(`Warning: proto folder not found in ${tag} archive`);
    }

    fs.rmSync(tempDir, { recursive: true, force: true });
}

function getLocalVersion() {
    try {
        const branch = execSync('git branch --show-current').toString().trim();
        if (branch === 'fuma-docs') return 'v4.11.0-beta';
        const tag = execSync('git describe --tags --abbrev=0').toString().trim();
        if (semver.valid(tag) && semver.gt(tag, CUTOFF)) return tag;
    } catch (e) {}
    return 'v4.11.0'; 
}

async function run() {
  console.log('Starting version discovery...');
  const tags = await fetchTags();
  const selectedTags = filterVersions(tags);
  
  let latestLabel = getLocalVersion();
  let others = selectedTags; // In our case, if local is latest, all filtered tags are others

  console.log(`Latest version (Local): ${latestLabel}`);
  console.log(`Older versions to fetch: ${others.join(', ') || 'None'}`);

  // Source Strategy: Latest (Local)
  const targetLatestProto = join(ROOT_DIR, '../proto/latest');
  fs.rmSync(targetLatestProto, { recursive: true, force: true });
  fs.mkdirSync(targetLatestProto, { recursive: true });
  
  const localProtoDir = join(ROOT_DIR, '../proto');
  const files = fs.readdirSync(localProtoDir);
  for (const file of files) {
      if (file === 'latest' || (file.startsWith('v') && semver.valid(file))) continue;
      
      const src = join(localProtoDir, file);
      const dest = join(targetLatestProto, file);
      
      if (fs.lstatSync(src).isDirectory()) {
          try {
              fs.symlinkSync(src, dest);
          } catch (e) {
              fs.cpSync(src, dest, { recursive: true });
          }
      } else {
          fs.copyFileSync(src, dest);
      }
  }

  // Download older versions
  for (const tag of others) {
    await downloadVersion(tag);
  }

  // Download Test version
  await downloadTestVersion();

  // Generate versions.json
  const versionsJson = [
    { param: 'latest', label: `${latestLabel} (Latest)`, url: '/docs' },
    { param: 'vTest', label: 'vTest (Branch)', url: '/docs/vTest' }
  ];

  for (const tag of others) {
    const v = `v${semver.major(tag)}.${semver.minor(tag)}`;
    versionsJson.push({ param: v, label: v, url: `/docs/${v}` });
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
