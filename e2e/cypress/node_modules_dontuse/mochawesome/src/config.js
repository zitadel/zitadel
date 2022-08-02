/**
 * Retrieve the value of a user supplied option.
 * Falls back to `defaultValue`
 * Order of precedence
 *  1. User-supplied option
 *  2. Environment variable
 *  3. Default value
 *
 * @param {string} optToGet  Option name
 * @param {object} options  User supplied options object
 * @param {boolean} isBool  Treat option as Boolean
 * @param {string|boolean} defaultValue  Fallback value
 *
 * @return {string|boolean}  Option value
 */
function _getOption(optToGet, options, isBool, defaultValue) {
  const envVar = `MOCHAWESOME_${optToGet.toUpperCase()}`;
  if (options && typeof options[optToGet] !== 'undefined') {
    return isBool && typeof options[optToGet] === 'string'
      ? options[optToGet] === 'true'
      : options[optToGet];
  }
  if (typeof process.env[envVar] !== 'undefined') {
    return isBool ? process.env[envVar] === 'true' : process.env[envVar];
  }
  return defaultValue;
}

module.exports = function (opts) {
  const reporterOpts = (opts && opts.reporterOptions) || {};
  const code = _getOption('code', reporterOpts, true, true);
  const noCode = _getOption('no-code', reporterOpts, true, false);

  return {
    quiet: _getOption('quiet', reporterOpts, true, false),
    reportFilename: _getOption(
      'reportFilename',
      reporterOpts,
      false,
      'mochawesome'
    ),
    saveHtml: _getOption('html', reporterOpts, true, true),
    saveJson: _getOption('json', reporterOpts, true, true),
    consoleReporter: _getOption('consoleReporter', reporterOpts, false, 'spec'),
    useInlineDiffs: !!opts.inlineDiffs,
    code: noCode ? false : code,
  };
};
