import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: ["src/index.ts"],
  format: ["esm", "cjs"],
  dts: true,
  splitting: false,
  sourcemap: true,
  clean: true,
  treeshake: true,
  external: ["@angular/core", "@angular/common", "@angular/router", "@angular/common/http", "rxjs", "@zitadel/zitadel-js"],
  ...options,
}));
