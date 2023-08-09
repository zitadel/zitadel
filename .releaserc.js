module.exports = {
  branches: [
    { name: "next" },
    { name: "next-rc", prerelease: "rc" },
    { name: "pipeline-image", prerelease: "ignore-me" }
  ],
  plugins: [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    [
      "@semantic-release/github",
      {
        "assets": [
          { "path": "./zitadel-linux-amd64.tar.gz", "label": "zitadel-linux-amd64.tar.gz" },
          { "path": "./zitadel-linux-arm64.tar.gz", "label": "zitadel-linux-arm64.tar.gz" },
          { "path": "./zitadel-windows-amd64.tar.gz", "label": "zitadel-windows-amd64.tar.gz" },
          { "path": "./zitadel-windows-arm64.tar.gz", "label": "zitadel-windows-arm64.tar.gz" },
          { "path": "./zitadel-darwin-amd64.tar.gz", "label": "zitadel-darwin-amd64.tar.gz" },
          { "path": "./zitadel-darwin-arm64.tar.gz", "label": "zitadel-darwin-arm64.tar.gz" },
        ]
      }
    ]
  ]
};
