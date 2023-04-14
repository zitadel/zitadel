module.exports = {
    branches: [
        {name: 'main'},
        {name: 'fix(eventstore)--search', prerelease: 'correct-creation-date'},
        {name: '1.87.x', range: '1.87.x', channel: '1.87.x'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
