/** @type {import('next').NextConfig} */
const nextConfig = {
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
  // Required for @connectrpc/connect-node (gRPC)
  serverExternalPackages: ["@connectrpc/connect-node"],
  // Allow importing components from the console app
  transpilePackages: ["@zitadel/client", "@zitadel/proto", "@zitadel/react"],
  async rewrites() {
    // Console app pages use links like /users, /organizations etc.
    // (because the standalone console has basePath: "/console").
    // Rewrite these to /console/* so they resolve in the cloud app.
    const consoleRoutes = [
      "overview",
      "users",
      "organizations",
      "projects",
      "applications",
      "actions",
      "sessions",
      "administrators",
      "activity",
      "settings",
      "getting-started",
      "account-settings",
      "feedback",
      "roles",
      "analytics",
      "billing",
      "support",
      "usage",
    ];
    return consoleRoutes.flatMap((route) => [
      { source: `/${route}`, destination: `/console/${route}` },
      { source: `/${route}/:path*`, destination: `/console/${route}/:path*` },
    ]);
  },
};

export default nextConfig;
