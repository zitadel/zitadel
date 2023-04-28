module.exports = {
    branches: [
        {name: 'main', channel: 'next'},
        {name: '1.87.x', range: '1.87.x', channel: '1.87.x'},
        {name: 'next', prerelease: true},
        {name: 'eventstore-phase-3', prerelease: 'alpha'},
        {name: 'custom-retry', prerelease: true}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
