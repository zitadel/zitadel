import tsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [tsconfigPaths()],
  test: {
    include: ["dockerized/**/*.test.ts"],
    // Increased timeouts are intentional: Docker image build and container startup
    // can be slow in CI and local environments, so we allow up to 3 minutes.
    testTimeout: 180000,
    hookTimeout: 180000,
    fileParallelism: false,
  },
});
