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
  turbopack: {
    rules: {
      '**/*.yaml': {
        loaders: ['raw-loader'],
        as: '*.js',
      },
      '**/*.Caddyfile': {
        loaders: ['raw-loader'],
        as: '*.js',
      },
      '**/*.conf': {
        loaders: ['raw-loader'],
        as: '*.js',
      },
    },
  },
  webpack: (config, { webpack }) => {
    config.plugins.push(
      new webpack.NormalModuleReplacementPlugin(
        /^fumadocs-mdx:collections\/server$/,
        (resource) => {
          resource.request = path.join(process.cwd(), '.source/server.ts');
        }
      )
    );

    config.module.rules.push({
      resourceQuery: /raw/,
      type: 'asset/source',
    });
    config.module.rules.push({
      test: /\.ya?ml$/,
      use: 'yaml-loader',
      resourceQuery: { not: [/raw/] },
    });
    config.module.rules.push({
      test: /\.(Caddyfile|conf)$/,
      type: 'asset/source',
    });
    return config;
  },
};

export default withMDX(config);
