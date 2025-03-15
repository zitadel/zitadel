import { defineConfig } from "tsup";

export default defineConfig({
  entry: ["index.ts"],
  dts: true,
  clean: true,
  minify: false,
  splitting: false,
  sourcemap: true,
  format: ["esm", "cjs"],
  platform: "neutral",
  target: "node16",
});
