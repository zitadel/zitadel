import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: {
    index: "src/index.ts",
    "auth/index": "src/auth/index.ts",
    "auth/oidc": "src/auth/oidc.ts",
    middleware: "src/middleware.ts",
    "server-action": "src/server-action.ts",
    api: "src/api.ts",
    webhook: "src/webhook.ts",
  },
  format: ["esm", "cjs"],
  dts: true,
  splitting: true,
  sourcemap: true,
  clean: true,
  treeshake: true,
  external: [
    "next",
    "next/headers",
    "next/navigation",
    "next/server",
    "react",
    "react-dom",
    "@zitadel/zitadel-js",
    "@zitadel/zitadel-js/v2",
    "@zitadel/zitadel-js/webhooks",
    "@zitadel/zitadel-js/node",
    "@zitadel/react",
    "oauth4webapi",
  ],
  ...options,
}));
