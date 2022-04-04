module.exports = {
    branches: [
        {name: 'main'},
        {name: '1.x.x', range: '1.x.x', channel: '1.x.x'},
        {name: 'workflow-dispatch', prerelease: true},
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
