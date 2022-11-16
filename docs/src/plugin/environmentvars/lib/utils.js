exports.settings = settings;
exports.hasData = hasData;

function settings(options) {
  var returns;

  returns = {};
  options = options || {};
  returns.fence = fences(options);
  returns.fail = is(options.fail, "Boolean") ? options.fail : false;
  returns.quiet = is(options.quiet, "Boolean") ? options.quiet : false;
  returns.name = is(options.name, "String") ? options.name : "variables";
  return returns;
}

function fences(value) {
  var defaults = ["{{", "}}"];

  if (value && value.fence) {
    return fences(value.fence);
  }

  if (is(value, "Array") && value.length) {
    var open = is(value[0], "String") ? value[0] : defaults[0];
    var closed;

    if (open === defaults[0]) {
      closed = defaults[1];
    } else {
      closed = is(value[1], "String") ? value[1] : value[0];
    }

    return [open, closed];
  }
  if (is(value, "String")) {
    return [value, value];
  }
  if (is(value, "Object") && value.open) {
    return [value.open, value.close || value.open];
  }

  return defaults;
}

function is(value, expected) {
  var toString = Object.prototype.toString;
  var match = toString.call(value).match(/\[object(.*?)\]/);

  /* istanbul ignore else */
  if (match) {
    actual = match[1].trim();
    return actual === expected;
  } else {
    return false;
  }
}

function hasData(string, data) {
  var splitter = /\.|\[(\d+)\]/;
  var value;
  var keys;
  var i;

  function empty(value) {
    return value && value.length;
  }

  i = -1;
  value = data;
  string = string.trim();
  keys = string.split(splitter).filter(empty);

  while (++i < keys.length) {
    val = value[keys[i]];
    if (val != null) {
      value = val;
      continue;
    } else {
      return null;
    }
  }

  return value;
}
