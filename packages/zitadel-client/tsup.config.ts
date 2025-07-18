import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: ["src/index.ts", "src/v1.ts", "src/v2.ts", "src/v3alpha.ts", "src/node.ts", "src/web.ts"],
  format: ["esm", "cjs"],
  treeshake: false,
  splitting: true,
  dts: true,
  minify: false,
  clean: true,
  sourcemap: true,
  ...options,
}));
