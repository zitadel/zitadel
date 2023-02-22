module.exports = {
    branches: [
        {name: 'main'},
        {name: 'v2.20.x', range: '2.20.x', channel: '2.20.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
