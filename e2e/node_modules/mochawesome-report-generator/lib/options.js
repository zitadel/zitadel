"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.getMergedOptions = exports.yargsOptions = void 0;

var path = require('path');

var isFunction = require('lodash.isfunction');
/** CLI Arguments
 *
 * @argument {string} test-data Data to use for rendering report
 *
 * @property {string}  reportFilename   Filename of saved report
 * @property {string}  reportDir        Path to save report to (default: cwd/mochawesome-report)
 * @property {string}  reportTitle      Title to use on the report (default: mochawesome)
 * @property {string}  reportPageTitle  Title of the report document (default: mochawesome-report)
 * @property {string}  assetsDir        Path to save report assets to (default: cwd/mochawesome-report/assets)
 * @property {boolean} inlineAssets     Should assets be inlined into HTML file (default: false)
 * @property {boolean} cdn              Should assets be loaded via CDN (default: false)
 * @property {boolean} charts           Should charts be enabled (default: false)
 * @property {boolean} code             Should test code output be enabled (default: true)
 * @property {boolean} autoOpen         Open the report after creation (default: false)
 * @property {boolean} overwrite        Overwrite existing files (default: true)
 * @property {string}  timestamp        Append timestamp in specified format to the filename.
 *                                      Ensures a new file is created on each run.
 *                                      Accepts any format string that dateformat can handle.
 *                                      See https://github.com/felixge/node-dateformat
 *                                      Defaults to 'isoDateTime' when no format is specified
 * @property {boolean} showPassed       Initial state of "Show Passed" filter (default: true)
 * @property {boolean} showFailed       Initial state of "Show Failed" filter (default: true)
 * @property {boolean} showPending      Initial state of "Show Pending" filter (default: true)
 * @property {boolean} showSkipped      Initial state of "Show Skipped" filter (default: false)
 * @property {string}  showHooks        Determines when hooks should display in the report
 *                                      Choices:
 *                                       - always: display all hooks
 *                                       - never: do not display hooks
 *                                       - failed: display only failed hooks (default)
 *                                       - context: display only hooks with context
 * @property {boolean} saveJson         Should report data be saved to JSON file (default: false)
 * @property {boolean} saveHtml         Should report be saved to HTML file (default: true)
 * @property {boolean} dev              Enable dev mode in the report,
 *                                      asssets loaded via webpack (default: false)
 */


var yargsOptions = {
  f: {
    alias: ['reportFilename'],
    describe: 'Filename of saved report',
    string: true,
    requiresArg: true
  },
  o: {
    alias: ['reportDir'],
    default: 'mochawesome-report',
    describe: 'Path to save report',
    string: true,
    normalize: true,
    requiresArg: true
  },
  t: {
    alias: ['reportTitle'],
    default: function _default() {
      return process.cwd().split(path.sep).pop();
    },
    describe: 'Report title',
    string: true,
    requiresArg: true
  },
  p: {
    alias: ['reportPageTitle'],
    default: 'Mochawesome Report',
    describe: 'Browser title',
    string: true,
    requiresArg: true
  },
  i: {
    alias: ['inline', 'inlineAssets'],
    default: false,
    describe: 'Inline report assets (styles, scripts)',
    boolean: true
  },
  assetsDir: {
    describe: 'Path to save assets',
    string: true,
    normalize: true,
    requiresArg: true
  },
  cdn: {
    default: false,
    describe: 'Load report assets via CDN',
    boolean: true
  },
  charts: {
    alias: ['enableCharts'],
    default: false,
    describe: 'Display charts',
    boolean: true
  },
  code: {
    alias: ['enableCode'],
    default: true,
    describe: 'Display test code',
    boolean: true
  },
  autoOpen: {
    default: false,
    describe: 'Automatically open the report HTML',
    boolean: true
  },
  overwrite: {
    default: true,
    describe: 'Overwrite existing files when saving',
    boolean: true
  },
  timestamp: {
    alias: ['ts'],
    default: false,
    describe: 'Append timestamp in specified format to filename',
    string: true
  },
  showPassed: {
    default: true,
    describe: 'Set intial state for "Show Passed" filter',
    boolean: true
  },
  showFailed: {
    default: true,
    describe: 'Set intial state for "Show Failed" filter',
    boolean: true
  },
  showPending: {
    default: true,
    describe: 'Set intial state for "Show Pending" filter',
    boolean: true
  },
  showSkipped: {
    default: false,
    describe: 'Set intial state for "Show Skipped" filter',
    boolean: true
  },
  showHooks: {
    default: 'failed',
    describe: 'Display hooks in the report',
    choices: ['always', 'never', 'failed', 'context']
  },
  saveJson: {
    default: false,
    describe: 'Save report data to JSON file',
    boolean: true
  },
  saveHtml: {
    default: true,
    describe: 'Save report to HTML file',
    boolean: true
  },
  dev: {
    default: false,
    describe: 'Enable dev mode',
    boolean: true
  }
};
/**
 * Retrieve the value of a user supplied option.
 * Order of precedence
 *  1. User-supplied option
 *  2. Environment variable
 *
 * @param {object}  userOptions  Options to parse through
 * @param {string}  optToGet     Option name
 * @param {boolean} isBool       Treat option as Boolean
 *
 * @return {string|boolean|undefined}  Option value
 */

exports.yargsOptions = yargsOptions;

function _getUserOption(userOptions, optToGet, isBool) {
  var envVar = "MOCHAWESOME_".concat(optToGet.toUpperCase());

  if (userOptions && typeof userOptions[optToGet] !== 'undefined') {
    return isBool && typeof userOptions[optToGet] === 'string' ? userOptions[optToGet] === 'true' : userOptions[optToGet];
  }

  if (typeof process.env[envVar] !== 'undefined') {
    return isBool ? process.env[envVar] === 'true' : process.env[envVar];
  }

  return undefined;
}
/*
 * Helper to create properties and assign values in an object.
 * Properties with `undefined` values are ignored. *mutative*

 * @param {object} obj Object to assign properties to
 */


function assignVal(obj, prop, userVal, defaultVal) {
  var val = userVal !== undefined ? userVal : defaultVal;

  if (val !== undefined) {
    obj[prop] = val; // eslint-disable-line
  }
}
/**
 * Return parsed user options merged with base config
 *
 * @param {Object} userOptions User-supplied options
 *
 * @return {Object} Merged options
 */


var getMergedOptions = function getMergedOptions(userOptions) {
  var mergedOptions = {};
  Object.keys(yargsOptions).forEach(function (optKey) {
    var yargOpt = yargsOptions[optKey];
    var aliases = yargOpt.alias;
    var defaultVal = isFunction(yargOpt.default) ? yargOpt.default() : yargOpt.default;
    var isBool = yargOpt.boolean;

    var userVal = _getUserOption(userOptions, optKey, isBool); // Most options are single-letter so we add the aliases as properties


    if (Array.isArray(aliases) && aliases.length) {
      // If the main prop does not have a user supplied value
      // we need to check the aliases, stopping if we get a user value
      if (userVal === undefined) {
        for (var i = 0; i < aliases.length; i += 1) {
          userVal = _getUserOption(userOptions, aliases[i], isBool);

          if (userVal !== undefined) {
            break;
          }
        }
      } // Handle cases where the main option is not a single letter


      if (optKey.length > 1) assignVal(mergedOptions, optKey, userVal, defaultVal); // Loop through aliases to set val

      aliases.forEach(function (alias) {
        return assignVal(mergedOptions, alias, userVal, defaultVal);
      });
    } else {
      // For options without aliases, use the option regardless of length
      assignVal(mergedOptions, optKey, userVal, defaultVal);
    }
  }); // Special handling for defining `assetsDir`

  if (!mergedOptions.assetsDir) {
    mergedOptions.assetsDir = path.join(mergedOptions.reportDir, 'assets');
  }

  return mergedOptions;
};

exports.getMergedOptions = getMergedOptions;