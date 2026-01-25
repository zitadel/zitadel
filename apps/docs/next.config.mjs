import { createMDX } from 'fumadocs-mdx/next';
import { rehypeCode } from 'fumadocs-core/mdx-plugins';
import path from 'path';

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  basePath: '/docs',
  outputFileTracingIncludes: {
    '/**': ['./openapi/**/*', './content/**/*'],
  },
  async redirects() {
    return [
      {
        source: '/',
        destination: '/docs',
        basePath: false,
        permanent: false,
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
