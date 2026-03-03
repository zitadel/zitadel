import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: {
    index: "src/index.ts",
    "auth/oidc": "src/auth/oidc.ts",
    "auth/session": "src/auth/session.ts",
    "auth/bearer-token": "src/api/bearer-token.ts",
    "api/v1": "src/api/v1.ts",
    "api/v2": "src/api/v2.ts",
    token: "src/token.ts",
    "actions/webhook": "src/webhooks.ts",
  },
  format: ["esm", "cjs"],
  dts: true,
  splitting: false,
  sourcemap: true,
  clean: !options.watch,
  treeshake: true,
  ...options,
}));
