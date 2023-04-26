module.exports = {
    branches: [
        {name: 'main'},
        {name: 'next'},
        {name: 'optimise-step-10', prerelease: true}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
