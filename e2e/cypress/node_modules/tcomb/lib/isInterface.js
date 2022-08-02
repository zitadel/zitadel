var isType = require('./isType');

module.exports = function isInterface(x) {
  return isType(x) && ( x.meta.kind === 'interface' );
};