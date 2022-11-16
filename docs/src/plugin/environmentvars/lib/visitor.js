var utils = require("./utils");
var test = require("./match");

module.exports = visitor;

function visitor(fence) {
  return function (node) {
    var file = this.file;
    var match;
    var found;

    match = test(node.value, fence);
    found = utils.hasData(match[1], file.data);

    if (found || found === 0) {
      return found;
    }

    return node.value;
  };
}
