module.exports = function isArray(x) {
  return Array.isArray ? Array.isArray(x) : x instanceof Array;
};