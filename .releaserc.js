module.exports = {
    branches: [
        {name: 'main'},
        {name: '2.8.x', range: '2.8.x', channel: '2.8.x'},
        {name: '1.87.x', range: '1.87.x', channel: '1.87.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
