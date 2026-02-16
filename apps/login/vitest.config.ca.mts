import tsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [tsconfigPaths()],
  test: {
    include: ["dockerized/ca/**/*.test.ts"],
    testTimeout: 30000,
    hookTimeout: 180000,
  },
});
