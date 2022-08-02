var assert = require('./assert');
var isFunction = require('./isFunction');
var isNil = require('./isNil');
var getTypeName = require('./getTypeName');
var isObject = require('./isObject');
var isArray = require('./isArray');
var isType = require('./isType');
var create = require('./create');
var assign = require('./assign');

function assignMany(values) {
  return values.reduce(function (acc, v) {
    return assign(acc, v);
  }, {});
}

function fromJSON(value, type, path) {
  if (process.env.NODE_ENV !== 'production') {
    assert(isFunction(type), function () {
      return 'Invalid argument type ' + assert.stringify(type) + ' supplied to fromJSON(value, type) (expected a type)';
    });
    path = path || [getTypeName(type)];
  }

  if (isFunction(type.fromJSON)) {
    return create(type, type.fromJSON(value), path);
  }

  if (!isType(type)) {
    return value instanceof type ? value : new type(value);
  }

  var kind = type.meta.kind;
  var k;
  var ret;

  switch (kind) {

    case 'maybe' :
      return isNil(value) ? null : fromJSON(value, type.meta.type, path);

    case 'subtype' : // the kind of a refinement is 'subtype' (for legacy reasons)
      ret = fromJSON(value, type.meta.type, path);
      if (process.env.NODE_ENV !== 'production') {
        assert(type.meta.predicate(ret), function () {
          return 'Invalid argument value ' + assert.stringify(value) + ' supplied to fromJSON(value, type) (expected a valid ' + getTypeName(type) + ')';
        });
      }
      return ret;

    case 'struct' :
      if (process.env.NODE_ENV !== 'production') {
        assert(isObject(value), function () {
          return 'Invalid argument value ' + assert.stringify(value) + ' supplied to fromJSON(value, type) (expected an object for type ' + getTypeName(type) + ')';
        });
      }
      var props = type.meta.props;
      var defaultProps = type.meta.defaultProps;
      ret = {};
      for (k in props) {
        var actual = value[k];
        if (actual === undefined) {
          actual = defaultProps[k];
        }
        if (props.hasOwnProperty(k)) {
          ret[k] = fromJSON(actual, props[k], ( process.env.NODE_ENV !== 'production' ? path.concat(k + ': ' + getTypeName(props[k])) : null ));
        }
      }
      return new type(ret);

    case 'interface' :
      if (process.env.NODE_ENV !== 'production') {
        assert(isObject(value), function () {
          return 'Invalid argument value ' + assert.stringify(value) + ' supplied to fromJSON(value, type) (expected an object)';
        });
      }
      var interProps = type.meta.props;
      ret = {};
      for (k in interProps) {
        if (interProps.hasOwnProperty(k)) {
          ret[k] = fromJSON(value[k], interProps[k], ( process.env.NODE_ENV !== 'production' ? path.concat(k + ': ' + getTypeName(interProps[k])) : null ));
        }
      }
      return ret;

    case 'list' :
      if (process.env.NODE_ENV !== 'production') {
        assert(isArray(value), function () {
          return 'Invalid argument value ' + assert.stringify(value) + ' supplied to fromJSON(value, type) (expected an array for type ' + getTypeName(type) + ')';
        });
      }
      var elementType = type.meta.type;
      var elementTypeName = getTypeName(elementType);
      return value.map(function (element, i) {
        return fromJSON(element, elementType, ( process.env.NODE_ENV !== 'production' ? path.concat(i + ': ' + elementTypeName) : null ));
      });

    case 'union' :
      var actualType = type.dispatch(value);
      if (process.env.NODE_ENV !== 'production') {
        assert(isFunction(actualType), function () {
          return 'Invalid argument value ' + assert.stringify(value) + ' supplied to fromJSON(value, type) (no constructor returned by dispatch of union ' + getTypeName(type) + ')';
        });
      }
      return fromJSON(value, actualType, path);

    case 'tuple' :
      if (process.env.NODE_ENV !== 'production') {
        assert(isArray(value), function () {
          return 'Invalid argument value ' + assert.stringify(value) + ' supplied to fromJSON(value, type) (expected an array for type ' + getTypeName(type) + ')';
        });
      }
      var types = type.meta.types;
      if (process.env.NODE_ENV !== 'production') {
        assert(isArray(value) && value.length === types.length, function () {
          return 'Invalid value ' + assert.stringify(value) + ' supplied to fromJSON(value, type) (expected an array of length ' + types.length + ' for type ' + getTypeName(type) + ')';
        });
      }
      return value.map(function (element, i) {
        return fromJSON(element, types[i], ( process.env.NODE_ENV !== 'production' ? path.concat(i + ': ' + getTypeName(types[i])) : null ));
      });

    case 'dict' :
      if (process.env.NODE_ENV !== 'production') {
        assert(isObject(value), function () {
          return 'Invalid argument value ' + assert.stringify(value) + ' supplied to fromJSON(value, type) (expected an object for type ' + getTypeName(type) + ')';
        });
      }
      var domain = type.meta.domain;
      var codomain = type.meta.codomain;
      var domainName = getTypeName(domain);
      var codomainName = getTypeName(codomain);
      ret = {};
      for (k in value) {
        if (value.hasOwnProperty(k)) {
          ret[domain(k, ( process.env.NODE_ENV !== 'production' ? path.concat(domainName) : null ))] = fromJSON(value[k], codomain, ( process.env.NODE_ENV !== 'production' ? path.concat(k + ': ' + codomainName) : null ));
        }
      }
      return ret;

    case 'intersection' :
      var values = type.meta.types.map(function (type, i) {
        return fromJSON(value, type, ( process.env.NODE_ENV !== 'production' ? path.concat(i + ': ' + getTypeName(type)) : null ));
      });
      return type(
        isObject(values[0]) ? assignMany(values) : value,
        path
      );

    default : // enums, irreducible
      return type(value, path);
  }
}

module.exports = fromJSON;
