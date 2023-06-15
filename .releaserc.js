module.exports = {
    branches: [
        { name: 'main' },
        { name: 'next' },
        { name: 'merge-eventstore', prerelease: 'eventstore-performance'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
