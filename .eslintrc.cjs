module.exports = {
  root: true,
  // This tells ESLint to load the config from the package `@zitadel/eslint-config`
  extends: ["zitadel"],
  settings: {
    next: {
      rootDir: ["apps/*/"],
    },
  },
};
