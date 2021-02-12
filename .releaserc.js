module.exports = {
    branches: ["master"],
    plugins: [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        ["@semantic-release/github", {
            "assets": [
                {
                    "path": "./artifacts/zitadelctl-darwin-amd64/zitadelctl-darwin-amd64",
                    "label": "Zitadelctl Darwin x86_64"
                },
                {
                    "path": "./artifacts/zitadelctl-linux-amd64/zitadelctl-linux-amd64",
                    "label": "Zitadelctl Linux x86_64"
                },
                {
                    "path": "./artifacts/zitadelctl-windows-amd64/zitadelctl-windows-amd64.exe",
                    "label": "Zitadelctl Windows x86_64"
                }
            ]
        }],
    ]
};
