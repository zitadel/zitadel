import { defineConfig } from "@eloqnt/cli";

export default defineConfig({
  srcPath: "./src",
  messages: {
    path: "./locales",
    locales: "infer",
    sourceLocale: "en",
    format: "json",
  },
  lint: {
    rules: {
      // Most messages are rendered via `<Translated i18nKey="…" namespace="…" />`
      // (src/components/translated.tsx), which the linter can't follow
      "orphan-message": "off",
    },
  },
});
