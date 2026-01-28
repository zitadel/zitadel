import { createMDX } from 'fumadocs-mdx/next';
import { rehypeCode } from 'fumadocs-core/mdx-plugins';
import path from 'path';

import { promises as fs } from 'fs';
import { URL } from 'url';

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  basePath: '/docs',
  outputFileTracingIncludes: {
    '/**': ['./openapi/**/*', './content/**/*', './.source/**/*'],
  },
  async redirects() {
    try {
        const redirectsPath = new URL('./redirects.json', import.meta.url);
        const redirectsContent = await fs.readFile(redirectsPath, 'utf-8');
        const generatedRedirects = JSON.parse(redirectsContent);
        
        return [
          ...generatedRedirects,
          {
            source: '/',
            destination: '/docs',
            basePath: false,
            permanent: false,
          },
        ];
    } catch (e) {
        console.warn('Could not load redirects.json', e);
        return [
            {
              source: '/',
              destination: '/docs',
              basePath: false,
              permanent: false,
            },
        ];
    }
  },
  async rewrites() {
    return [
      {
        source: '/mp/lib.min.js',
        destination: 'https://cdn.mxpnl.com/libs/mixpanel-2-latest.min.js',
      },
      {
        source: '/mp/lib.js',
        destination: 'https://cdn.mxpnl.com/libs/mixpanel-2-latest.js',
      },
      {
        source: '/mp/decide',
        destination: 'https://decide.mixpanel.com/decide',
      },
      {
        source: '/mp/:slug*',
        destination: 'https://api-eu.mixpanel.com/:slug*',
      },
    ];
  },
  turbopack: {
    rules: {
      '*.{go,yaml,Caddyfile,conf}': {
        loaders: ['raw-loader'],
        as: '*.js',
      },
    },
  },
  webpack: (config) => {
    config.module.rules.push({
      test: /\.(go|yaml|Caddyfile|conf)$/,
      type: 'asset/source',
    });
    return config;
  },
};

const withMDX = createMDX({
  mdxOptions: {
    rehypePlugins: [
      [
        rehypeCode,
        {
          langs: [
            'bash',
            'yaml',
            'json',
            'go',
            'typescript',
            'javascript',
            'sql',
            'prometheus',
            'promql',
          ],
        },
      ],
    ],
  },
});

export default withMDX(nextConfig);
