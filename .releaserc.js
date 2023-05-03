module.exports = {
    branches: [
        {name: 'next'},
        {name: 'v2.25.x', range: '2.25.x', channel: '2.25.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
