module.exports = {
    branches: [
        {name: 'main'},
        {name: 'next'},
        {name: 'v2.22.x', range: '2.22.x', channel: '2.22.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
