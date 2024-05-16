import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: ["src/index.tsx"],
  format: ["esm", "cjs"],
  dts: true,
  external: ["react"],
  ...options,
}));
