import { createMDX } from 'fumadocs-mdx/next';
import { rehypeCode } from 'fumadocs-core/mdx-plugins';

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
/*
  experimental: {
    turbo: {
      rules: {
        '** / *.{go,yaml,Caddyfile,conf}': {
          as: 'raw',
        },
      },
    },
  },
*/
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
