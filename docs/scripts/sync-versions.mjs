import fs from 'node:fs';
import path from 'node:path';
import { execSync } from 'node:child_process';

const versionsPath = path.join(process.cwd(), 'versions.json');
if (!fs.existsSync(versionsPath)) {
  console.log('No versions.json found, skipping sync.');
  process.exit(0);
}

const versions = JSON.parse(fs.readFileSync(versionsPath, 'utf8'));
const outDir = path.join(process.cwd(), 'content/versions');

// Clean up previous versions
if (fs.existsSync(outDir)) {
  fs.rmSync(outDir, { recursive: true, force: true });
}
fs.mkdirSync(outDir, { recursive: true });

versions.forEach((v) => {
  if (v.type === 'remote') {
    console.log(`Syncing ${v.version} (${v.label}) from ref: ${v.ref}...`);
    const versionDir = path.join(outDir, v.version);
    fs.mkdirSync(versionDir, { recursive: true });

    try {
      // Fetch tarball from GitHub and extract only the docs/content/docs folder
      // We use --strip-components to flatten the structure into content/versions/<version>/
      // Structure in repo: <root>/docs/content/docs
      // Structure in tarball: <repo-name>-<ref>/docs/content/docs
      // We strip 3 levels: <repo-name>-<ref>, docs, content
      
      const url = `https://codeload.github.com/zitadel/zitadel/tar.gz/${v.ref}`;
      // Note: This assumes the remote repo has the same structure docs/content/docs
      // If the remote version is older and has a different structure, this might need adjustment.
      execSync(`curl -L ${url} | tar -xz -C "${versionDir}" --strip-components=3 "*/docs/content/docs"`, {
        stdio: 'inherit',
      });
      console.log(`Successfully synced ${v.version}`);
    } catch (error) {
      console.error(`Failed to sync ${v.version}:`, error);
      // Optional: Fail build if a version is missing
      // process.exit(1); 
    }
  }
});
