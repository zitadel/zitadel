module.exports = {
    branches: [
        {name: 'main'},
        {name: 'v2.22-rc', prerelease: 'rc'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
