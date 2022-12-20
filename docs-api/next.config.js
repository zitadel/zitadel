const withMarkdoc = require("@markdoc/next.js");

module.exports = withMarkdoc({
  mode: "server",
  nodes: {
    variables: {
      protocol: "rest",
      language: "js",
    },
  },
})({
  basePath: "/docs/api",
  pageExtensions: ["js", "jsx", "ts", "tsx", "md", "mdoc"],
});
