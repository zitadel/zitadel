module.exports = {
  branches: [
    { name: "next" },
    { name: "next-rc", prerelease: "rc" },
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    [
      "@semantic-release/github",
      {
        draftRelease: true,
        successComment: false,
        assets: [
          {
            path: ".artifacts/zitadel-linux-amd64/zitadel-linux-amd64.tar.gz",
            label: "zitadel-linux-amd64.tar.gz",
          },
          {
            path: ".artifacts/zitadel-linux-arm64/zitadel-linux-arm64.tar.gz",
            label: "zitadel-linux-arm64.tar.gz",
          },
          {
            path: ".artifacts/zitadel-windows-amd64/zitadel-windows-amd64.tar.gz",
            label: "zitadel-windows-amd64.tar.gz",
          },
          {
            path: ".artifacts/zitadel-windows-arm64/zitadel-windows-arm64.tar.gz",
            label: "zitadel-windows-arm64.tar.gz",
          },
          {
            path: ".artifacts/zitadel-darwin-amd64/zitadel-darwin-amd64.tar.gz",
            label: "zitadel-darwin-amd64.tar.gz",
          },
          {
            path: ".artifacts/zitadel-darwin-arm64/zitadel-darwin-arm64.tar.gz",
            label: "zitadel-darwin-arm64.tar.gz",
          },
          {
            path: ".artifacts/checksums.txt",
            label: "checksums.txt",
          }
        ],
      },
    ],
  ],
};
