const Mocha = require('mocha');

const mochaSerializeSuite = Mocha.Suite.prototype.serialize;

// Serialize the full root suite state to count `Skipped` tests.
Mocha.Suite.prototype.serialize = function (...args) {
  if (this.root) {
    // Skipping the EVENT_SUITE_BEGIN serialization can reduce data transfer via IPC by 25%
    return serializeSuite(this);
  }
  return mochaSerializeSuite.apply(this, args)
}

const serializeSuite = (suite) => {
  const result = suite.root ? mochaSerializeSuite.call(suite) : serializeObject(suite, ['file']);
  result.suites = suite.suites.map(it => serializeSuite(it));
  result.tests = suite.tests.map(it => serializeTest(it));
  ['_beforeAll', '_beforeEach', '_afterEach', '_afterAll'].forEach(hookName => {
    result[hookName] = suite[hookName].map(it => serializeHook(it))
  });
  return result;
}

const serializeHook = hook => {
  return serializeObject(hook, ['body', 'state', 'err', 'context', '$$fullTitle']);
}

const serializeTest = test => {
  const result = serializeObject(test, ['pending', 'context']);
  // to remove a circular dependency: https://github.com/adamgruber/mochawesome/issues/356
  result["$$retriedTest"] = null;
  return result;
}

const serializeObject = (obj, fields) => {
  const result = obj.serialize();
  for (let field of fields) {
    // The field's started with `$$` are results of methods
    result[field] = field.startsWith('$$') ? obj[field.slice(2)]() : obj[field];
  }
  if (result.err instanceof Error) {
    result.err = serializeError(result.err);
  }
  return result;
};

const serializeError = error => {
  /* The default properties of Error class: name, message and stack; are excluded from the enumeration.
     It causes the following: JSON.stringify(new Error("FAKE")) === '{}'
     So, we need to provide explicitly these properties to the JSON serializer. */
  if (error instanceof Error) {
    return {
      message: error.message,
      stack: error.stack,
      name: error.name,
      ...error,
    };
  }
  return error;
};

module.exports = { serializeSuite, serializeHook, serializeTest, serializeObject, serializeError };
