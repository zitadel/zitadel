module.exports = function getFunctionName(f) {
  return f.displayName || f.name || '<function' + f.length + '>';
};