var isType = require('./isType');

module.exports = function isStruct(x) {
  return isType(x) && ( x.meta.kind === 'struct' );
};