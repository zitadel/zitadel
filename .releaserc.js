module.exports = {
    branches: [
        {name: 'main'},
        {name: '1.x.x', range: '1.x.x', channel: '1.x.x'}, 
        {name: 'v2', channel: 'next', prerelease: true},
      ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
