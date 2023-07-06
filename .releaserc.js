module.exports = {
  branches: [{ name: "main" }, { name: "next" }, {name: "usage-telemetry", prerelease: true}],
  plugins: ["@semantic-release/commit-analyzer"],
};
