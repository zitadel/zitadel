import { createMDX } from 'fumadocs-mdx/next';
import fs from 'node:fs';
import path from 'node:path';

const withMDX = createMDX();

// Load versions to configure rewrites
const versionsPath = path.join(process.cwd(), 'versions.json');
const versionsData = fs.existsSync(versionsPath) 
  ? JSON.parse(fs.readFileSync(versionsPath, 'utf8')) 
  : [];

// Load redirects
const redirectsPath = path.join(process.cwd(), 'redirects.json');
const redirectsData = fs.existsSync(redirectsPath)
  ? JSON.parse(fs.readFileSync(redirectsPath, 'utf8'))
  : [];

/** @type {import('next').NextConfig} */
const config = {
  reactStrictMode: true,
  serverExternalPackages: ['shiki'],
  async redirects() {
    return [
      ...redirectsData.map((r) => ({
        source: r.source,
        destination: r.destination,
        permanent: true,
      })),
      // Wildcard redirects for API sections
      {
        source: '/docs/apis/resources/admin/:path*',
        destination: '/docs/references/api-v1/admin',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/auth/:path*',
        destination: '/docs/references/api-v1/auth',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/mgmt/:path*',
        destination: '/docs/references/api-v1/management',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/system/:path*',
        destination: '/docs/references/api-v1/system',
        permanent: true,
      },
      // V2 API redirects
      {
        source: '/docs/apis/resources/user_service_v2/:path*',
        destination: '/docs/references/api/user',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/session_service_v2/:path*',
        destination: '/docs/references/api/session',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/org_service_v2/:path*',
        destination: '/docs/references/api/org',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/settings_service_v2/:path*',
        destination: '/docs/references/api/settings',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/action_service_v2/:path*',
        destination: '/docs/references/api/action',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/feature_service_v2/:path*',
        destination: '/docs/references/api/feature',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/idp_service_v2/:path*',
        destination: '/docs/references/api/idp',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/instance_service_v2/:path*',
        destination: '/docs/references/api/instance',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/project_service_v2/:path*',
        destination: '/docs/references/api/project',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/saml_service_v2/:path*',
        destination: '/docs/references/api/saml',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/oidc_service_v2/:path*',
        destination: '/docs/references/api/oidc',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/webkey_service_v2/:path*',
        destination: '/docs/references/api/webkey',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/authorization_service_v2/:path*',
        destination: '/docs/references/api/authorization',
        permanent: true,
      },
      {
        source: '/docs/apis/resources/application_service_v2/:path*',
        destination: '/docs/references/api/application',
        permanent: true,
      },
    ];
  },
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
