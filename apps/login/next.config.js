/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true, // Recommended for the `pages` directory, default in `app`.
  swcMinify: true,
  experimental: {
    serverActions: true,
  },
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "zitadel.com",
        port: "",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "zitadel.cloud",
        port: "",
        pathname: "/**",
      },
    ],
  },
};

module.exports = nextConfig;
