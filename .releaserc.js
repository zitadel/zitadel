module.exports = {
    branches: [
        {name: 'main', channel: 'next'},
        {name: 'next', prerelease: true}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
