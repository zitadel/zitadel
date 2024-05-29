import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  treeshake: false,
  splitting: true,
  entry: ["src/index.ts"],
  format: ["esm", "cjs"],
  dts: true,
  minify: false,
  clean: true,
  sourcemap: true,
  ...options,
}));
