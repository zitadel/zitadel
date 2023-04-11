module.exports = {
    branches: [
        {name: 'main'},
        {name: "next", prerelease: true},
        {name: '1.87.x', range: '1.87.x', channel: '1.87.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
