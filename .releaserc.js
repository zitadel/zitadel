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
            path: ".artifacts/pack/zitadel-api-linux-amd64.tar.gz",
            label: "zitadel-api-linux-amd64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-api-linux-arm64.tar.gz",
            label: "zitadel-api-linux-arm64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-api-windows-amd64.tar.gz",
            label: "zitadel-api-windows-amd64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-api-windows-arm64.tar.gz",
            label: "zitadel-api-windows-arm64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-api-darwin-amd64.tar.gz",
            label: "zitadel-api-darwin-amd64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-api-darwin-arm64.tar.gz",
            label: "zitadel-api-darwin-arm64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-login.tar.gz",
            label: "zitadel-login.tar.gz",
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
