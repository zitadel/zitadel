import createNextIntlPlugin from "next-intl/plugin";
import { DEFAULT_CSP } from "./constants/csp.js";
import path from "path";
import { fileURLToPath } from "url";

// Fix for ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const withNextIntl = createNextIntlPlugin();

/** @type {import('next').NextConfig} */

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
    key: "X-Frame-Options",
    value: "SAMEORIGIN",
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
    value: `${DEFAULT_CSP} frame-ancestors 'none'`,
  },
  { key: "X-Frame-Options", value: "deny" },
];

const nextConfig = {
  basePath: "/ui/v2/login",
  output: "standalone",
  outputFileTracingRoot: path.join(__dirname, "../../"),
  reactStrictMode: true,
  experimental: {
    dynamicIO: true,
    // Add React 19 compatibility optimizations
    optimizePackageImports: ['@radix-ui/react-tooltip', '@heroicons/react'],
  },
  images: {
    unoptimized: true
  },
  eslint: {
    ignoreDuringBuilds: true,
  },
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
