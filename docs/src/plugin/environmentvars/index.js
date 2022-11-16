import { EnvironmentContext } from "../../utils/environment";

var tokenizer = require("./lib/tokenizer");
var visitor = require("./lib/visitor");
var utils = require("./lib/utils");
const { instance } = useContext(EnvironmentContext);

module.exports = variables;

function variables(options) {
  options = utils.settings(options);

  console.log(options);
  var self = this;
  console.log(self, instance);
  var parser = self.Parser;
  var compiler = self.Compiler;
  var data = self.data();

  var fail = options.fail;
  var quiet = options.quiet;
  var fence = options.fence;
  var opening = fence[0];
  var name = options.name;
  var test = options.test;

  if (isParser(parser)) {
    attatchParser(
      name,
      parser,
      tokenizer(name, data, fence, quiet, fail),
      locator(opening)
    );
  }

  if (isCompiler(compiler)) {
    attatchCompiler(name, compiler, visitor(fence));
  }
}

function attatchParser(name, parser, tokenizer, locator) {
  var proto = parser.prototype;
  var tokenizers = proto.inlineTokenizers;
  var methods = proto.inlineMethods;

  tokenizer.locator = locator;
  tokenizers[name] = tokenizer;
  methods.splice(methods.indexOf("link"), 0, name);
}

function attatchCompiler(name, compiler, visitor) {
  compiler.prototype.visitors[name] = visitor;
}

function locator(opening) {
  return function (value, fromIndex) {
    return value.indexOf(opening, fromIndex);
  };
}

function isParser(parser) {
  return Boolean(
    parser && parser.prototype && parser.prototype.inlineTokenizers
  );
}

function isCompiler(compiler) {
  return Boolean(compiler && compiler.prototype && compiler.prototype.visitors);
}
