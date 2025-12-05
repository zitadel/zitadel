const next = require("eslint-config-next");
const prettier = require("eslint-config-prettier");

module.exports = [
    {
        ignores: [
            "**/.next/**",
            "**/node_modules/**",
            "**/playwright-report/**",
            "**/test-results/**",
            "**/dist/**",
            "**/build/**",
            "**/acceptance/**",
            "**/integration/**",
        ],
    },
    ...next,
    prettier,
    {
        rules: {
            "@next/next/no-html-link-for-pages": "off",
            "@next/next/no-img-element": "off",
            "react/no-unescaped-entities": "off",
            "no-unused-vars": "off",
            "no-undef": "off",

            "no-restricted-imports": ["error", {
                paths: [{
                    name: "next/image",
                    message: "Use of next/image is forbidden. Use regular <img> elements instead.",
                }],
            }],
        },
    },
    {
        files: ["**/*.ts", "**/*.tsx"],
        rules: {
            "@typescript-eslint/no-unused-vars": ["error", {
                argsIgnorePattern: "^_",
                varsIgnorePattern: "^_",
            }],
        },
    }
];
