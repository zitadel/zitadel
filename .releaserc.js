module.exports = {
  branches: [
    { name: "next" },
    { name: "next-rc", prerelease: "rc" },
    { name: "next-auth-spans", prerelease: "rc" }
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github"
  ],
};
