import js from "@eslint/js";
import nextPlugin from "@next/eslint-plugin-next";
import tsPlugin from "@typescript-eslint/eslint-plugin";
import tsParser from "@typescript-eslint/parser";
import prettierConfig from "eslint-config-prettier/flat";
import importPlugin from "eslint-plugin-import";
import reactHooksPlugin from "eslint-plugin-react-hooks";
import globals from "globals";

export default [
  {
    ignores: [
      "node_modules/**",
      ".next/**",
      "dockerized/**",
      "cypress/**",
      "acceptance/**",
      "integration/**",
      "vitest.config*.ts",
      "next-env.d.ts",
    ],
  },
  js.configs.recommended,
  nextPlugin.configs["core-web-vitals"],
  {
    files: ["**/*.{ts,tsx,jsx}"],
    plugins: {
      "react-hooks": reactHooksPlugin,
    },
    rules: {
      "react-hooks/rules-of-hooks": "error",
      "react-hooks/exhaustive-deps": "error",
    },
  },
  {
    files: ["**/*.{js,mjs,cjs}"],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
  },
  {
    files: ["**/*.{ts,tsx}"],
    plugins: {
      "@typescript-eslint": tsPlugin,
      import: importPlugin,
    },
    languageOptions: {
      parser: tsParser,
      ecmaVersion: "latest",
      sourceType: "module",
      parserOptions: {
        ecmaFeatures: { jsx: true },
        project: "./tsconfig.json",
      },
    },
    settings: {
      "import/resolver": {
        typescript: true,
      },
    },
    rules: {
      ...tsPlugin.configs.recommended.rules,
      ...importPlugin.configs.recommended.rules,
      "no-unused-vars": "off",
      "@typescript-eslint/no-unused-vars": [
        "error",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
        },
      ],
      "no-undef": "off",
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/ban-ts-comment": "off",
      "@next/next/no-img-element": "off",
      "no-restricted-imports": [
        "error",
        {
          paths: [
            {
              name: "next/image",
              message:
                "Use of next/image is forbidden. Use regular <img> elements instead.",
            },
          ],
        },
      ],
    },
  },
  // Server components and server actions: must use logger, not console.*
  {
    files: [
      "src/lib/server/**/*.ts",
      "src/app/**/route.ts",
      "src/middleware.ts",
      "src/lib/cookies.ts",
    ],
    rules: {
      "no-console": "warn",
    },
  },
  // instrumentation.ts uses console.log during OTEL initialization (before logger is available)
  {
    files: ["src/instrumentation.ts"],
    rules: {
      "no-console": "off",
    },
  },
  // Client components: logger import is forbidden (enforced by server-only at build time)
  {
    files: ["src/components/**/*.tsx"],
    rules: {
      "no-restricted-imports": [
        "error",
        {
          paths: [
            {
              name: "next/image",
              message:
                "Use of next/image is forbidden. Use regular <img> elements instead.",
            },
            {
              name: "@/lib/logger",
              message:
                "Cannot import logger in client components. Use console.* instead. The logger uses 'server-only' and will fail at build time.",
            },
          ],
        },
      ],
    },
  },
  // Test files: allow console.*
  {
    files: [
      "**/*.test.ts",
      "**/*.test.tsx",
      "**/*.spec.ts",
      "**/*.spec.tsx",
      "**/tests/**/*",
    ],
    rules: {
      "no-console": "off",
    },
  },
  prettierConfig,
];
