/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: '/openapi/:path*',
        destination: '/api/openapi/:path*',
      },
    ];
  },
};

export default nextConfig;
