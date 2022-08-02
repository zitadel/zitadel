var isType = require('./isType');
var getFunctionName = require('./getFunctionName');

module.exports = function getTypeName(ctor) {
  if (isType(ctor)) {
    return ctor.displayName;
  }
  return getFunctionName(ctor);
};