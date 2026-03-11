import createNextIntlPlugin from "next-intl/plugin";
import { DEFAULT_CSP } from "./constants/csp.js";

const withNextIntl = createNextIntlPlugin();

const secureHeaders = [
  {
    key: "Strict-Transport-Security",
    value: "max-age=63072000; includeSubDomains; preload",
  },
  {
    key: "Referrer-Policy",
    value: "origin-when-cross-origin",
  },
  {
    key: "X-Content-Type-Options",
    value: "nosniff",
  },
  {
    key: "X-XSS-Protection",
    value: "1; mode=block",
  },
  {
    key: "Content-Security-Policy",
    value: DEFAULT_CSP,
  },
  { key: "X-Frame-Options", value: "deny" },
];

/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: process.env.NEXT_PUBLIC_BASE_PATH,
  output: process.env.NEXT_OUTPUT_MODE || undefined,
  // Set the tracing root to the monorepo workspace root so that the pnpm
  // virtual store (.pnpm) is included in the standalone output. Without this,
  // the standalone node_modules contain dangling symlinks that point outside
  // the Docker image. The server entry lands at apps/login/server.js inside
  // the standalone, which the build script and Dockerfile are adjusted for.
  outputFileTracingRoot: new URL("../..", import.meta.url).pathname,
  reactStrictMode: true,
  experimental: {
    // Add React 19 compatibility optimizations
    optimizePackageImports: ["@radix-ui/react-tooltip", "@heroicons/react"],
    useCache: true,
    serverActions: {
      ...(process.env.SERVER_ACTION_ALLOWED_ORIGINS
        ? { allowedOrigins: process.env.SERVER_ACTION_ALLOWED_ORIGINS.split(",").map((o) => o.trim()) }
        : {}),
    },
  },
  // Packages that must not be bundled by webpack and should remain as external
  // requires at runtime. These packages use native modules or have bundling
  // incompatibilities. Keep this list in sync with package.json dependencies
  // when adding new OpenTelemetry or logging packages.
  serverExternalPackages: [
    'winston',
    '@opentelemetry/api',
    '@opentelemetry/api-logs',
    '@opentelemetry/sdk-node',
    '@opentelemetry/sdk-metrics',
    '@opentelemetry/sdk-logs',
    '@opentelemetry/exporter-metrics-otlp-http',
    '@opentelemetry/exporter-logs-otlp-http',
    '@opentelemetry/exporter-prometheus',
    '@opentelemetry/resources',
    '@opentelemetry/semantic-conventions',
    '@opentelemetry/auto-instrumentations-node',
    '@opentelemetry/winston-transport',
    '@opentelemetry/resource-detector-container',
    '@opentelemetry/resource-detector-gcp',
  ],
  // Improve SSR stability - not actually needed for React 19 SSR issues
  // onDemandEntries: {
  //   maxInactiveAge: 25 * 1000,
  //   pagesBufferLength: 2,
  // },
  // Better error handling for production builds
  poweredByHeader: false,
  async headers() {
    return [
      {
        source: "/:path*",
        headers: secureHeaders,
      },
    ];
  },
};

export default withNextIntl(nextConfig);
