module.exports = {
    branches: [
        // branch to create latest version
        'next',
        // main is the default branch and should create a draft release for each merge of frontend/backend
        {name: 'main', channel: 'next', prerelease: 'rc'},
        // maintainance branches for v2 minor releases
        // TODO: replace ? with actual minor version
        // {name: 'v2.?.x', range: '2.?.x', channel: '2.?.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
