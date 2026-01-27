import react from "@vitejs/plugin-react";
import tsconfigPaths from "vite-tsconfig-paths";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [tsconfigPaths(), react()],
  resolve: {
    alias: {
      // Mock server-only package in tests (it throws an error outside of RSC context)
      "server-only": new URL("./test-mocks/server-only.ts", import.meta.url).pathname,
    },
  },
  test: {
    include: ["src/**/*.test.ts", "src/**/*.test.tsx", "tests/**/*.test.ts"],
    environment: "jsdom",
    setupFiles: ["./test-setup.ts"],
  },
});
