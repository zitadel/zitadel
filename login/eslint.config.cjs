const { fixupConfigRules } = require("@eslint/compat");
const { FlatCompat } = require("@eslint/eslintrc");

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

module.exports = [
  {
    ignores: [
      "**/.next/**",
      "**/node_modules/**",
      "**/dist/**",
      "**/build/**",
      "**/*_pb.js",
      "**/*_pb.d.ts",
      "**/packages/zitadel-proto/**",
    ],
  },
  ...fixupConfigRules(compat.extends("next/core-web-vitals", "prettier")),
  {
    languageOptions: {
      parserOptions: {
        requireConfigFile: false,
        babelOptions: {
          presets: ["next/babel"],
        },
      },
    },
    rules: {
      "@next/next/no-html-link-for-pages": "off",
    },
    settings: {
      next: {
        rootDir: ["apps/*/"],
      },
    },
  },
];
