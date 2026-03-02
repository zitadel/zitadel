import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: {
    index: "src/index.ts",
    "auth/oidc": "src/auth/oidc.ts",
    "auth/session": "src/auth/session.ts",
    "api/bearer-token": "src/api/bearer-token.ts",
    token: "src/token.ts",
    webhooks: "src/webhooks.ts",
    v2: "src/v2.ts",
  },
  format: ["esm", "cjs"],
  dts: true,
  splitting: false,
  sourcemap: true,
  clean: !options.watch,
  treeshake: true,
  ...options,
}));
