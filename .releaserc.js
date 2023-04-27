module.exports = {
    branches: [
        {name: 'main'},
        {name: 'next'},
        {name: 'v2.25.x', channel: 'next'}
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
