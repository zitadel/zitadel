module.exports = {
  branches: [
    { name: "next" },
    { name: "next-rc", prerelease: "rc" },
    { name: "debug-artifacts-upload", prerelease: "ignore-me-ffo" }
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github",
    {
      "draftRelease": true,
      "successComment": false,
      "failComment": false,
      "labels": false,
      "releasedLabels": false,
      "addReleases": false,
      "failTitle": false,
    }
  ]
};
