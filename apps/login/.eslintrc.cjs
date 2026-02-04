module.exports = {
  parser: "@typescript-eslint/parser",
  extends: ["next", "prettier"],
  plugins: ["@typescript-eslint"],
  rules: {
    "@next/next/no-html-link-for-pages": "off",
    "@next/next/no-img-element": "off",
    "react/no-unescaped-entities": "off",
    "no-unused-vars": "off",
    "@typescript-eslint/no-unused-vars": ["error", {
      argsIgnorePattern: "^_" ,
      varsIgnorePattern: "^_" ,
    }],
    "no-undef": "off",
    "no-restricted-imports": ["error", {
      "paths": [{
        "name": "next/image",
        "message": "Use of next/image is forbidden. Use regular <img> elements instead."
      }]
    }],
  },
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
    ecmaFeatures: {
      jsx: true,
    },
    project: "./tsconfig.json",
  },
  overrides: [
    {
      // Server components and server actions: must use logger, not console.*
      // Files in src/lib/server/*, src/app/**/route.ts, src/middleware.ts, and files with "use server"
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
    {
      // instrumentation.ts uses console.log during OTEL initialization (before logger is available)
      files: ["src/instrumentation.ts"],
      rules: {
        "no-console": "off",
      },
    },
    {
      // Client components: logger import is forbidden (enforced by server-only at build time)
      // This override documents the expectation for client components
      files: ["src/components/**/*.tsx"],
      rules: {
        "no-restricted-imports": ["error", {
          "paths": [
            {
              "name": "next/image",
              "message": "Use of next/image is forbidden. Use regular <img> elements instead."
            },
            {
              "name": "@/lib/logger",
              "message": "Cannot import logger in client components. Use console.* instead. The logger uses 'server-only' and will fail at build time."
            }
          ]
        }],
      },
    },
    {
      // Test files: allow console.*
      files: ["**/*.test.ts", "**/*.test.tsx", "**/*.spec.ts", "**/*.spec.tsx", "**/tests/**/*"],
      rules: {
        "no-console": "off",
      },
    },
  ],
};
