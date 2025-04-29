import { defineConfig } from "tsup";

export default defineConfig({
  entry: ["src/index.ts", "src/v1.ts", "src/v2.ts", "src/v3alpha.ts"],
  dts: true,
  clean: true,
  minify: false,
  splitting: false,
  treeshake: false,
  sourcemap: true,
  format: ["esm", "cjs"],
  platform: "neutral",
  target: "node16",
});
