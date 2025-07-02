module.exports = {
  extends: ["next/core-web-vitals"],
  ignorePatterns: ["external/**/*.ts"],
  rules: {
    "@next/next/no-html-link-for-pages": "off",
  },
  settings: {
    react: {
      version: "detect",
    },
  },
};
