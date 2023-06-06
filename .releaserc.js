module.exports = {
    branches: [
        { name: 'main' },
        { name: 'next' },
        { name: 'urlsafebase64', prerelease: true },
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
