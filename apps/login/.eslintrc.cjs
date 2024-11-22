module.exports = {
  extends: ["next/babel", "next/core-web-vitals"],
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
