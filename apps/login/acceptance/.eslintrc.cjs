module.exports = {
  root: true,
  // Use basic ESLint config since the login app has its own detailed config
  extends: ["eslint:recommended"],
  settings: {
    next: {
      rootDir: ["apps/*/"],
    },
  },
};
