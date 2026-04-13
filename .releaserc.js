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
            path: ".artifacts/pack/zitadel-linux-amd64.tar.gz",
            label: "zitadel-linux-amd64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-linux-arm64.tar.gz",
            label: "zitadel-linux-arm64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-windows-amd64.tar.gz",
            label: "zitadel-windows-amd64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-windows-arm64.tar.gz",
            label: "zitadel-windows-arm64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-darwin-amd64.tar.gz",
            label: "zitadel-darwin-amd64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-darwin-arm64.tar.gz",
            label: "zitadel-darwin-arm64.tar.gz",
          },
          {
            path: ".artifacts/pack/zitadel-login.tar.gz",
            label: "zitadel-login.tar.gz",
          },
          {
            path: ".artifacts/pack/checksums.txt",
            label: "checksums.txt",
          }
        ],
      },
    ],
  ],
};
