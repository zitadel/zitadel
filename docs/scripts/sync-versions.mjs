import fs from 'node:fs';
import path from 'node:path';
import { execSync } from 'node:child_process';

// Check if running from root or docs dir
const versionsPath = fs.existsSync(path.join(process.cwd(), 'docs/versions.json')) 
  ? path.join(process.cwd(), 'docs/versions.json')
  : path.join(process.cwd(), 'versions.json');

if (!fs.existsSync(versionsPath)) {
  console.log(`No versions.json found at ${versionsPath}, skipping sync.`);
  process.exit(0);
}

// --- Dynamic Version Discovery ---
try {
  console.log('Reading git tags...');

  const tagsOutput = execSync('git tag -l "v*"', { encoding: 'utf8' });
  const tags = tagsOutput.split('\n').filter(Boolean);

  const semverRegex = /^v(\d+)\.(\d+)\.(\d+)$/;
  const parsedTags = tags.map(t => {
      const match = t.match(semverRegex);
      if (!match) return null;
      return { original: t, major: parseInt(match[1]), minor: parseInt(match[2]), patch: parseInt(match[3]) };
  }).filter(Boolean);

  // Group by major.minor, keep highest Patch
  const latestByMinor = {};
  parsedTags.forEach(t => {
      const key = `${t.major}.${t.minor}`;
      if (!latestByMinor[key] || 
          (t.patch > latestByMinor[key].patch)) {
          latestByMinor[key] = t;
      }
  });

  // Sort by Major DESC, Minor DESC
  const sorted = Object.values(latestByMinor).sort((a, b) => {
      if (a.major !== b.major) return b.major - a.major;
      return b.minor - a.minor;
  });

  // Take top 4 (Highest + 3 prior)
  const topVersions = sorted.slice(0, 4);

  const newRemoteVersions = topVersions.map(v => {
    // Generate slug for Vercel preview URL
    // Pattern: v4.6.x -> v46x
    const versionSlug = `v${v.major}${v.minor}x`;
    const targetUrl = `https://docs-git-${versionSlug}-zitadel.vercel.app/docs`;

    return {
      version: `v${v.major}.${v.minor}`,
      label: `v${v.major}.${v.minor}`,
      type: 'external',
      url: `/docs/v${v.major}.${v.minor}`,
      target: targetUrl
    };
  });

  // Read existing to preserve manually added external links
  const currentVersions = JSON.parse(fs.readFileSync(versionsPath, 'utf8'));
  
  // Preserve existing external versions that are not covered by the new discovery
  const newVersionKeys = new Set(newRemoteVersions.map(v => v.version));
  const existingExternalVersions = currentVersions.filter(v => v.type === 'external' && !newVersionKeys.has(v.version));

  // Merge: New Generated + Existing Manual External
  const finalVersions = [...newRemoteVersions, ...existingExternalVersions];

  // Write back to versions.json
  fs.writeFileSync(versionsPath, JSON.stringify(finalVersions, null, 2));
  console.log('Updated versions.json with:', finalVersions.map(v => v.version).join(', '));

} catch (error) {
  console.error('Failed to update versions from git tags:', error);
  // Continue with existing versions.json if dynamic update fails
}
// ---------------------------------

console.log('Version sync complete (configured for external proxy).');
