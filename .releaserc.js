module.exports = {
  branches: [
    { name: "main" },
    { name: "next" },
    { name: "rc", prerelease: true },
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github"
  ],
};
