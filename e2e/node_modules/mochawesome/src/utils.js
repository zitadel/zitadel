const isString = require('lodash.isstring');
const isFunction = require('lodash.isfunction');
const isEmpty = require('lodash.isempty');
const chalk = require('chalk');
const uuid = require('uuid');
const mochaUtils = require('mocha/lib/utils');
const stringify = require('json-stringify-safe');
const diff = require('diff');
const stripAnsi = require('strip-ansi');
const stripFnStart = require('./stripFnStart');

/**
 * Logger utility
 *
 * @param {String} msg - message to log
 * @param {String} level - log level [log, info, warn, error]
 * @param {Object} config - configuration object
 */
function log(msg, level, config) {
  // Don't log messages in quiet mode
  if (config && config.quiet) return;
  const logMethod = console[level] || console.log;
  let out = msg;
  if (typeof msg === 'object') {
    out = stringify(msg, null, 2);
  }
  logMethod(`[${chalk.gray('mochawesome')}] ${out}\n`);
}

/**
 * Strip the function definition from `str`,
 * and re-indent for pre whitespace.
 *
 * @param {String} str - code in
 *
 * @return {String} cleaned code string
 */
function cleanCode(str) {
  str = str
    .replace(/\r\n|[\r\n\u2028\u2029]/g, '\n') // unify linebreaks
    .replace(/^\uFEFF/, ''); // replace zero-width no-break space

  str = stripFnStart(str) // replace function declaration
    .replace(/\)\s*\)\s*$/, ')') // replace closing paren
    .replace(/\s*};?\s*$/, ''); // replace closing bracket

  // Preserve indentation by finding leading tabs/spaces
  // and removing that amount of space from each line
  const spaces = str.match(/^\n?( *)/)[1].length;
  const tabs = str.match(/^\n?(\t*)/)[1].length;
  /* istanbul ignore next */
  const indentRegex = new RegExp(
    `^\n?${tabs ? '\t' : ' '}{${tabs || spaces}}`,
    'gm'
  );

  str = str.replace(indentRegex, '').trim();
  return str;
}

/**
 * Create a unified diff between two strings
 *
 * @param {Error}  err          Error object
 * @param {string} err.actual   Actual result returned
 * @param {string} err.expected Result expected
 *
 * @return {string} diff
 */
function createUnifiedDiff({ actual, expected }) {
  return diff
    .createPatch('string', actual, expected)
    .split('\n')
    .splice(4)
    .map(line => {
      if (line.match(/@@/)) {
        return null;
      }
      if (line.match(/\\ No newline/)) {
        return null;
      }
      return line.replace(/^(-|\+)/, '$1 ');
    })
    .filter(line => typeof line !== 'undefined' && line !== null)
    .join('\n');
}

/**
 * Create an inline diff between two strings
 *
 * @param {Error}  err          Error object
 * @param {string} err.actual   Actual result returned
 * @param {string} err.expected Result expected
 *
 * @return {array} diff string objects
 */
function createInlineDiff({ actual, expected }) {
  return diff.diffWordsWithSpace(actual, expected);
}

/**
 * Return a normalized error object
 *
 * @param {Error} err Error object
 *
 * @return {Object} normalized error
 */
function normalizeErr(err, config) {
  const { name, message, actual, expected, stack, showDiff } = err;
  let errMessage;
  let errDiff;

  /**
   * Check that a / b have the same type.
   */
  function sameType(a, b) {
    const objToString = Object.prototype.toString;
    return objToString.call(a) === objToString.call(b);
  }

  // Format actual/expected for creating diff
  if (
    showDiff !== false &&
    sameType(actual, expected) &&
    expected !== undefined
  ) {
    /* istanbul ignore if */
    if (!(isString(actual) && isString(expected))) {
      err.actual = mochaUtils.stringify(actual);
      err.expected = mochaUtils.stringify(expected);
    }
    errDiff = config.useInlineDiffs
      ? createInlineDiff(err)
      : createUnifiedDiff(err);
  }

  // Assertion libraries do not output consitent error objects so in order to
  // get a consistent message object we need to create it ourselves
  if (name && message) {
    errMessage = `${name}: ${stripAnsi(message)}`;
  } else if (stack) {
    errMessage = stack.replace(/\n.*/g, '');
  }

  return {
    message: errMessage,
    estack: stack && stripAnsi(stack),
    diff: errDiff,
  };
}

/**
 * Return a plain-object representation of `test`
 * free of cyclic properties etc.
 *
 * @param {Object} test
 *
 * @return {Object} cleaned test
 */
function cleanTest(test, config) {
  const code = config.code ? test.body || '' : '';

  const fullTitle = isFunction(test.fullTitle)
    ? stripAnsi(test.fullTitle())
    : stripAnsi(test.title);

  const cleaned = {
    title: stripAnsi(test.title),
    fullTitle,
    timedOut: test.timedOut,
    duration: test.duration || 0,
    state: test.state,
    speed: test.speed,
    pass: test.state === 'passed',
    fail: test.state === 'failed',
    pending: test.pending,
    context: stringify(test.context, null, 2),
    code: code && cleanCode(code),
    err: (test.err && normalizeErr(test.err, config)) || {},
    uuid: test.uuid || /* istanbul ignore next: default */ uuid.v4(),
    parentUUID: test.parent && test.parent.uuid,
    isHook: test.type === 'hook',
  };

  cleaned.skipped =
    !cleaned.pass && !cleaned.fail && !cleaned.pending && !cleaned.isHook;

  return cleaned;
}

/**
 * Return a plain-object representation of `suite` with additional properties for rendering.
 *
 * @param {Object} suite
 * @param {Object} testTotals  Cumulative count of tests registered/skipped
 * @param {Integer} testTotals.registered
 * @param {Integer} testTotals.skipped
 *
 * @return {Object|boolean} cleaned suite or false if suite is empty
 */
function cleanSuite(suite, testTotals, config) {
  let duration = 0;
  const passingTests = [];
  const failingTests = [];
  const pendingTests = [];
  const skippedTests = [];

  const beforeHooks = []
    .concat(suite._beforeAll, suite._beforeEach)
    .map(test => cleanTest(test, config));

  const afterHooks = []
    .concat(suite._afterAll, suite._afterEach)
    .map(test => cleanTest(test, config));

  const tests = suite.tests.map(test => {
    const cleanedTest = cleanTest(test, config);
    duration += test.duration || 0;
    if (cleanedTest.state === 'passed') passingTests.push(cleanedTest.uuid);
    if (cleanedTest.state === 'failed') failingTests.push(cleanedTest.uuid);
    if (cleanedTest.pending) pendingTests.push(cleanedTest.uuid);
    if (cleanedTest.skipped) skippedTests.push(cleanedTest.uuid);
    return cleanedTest;
  });

  testTotals.registered += tests.length;
  testTotals.skipped += skippedTests.length;

  const cleaned = {
    uuid: suite.uuid || /* istanbul ignore next: default */ uuid.v4(),
    title: stripAnsi(suite.title),
    fullFile: suite.file || '',
    file: suite.file ? suite.file.replace(process.cwd(), '') : '',
    beforeHooks,
    afterHooks,
    tests,
    suites: suite.suites,
    passes: passingTests,
    failures: failingTests,
    pending: pendingTests,
    skipped: skippedTests,
    duration,
    root: suite.root,
    rootEmpty: suite.root && tests.length === 0,
    _timeout: suite._timeout,
  };

  const isEmptySuite =
    isEmpty(cleaned.suites) &&
    isEmpty(cleaned.tests) &&
    isEmpty(cleaned.beforeHooks) &&
    isEmpty(cleaned.afterHooks);

  return !isEmptySuite && cleaned;
}

/**
 * Map over a suite, returning a cleaned suite object
 * and recursively cleaning any nested suites.
 *
 * @param {Object} suite          Suite to map over
 * @param {Object} testTotals  Cumulative count of tests registered/skipped
 * @param {Integer} testTotals.registered
 * @param {Integer} testTotals.skipped
 * @param {Object} config         Reporter configuration
 */
function mapSuites(suite, testTotals, config) {
  const suites = suite.suites.reduce((acc, subSuite) => {
    const mappedSuites = mapSuites(subSuite, testTotals, config);
    if (mappedSuites) {
      acc.push(mappedSuites);
    }
    return acc;
  }, []);
  const toBeCleaned = { ...suite, suites };
  return cleanSuite(toBeCleaned, testTotals, config);
}

module.exports = {
  log,
  cleanCode,
  cleanTest,
  cleanSuite,
  mapSuites,
};
