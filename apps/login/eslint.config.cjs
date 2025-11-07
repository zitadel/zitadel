const tsParser = require("@typescript-eslint/parser");
const typescriptEslint = require("@typescript-eslint/eslint-plugin");
const js = require("@eslint/js");

const { FlatCompat } = require("@eslint/eslintrc");

const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all
});

module.exports = [
    ...compat.extends("next", "prettier"),
    {
        languageOptions: {
            parser: tsParser,
            ecmaVersion: "latest",
            sourceType: "module",

            parserOptions: {
                ecmaFeatures: {
                    jsx: true,
                },
                project: "./tsconfig.json",
            },
        },

        plugins: {
            "@typescript-eslint": typescriptEslint,
        },

        rules: {
            "@next/next/no-html-link-for-pages": "off",
            "@next/next/no-img-element": "off",
            "react/no-unescaped-entities": "off",
            "no-unused-vars": "off",

            "@typescript-eslint/no-unused-vars": ["error", {
                argsIgnorePattern: "^_",
                varsIgnorePattern: "^_",
            }],

            "no-undef": "off",

            "no-restricted-imports": ["error", {
                paths: [{
                    name: "next/image",
                    message: "Use of next/image is forbidden. Use regular <img> elements instead.",
                }],
            }],
        },
    },
];
