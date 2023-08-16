module.exports = {
  branches: [
    { name: "next" },
    { name: "v2.31.x-spans", prerelease: "rc" }
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github"
  ],
};
