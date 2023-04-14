import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  treeshake: true,
  splitting: true,
  publicDir: true,
  entry: ["src/**/*.{tsx,png,svg}"],
  format: ["esm"],
  dts: true,
  minify: true,
  clean: true,
  external: ["react"],
  ...options,
}));
