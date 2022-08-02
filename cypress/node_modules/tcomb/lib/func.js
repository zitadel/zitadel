var assert = require('./assert');
var isTypeName = require('./isTypeName');
var FunctionType = require('./Function');
var isArray = require('./isArray');
var list = require('./list');
var isObject = require('./isObject');
var create = require('./create');
var isNil = require('./isNil');
var isBoolean = require('./isBoolean');
var tuple = require('./tuple');
var getFunctionName = require('./getFunctionName');
var getTypeName = require('./getTypeName');
var isType = require('./isType');

function getDefaultName(domain, codomain) {
  return '(' + domain.map(getTypeName).join(', ') + ') => ' + getTypeName(codomain);
}

function isInstrumented(f) {
  return FunctionType.is(f) && isObject(f.instrumentation);
}

function getOptionalArgumentsIndex(types) {
  var end = types.length;
  var areAllMaybes = false;
  for (var i = end - 1; i >= 0; i--) {
    var type = types[i];
    if (!isType(type) || type.meta.kind !== 'maybe') {
      return (i + 1);
    } else {
      areAllMaybes = true;
    }
  }
  return areAllMaybes ? 0 : end;
}

function func(domain, codomain, name) {

  domain = isArray(domain) ? domain : [domain]; // handle handy syntax for unary functions

  if (process.env.NODE_ENV !== 'production') {
    assert(list(FunctionType).is(domain), function () { return 'Invalid argument domain ' + assert.stringify(domain) + ' supplied to func(domain, codomain, [name]) combinator (expected an array of types)'; });
    assert(FunctionType.is(codomain), function () { return 'Invalid argument codomain ' + assert.stringify(codomain) + ' supplied to func(domain, codomain, [name]) combinator (expected a type)'; });
    assert(isTypeName(name), function () { return 'Invalid argument name ' + assert.stringify(name) + ' supplied to func(domain, codomain, [name]) combinator (expected a string)'; });
  }

  var displayName = name || getDefaultName(domain, codomain);
  var domainLength = domain.length;
  var optionalArgumentsIndex = getOptionalArgumentsIndex(domain);

  function FuncType(value, path) {

    if (!isInstrumented(value)) { // automatically instrument the function
      return FuncType.of(value);
    }

    if (process.env.NODE_ENV !== 'production') {
      path = path || [displayName];
      assert(FuncType.is(value), function () { return 'Invalid value ' + assert.stringify(value) + ' supplied to ' + path.join('/'); });
    }

    return value;
  }

  FuncType.meta = {
    kind: 'func',
    domain: domain,
    codomain: codomain,
    name: name,
    identity: true
  };

  FuncType.displayName = displayName;

  FuncType.is = function (x) {
    return isInstrumented(x) &&
      x.instrumentation.domain.length === domainLength &&
      x.instrumentation.domain.every(function (type, i) {
        return type === domain[i];
      }) &&
      x.instrumentation.codomain === codomain;
  };

  FuncType.of = function (f, curried) {

    if (process.env.NODE_ENV !== 'production') {
      assert(FunctionType.is(f), function () { return 'Invalid argument f supplied to func.of ' + displayName + ' (expected a function)'; });
      assert(isNil(curried) || isBoolean(curried), function () { return 'Invalid argument curried ' + assert.stringify(curried) + ' supplied to func.of ' + displayName + ' (expected a boolean)'; });
    }

    if (FuncType.is(f)) { // makes FuncType.of idempotent
      return f;
    }

    function fn() {
      var args = Array.prototype.slice.call(arguments);
      var argsLength = args.length;

      if (process.env.NODE_ENV !== 'production') {
        // type-check arguments
        var tupleLength = curried ? argsLength : Math.max(argsLength, optionalArgumentsIndex);
        tuple(domain.slice(0, tupleLength), 'arguments of function ' + displayName)(args);
      }

      if (curried && argsLength < domainLength) {
        if (process.env.NODE_ENV !== 'production') {
          assert(argsLength > 0, 'Invalid arguments.length = 0 for curried function ' + displayName);
        }
        var g = Function.prototype.bind.apply(f, [this].concat(args));
        var newDomain = func(domain.slice(argsLength), codomain);
        return newDomain.of(g, true);
      }
      else {
        return create(codomain, f.apply(this, args));
      }
    }

    fn.instrumentation = {
      domain: domain,
      codomain: codomain,
      f: f
    };

    fn.displayName = getFunctionName(f);

    return fn;

  };

  return FuncType;

}

func.getDefaultName = getDefaultName;
func.getOptionalArgumentsIndex = getOptionalArgumentsIndex;
module.exports = func;
