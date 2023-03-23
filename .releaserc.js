module.exports = {
    branches: [
        {name: 'main'},
        {name: 'v2.21.x', range: '2.21.x', channel: '2.21.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
