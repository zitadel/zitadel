import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  treeshake: true,
  splitting: true,
  publicDir: true,
  entry: ["src/index.ts", "src/**/index.ts"],
  format: ["esm", "cjs"],
  dts: true,
  minify: true,
  clean: true,
  ...options,
}));
