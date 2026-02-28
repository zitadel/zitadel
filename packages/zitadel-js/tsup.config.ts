import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: { index: "src/index.ts", token: "src/token.ts", webhooks: "src/webhooks.ts", v2: "src/v2.ts" },
  format: ["esm", "cjs"],
  dts: true,
  splitting: false,
  sourcemap: true,
  clean: !options.watch,
  treeshake: true,
  ...options,
}));
