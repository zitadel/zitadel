module.exports = {
  parser: "@babel/eslint-parser",
  extends: ["next", "prettier"],
  rules: {
    "@next/next/no-html-link-for-pages": "off",
    "@next/next/no-img-element": "off",
    "react/no-unescaped-entities": "off"
  },
  parserOptions: {
    requireConfigFile: false,
    babelOptions: {
      presets: ["next/babel"],
    },
  },
};
