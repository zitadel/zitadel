module.exports = {
    branches: ["master", "docker-semrel"],
    plugins: [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        "@semantic-release/github",
        ["@semantic-release/exec", {
            "publishCmd": "echo '::set-env name=CAOS_NEXT_VERSION::${nextRelease.version}'"
            }],
    ]
};