module.exports = {
    branches: [
        { name: 'main' },
        { name: 'next' },
        { name: 'friendly-quota-depleted-screen-acceptance', prerelease: true },
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
