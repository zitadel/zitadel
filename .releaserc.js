module.exports = {
  branches: [
    { name: "main" }, 
    { name: "next" },
    { name: "ci/improve-make", channel: "ignore-me", prerelease: true }
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github"
  ],
};
