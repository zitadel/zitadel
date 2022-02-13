module.exports = {
    branches: [
        {name: 'main'},
        {name: 'v2-ci', channel: 'next'},
      ]
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
