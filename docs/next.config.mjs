import { createMDX } from 'fumadocs-mdx/next';
import fs from 'node:fs';
import path from 'node:path';

const withMDX = createMDX();

// Load versions to configure rewrites
const versionsPath = path.join(process.cwd(), 'versions.json');
const versionsData = fs.existsSync(versionsPath) 
  ? JSON.parse(fs.readFileSync(versionsPath, 'utf8')) 
  : [];

/** @type {import('next').NextConfig} */
const config = {
  reactStrictMode: true,
  serverExternalPackages: ['shiki'],
  async rewrites() {
    return versionsData
      .filter(v => v.type === 'external' && v.url && v.target)
      .map(v => ({
        source: `${v.url}/:path*`,
        destination: `${v.target}/:path*`,
      }));
  },
};

export default withMDX(config);
