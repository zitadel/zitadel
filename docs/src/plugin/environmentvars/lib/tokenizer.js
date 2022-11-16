var utils = require("./utils");
var test = require("./match");

module.exports = tokenizer;

function tokenizer(name, data, fence, quiet, fail) {
  return function (eat, value, silent) {
    var tokenized;
    var fromFile;
    var fromData;
    var subvalue;
    var message;
    var found;
    var match;
    var file;
    var self;
    var node;
    var add;
    var val;
    var now;
    var i;

    self = this;
    file = self.file;
    match = test(value, fence);

    if (match) {
      subvalue = match[0];
      sub = match[1].trim();

      fromFile = utils.hasData(sub, file.data);
      fromData = utils.hasData(sub, data);
      found = fromFile || fromData;
      now = eat.now();

      /* istanbul ignore if */
      if (silent) {
        return true;
      }

      add = eat(subvalue);

      if (found != null) {
        tokenized = self.tokenizeInline(found.toString(), now);
        i = -1;

        while (++i < tokenized.length) {
          node = add(tokenized[i]);
        }

        return node;
      }

      if (!found && !quiet) {
        sub = sub.indexOf(".") === 0 ? sub : "." + sub;
        message = "Could not resolve `data" + sub + "` in VFile or Processor.";

        if (fail) {
          return file.fail(message, now, "variables:undef-variable");
        } else {
          file.message(message, now, "variables:undef-variable");
        }
      }

      return add({
        type: name,
        value: subvalue,
      });
    }
  };
}
