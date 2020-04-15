module.exports = {
    branches: ["master", "docker-semrel"],
    verifyConditions: [
        ["semantic-release-docker", {
            "registryUrl": "docker.pkg.github.com",
            }
        ]
    ],
    plugins: [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        "@semantic-release/github",
        ["semantic-release-docker", {
            "registryUrl": "docker.pkg.github.com",
            "name": "caos/zitadel/zitadel"
            }
        ],
    ]
};