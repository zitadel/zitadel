import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: ["src/index.ts"],
  format: ["esm", "cjs"],
  dts: true,
  splitting: false,
  sourcemap: true,
  clean: true,
  treeshake: true,
  external: ["next", "react", "react-dom", "@zitadel/zitadel-js", "@zitadel/zitadel-js/v2", "@zitadel/zitadel-js/webhooks", "@zitadel/zitadel-js/node", "@zitadel/react", "oauth4webapi"],
  ...options,
}));
