module.exports = {
    branches: [
        {name: 'main'},
        {name: '1.87.x', range: '1.87.x', channel: '1.87.x'},
        {name: 'v2.19.x', range: '2.19.x', channel: '2.19.x'},
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
