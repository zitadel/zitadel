module.exports = {
    branches: [
        {name: 'next'},
        {name: 'v2.27-rc', prerelease: 'rc', channel: 'next'},
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
