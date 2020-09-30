module.exports = {
    branches: ["master"],
    plugins: [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        ["@semantic-release/github", {
            "assets": [
              {"path": ".artifacts/zitadel-darwin-amd64/zitadel-darwin-amd64", "label": "Zitadel Darwin x86_64"},
              {"path": ".artifacts/zitadel-linux-amd64/zitadel-linux-amd64", "label": "Zitadel Linux x86_64"},
              {"path": ".artifacts/zitadel-windows-amd64/zitadel-windows-amd64", "label": "Zitadel Windows x86_64"},
              {"path": ".artifacts/zitadel-darwin-amd64/zitadelctl-darwin-amd64", "label": "Zitadelctl Darwin x86_64"},
              {"path": ".artifacts/zitadel-linux-amd64/zitadelctl-linux-amd64", "label": "Zitadelctl Linux x86_64"},
              {"path": ".artifacts/zitadel-windows-amd64/zitadelctl-windows-amd64", "label": "Zitadelctl Windows x86_64"}
            ]
        }],
    ]
};