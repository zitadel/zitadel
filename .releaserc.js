module.exports = {
    branches: [
        {name: 'main', channel: 'next'},
        {name: 'next', prerelease: true},
        {name: 'eventstore-phase-3', prerelease: 'alpha'},
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
