module.exports = {
    branch: 'master',
    plugins: [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        "@semantic-release/github",
        ["@semantic-release/exec", {
            "prepareCmd": "echo '::set-env name=CAOS_NEXT_VERSION::v${nextRelease.version}'"
        }],
        ["semantic-release-docker", {
            "verifyConditions": {
                "registryUrl": "docker.pkg.github.com"
            },
            "publish": {
                "name": "caos/zitadel/zitadel"
            }
        }],
    ]
};