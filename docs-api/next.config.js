const withMarkdoc = require("@markdoc/next.js");

module.exports = withMarkdoc({
  basePath: "/docs/api",
  target: "serverless",
})({
  pageExtensions: ["js", "jsx", "ts", "tsx", "md", "mdoc"],
});
