module.exports = {
  branches: [
    { name: "next" },
    { name: "next-rc", prerelease: "rc" },
    { name: 'merge-eventstore', prerelease: 'eventstore-performance'}
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github"
  ],
};
