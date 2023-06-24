module.exports = {
    branches: [
        { name: 'next' },
        { name: 'v2.29.x', range: '2.29.x', channel: '2.29.x' },
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
