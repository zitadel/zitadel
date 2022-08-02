module.exports = function fail(message) {
  throw new TypeError('[tcomb] ' + message);
};