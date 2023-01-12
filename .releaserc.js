module.exports = {
    branches: [
        {name: 'main'},
        {name: '1.87.x', range: '1.87.x', channel: '1.87.x'},
        {name: 'system-listiammembers', prerelease: 'iammembers'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
