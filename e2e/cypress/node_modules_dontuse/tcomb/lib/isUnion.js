var isType = require('./isType');

module.exports = function isUnion(x) {
  return isType(x) && ( x.meta.kind === 'union' );
};