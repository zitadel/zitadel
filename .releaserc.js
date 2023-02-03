module.exports = {
    branches: [
        {name: 'main'},
        {name: '1.87.x', range: '1.87.x', channel: '1.87.x'},
        {name: '2.18.x', range: '2.18.x', channel: '2.18.x'},
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
