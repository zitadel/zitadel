module.exports = {
    branches: [
        {name: 'main'},
        {name: '1.87.x', range: '1.87.x', channel: '1.87.x'},
        {name: '2.16.x', range: '2.16.x', channel: '2.16.x'},
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
