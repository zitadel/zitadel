module.exports = {
    branches: [
        {name: 'main'},
        {name: 'next'},
        {name: 'v2.27.x', range: '2.27.x', channel: '2.27.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
