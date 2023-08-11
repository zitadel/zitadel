module.exports = {
  branches: [
    { name: "next" },
    { name: "next-rc", prerelease: "rc" },
    { name: "pipeline-upload-assets", prerelease: "ignore-me2" }
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    [
      "@semantic-release/github",
      {
        assets: [
          {
            path: ".artifacts/zitadel-linux-amd64/zitadel-linux-amd64.tar",
            label: "zitadel-linux-amd64.tar",
          },
          {
            path: ".artifacts/zitadel-linux-arm64/zitadel-linux-arm64.tar",
            label: "zitadel-linux-arm64.tar",
          },
          {
            path: ".artifacts/zitadel-windows-amd64/zitadel-windows-amd64.tar",
            label: "zitadel-windows-amd64.tar",
          },
          {
            path: ".artifacts/zitadel-windows-arm64/zitadel-windows-arm64.tar",
            label: "zitadel-windows-arm64.tar",
          },
          {
            path: ".artifacts/zitadel-darwin-amd64/zitadel-darwin-amd64.tar",
            label: "zitadel-darwin-amd64.tar",
          },
          {
            path: ".artifacts/zitadel-darwin-arm64/zitadel-darwin-arm64.tar",
            label: "zitadel-darwin-arm64.tar",
          },
        ],
        draftRelease: true,
        successComment: false,
        failComment: false,
        labels: false,
        releasedLabels: false,
        addReleases: false,
        failTitle: false,
      },
    ],
  ],
};
