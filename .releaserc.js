module.exports = {
    branches: [
        { name: 'next' },
        { name: 'v2.28.x', range: '2.28.x', channel: '2.28.x' },
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
