const withMarkdoc = require("@markdoc/next.js");

module.exports = withMarkdoc({ mode: "server" })({
  basePath: "/docs/api",
  pageExtensions: ["js", "jsx", "ts", "tsx", "md", "mdoc"],
});
