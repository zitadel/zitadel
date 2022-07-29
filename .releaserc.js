module.exports = {
    branches: [
        {name: 'v1',  range: '1.x.x', channel: '1.x.x'},
        {name: '1.x.x', range: '1.x.x', channel: '1.x.x'},
        {name: 'v2-alpha', prerelease: true},
        {name: 'v2-alpha-import', prerelease: true},
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
