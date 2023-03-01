module.exports = {
    branches: [
        {name: 'main'},
        {name: '2.21-rc', prerelease: 'rc'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
