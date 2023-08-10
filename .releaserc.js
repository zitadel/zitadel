module.exports = {
  branches: [
    { name: "next" },
    { name: "next-rc", prerelease: "rc" },
    { name: "debug-artifacts-upload", prerelease: "ignore-me" }
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github",
    {
      "draftRelease": true
    }
  ]
};
