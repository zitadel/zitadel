import tsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [tsconfigPaths()],
  test: {
    include: ["dockerized/**/*.test.ts"],
    testTimeout: 180000,
    hookTimeout: 180000,
    fileParallelism: false,
  },
});
