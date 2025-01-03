module.exports = {
  root: true,
  // This tells ESLint to load the config from the package `eslint-config-zitadel`
  extends: ["zitadel"],
  settings: {
    next: {
      rootDir: ["apps/*/"],
    },
  },
};
