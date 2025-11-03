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
};
