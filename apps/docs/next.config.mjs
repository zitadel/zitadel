import { createMDX } from 'fumadocs-mdx/next';
import { rehypeCode } from 'fumadocs-core/mdx-plugins';


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
        {
          source: '/img/:path*',
          destination: '/docs/img/:path*',
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
        {
          source: '/img/:path*',
          destination: '/docs/img/:path*',
          basePath: false,
          permanent: false,
        },
      ];
    }
  },
  async rewrites() {
    return [
      {
        source: '/favicon.ico',
        destination: 'https://zitadel.com/favicon.ico',
        basePath: false,
      },
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
      {
        source: '/pl/js/script.js',
        destination: 'https://plausible.io/js/script.js',
      },
      {
        source: '/pl/api/event',
        destination: 'https://plausible.io/api/event',
      },
    ];
  },
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'Content-Security-Policy',
            value: `default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://www.youtube.com https://www.google.com/recaptcha/ https://www.gstatic.com/recaptcha/ https://www.gstatic.com/charts/ https://www.youtube.com/; child-src zitadel.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://www.gstatic.com/charts/ zitadel.com; font-src 'self' https://fonts.gstatic.com https://fonts.googleapis.com; object-src 'none'; connect-src 'self' https://raw.githubusercontent.com/zitadel/ https://api.inkeep.com https://api.io.inkeep.com https://www.youtube.com; frame-src https://www.youtube.com/ https://www.google.com/recaptcha/ https://recaptcha.google.com/recaptcha/; img-src 'self' https://raw.githubusercontent.com/devicons/devicon/master/icons/ https://i.ytimg.com https://yt3.ggpht.com data: `,
          },
          {
            key: 'Strict-Transport-Security',
            value: 'max-age=63072000; includeSubDomains; preload',
          },
          {
            key: 'Permissions-Policy',
            value: 'payment=(self "https://js.stripe.com")',
          },
          {
            key: 'Referrer-Policy',
            value: 'origin-when-cross-origin',
          },
          {
            key: 'X-Frame-Options',
            value: 'SAMEORIGIN',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'X-XSS-Protection',
            value: '1; mode=block',
          },
        ],
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

const withMDX = createMDX();

export default withMDX(nextConfig);
