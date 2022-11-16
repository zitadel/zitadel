module.exports = match;

function match(value, fence) {
  var offset;
  var start;
  var close;
  var open;
  var same;
  var left;
  var nl;

  if (!fence) {
    return;
  }

  open = fence[0];
  close = fence[1];
  same = open === close;
  start = value.indexOf(open);
  offset = value.indexOf(close, start);
  nl = "\n";

  if (start !== 0 || offset < 0) {
    return;
  }

  index = start + open.length - 1;
  left = 1;

  while (left > 0 && value.charAt(index) !== nl) {
    index += 1;

    if (value.slice(index, index + open.length) === open && !same) {
      left += 1;
    }

    if (value.slice(index, index + close.length) !== close) {
      continue;
    }

    left -= 1;

    if (left < 1) {
      return [
        value.slice(start, index + close.length),
        value.slice(start + open.length, index).trim(),
      ];
    }
  }
}
