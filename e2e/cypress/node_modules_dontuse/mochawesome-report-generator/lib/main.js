"use strict";

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; var ownKeys = Object.keys(source); if (typeof Object.getOwnPropertySymbols === 'function') { ownKeys = ownKeys.concat(Object.getOwnPropertySymbols(source).filter(function (sym) { return Object.getOwnPropertyDescriptor(source, sym).enumerable; })); } ownKeys.forEach(function (key) { _defineProperty(target, key, source[key]); }); } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

var fs = require('fs-extra');

var fsu = require('fsu');

var path = require('path');

var opener = require('opener');

var dateFormat = require('dateformat');

var renderMainHTML = require('./main-html');

var pkg = require('../package.json');

var _require = require('./options'),
    getMergedOptions = _require.getMergedOptions;

var distDir = path.join(__dirname, '..', 'dist');
var fileExtRegex = /\.[^.]*?$/;
var htmlJsonExtRegex = /\.(?:html|json)$/;
var semverRegex = /\d+\.\d+\.\d+(?:-(alpha|beta)\.\d+)?/;
/**
 * Saves a file
 *
 * @param {string} filename Name of file to save
 * @param {string} data Data to be saved
 * @param {boolean} overwrite Overwrite existing files (default: true)
 *
 * @return {Promise} Resolves with filename if successfully saved
 */

function saveFile(filename, data, overwrite) {
  if (overwrite) {
    return fs.outputFile(filename, data).then(function () {
      return filename;
    });
  }

  return new Promise(function (resolve, reject) {
    fsu.writeFileUnique(filename.replace(fileExtRegex, '{_###}$&'), data, {
      force: true
    }, function (err, savedFile) {
      return err === null ? resolve(savedFile) : reject(err);
    });
  });
}
/**
 * Opens a file
 *
 * @param {string} filename Name of file to open
 *
 * @return {Promise} Resolves with filename if successfully opened
 */


function openFile(filename) {
  return new Promise(function (resolve, reject) {
    opener(filename, null, function (err) {
      return err === null ? resolve(filename) : reject(err);
    });
  });
}
/**
 * Synchronously loads a file with utf8 encoding
 *
 * @param {string} filename Name of file to load
 *
 * @return {string} File data as string
 */


function loadFile(filename) {
  return fs.readFileSync(filename, 'utf8');
}
/**
 * Get the dateformat format string based on the timestamp option
 *
 * @param {string|boolean} ts Timestamp option value
 *
 * @return {string} Valid dateformat format string
 */


function getTimestampFormat(ts) {
  return ts === '' || ts === true || ts === 'true' || ts === false || ts === 'false' ? 'isoDateTime' : ts;
}
/**
 * Construct the path/name of the HTML/JSON file to be saved
 *
 * @param {object} reportOptions Options object
 * @param {string} reportOptions.reportDir Directory to save report to
 * @param {string} reportOptions.reportFilename Filename to save report to
 * @param {string} reportOptions.timestamp Timestamp format to be appended to the filename
 * @param {object} reportData JSON test data
 *
 * @return {string} Fully resolved path without extension
 */


function getFilename(_ref, reportData) {
  var reportDir = _ref.reportDir,
      reportFilename = _ref.reportFilename,
      timestamp = _ref.timestamp;
  var DEFAULT_FILENAME = 'mochawesome';
  var NAME_REPLACE = '[name]';
  var STATUS_REPLACE = '[status]';
  var DATETIME_REPLACE = '[datetime]';
  var STATUSES = {
    Pass: 'pass',
    Fail: 'fail'
  };
  var filename = reportFilename || DEFAULT_FILENAME;
  var hasDatetimeReplacement = filename.includes(DATETIME_REPLACE);
  var tsFormat = getTimestampFormat(timestamp);
  var ts = dateFormat(new Date(), tsFormat) // replace commas, spaces or comma-space combinations with underscores
  .replace(/(,\s*)|,|\s+/g, '_') // replace forward and back slashes with hyphens
  .replace(/\\|\//g, '-') // remove colons
  .replace(/:/g, '');

  if (timestamp !== false && timestamp !== 'false') {
    if (!hasDatetimeReplacement) {
      filename = "".concat(filename, "_").concat(DATETIME_REPLACE);
    }
  }

  var specFilename = path.basename(reportData.results[0].file || '').replace(/\..+/, '');
  var status = reportData.stats.failures > 0 ? STATUSES.Fail : STATUSES.Pass;
  filename = filename.replace(NAME_REPLACE, specFilename || DEFAULT_FILENAME).replace(STATUS_REPLACE, status).replace(DATETIME_REPLACE, ts).replace(htmlJsonExtRegex, '');
  return path.resolve(process.cwd(), reportDir, filename);
}
/**
 * Get report options by extending base options
 * with user provided options
 *
 * @param {object} opts Report options
 * @param {object} reportData JSON test data
 *
 * @return {object} User options merged with default options
 */


function getOptions(opts, reportData) {
  var mergedOptions = getMergedOptions(opts || {});
  var filename = getFilename(mergedOptions, reportData); // For saving JSON from mochawesome reporter

  if (mergedOptions.saveJson) {
    mergedOptions.jsonFile = "".concat(filename, ".json");
  }

  mergedOptions.htmlFile = "".concat(filename, ".html");
  return mergedOptions;
}
/**
 * Determine if assets should be copied following below logic:
 * - Assets folder does not exist -> copy assets
 * - Assets folder exists -> load the css asset to inspect the banner
 * - Error loading css file -> copy assets
 * - Read the package version from the css asset
 * - Asset version is not found -> copy assets
 * - Asset version differs from current version -> copy assets
 *
 * @param {string} assetsDir Directory where assets should be saved
 *
 * @return {boolean} Should assets be copied
 */


function _shouldCopyAssets(assetsDir) {
  if (!fs.existsSync(assetsDir)) {
    return true;
  }

  try {
    var appCss = loadFile(path.join(assetsDir, 'app.css'));
    var appCssVersion = semverRegex.exec(appCss);

    if (!appCssVersion || appCssVersion[0] !== pkg.version) {
      return true;
    }
  } catch (e) {
    return true;
  }

  return false;
}
/**
 * Copy the report assets to the report dir, ignoring inline assets
 *
 * @param {object} opts Report options
 */


function copyAssets(_ref2) {
  var assetsDir = _ref2.assetsDir;

  if (_shouldCopyAssets(assetsDir)) {
    fs.copySync(distDir, assetsDir, {
      filter: function filter(src) {
        return !/inline/.test(src);
      }
    });
  }
}
/**
 * Get the report assets object
 *
 * @param {object} reportOptions Options
 * @return {object} Object with assets props
 */


function getAssets(reportOptions) {
  var assetsDir = reportOptions.assetsDir,
      cdn = reportOptions.cdn,
      dev = reportOptions.dev,
      inlineAssets = reportOptions.inlineAssets,
      reportDir = reportOptions.reportDir;
  var relativeAssetsDir = path.relative(reportDir, assetsDir); // Default URLs to assets path

  var assets = {
    inlineScripts: null,
    inlineStyles: null,
    scriptsUrl: path.join(relativeAssetsDir, 'app.js'),
    stylesUrl: path.join(relativeAssetsDir, 'app.css')
  }; // If using inline assets, load files and strings

  if (inlineAssets) {
    assets.inlineScripts = loadFile(path.join(distDir, 'app.js'));
    assets.inlineStyles = loadFile(path.join(distDir, 'app.inline.css'));
  } // If using CDN, return remote urls


  if (cdn) {
    assets.scriptsUrl = "https://unpkg.com/mochawesome-report-generator@".concat(pkg.version, "/dist/app.js");
    assets.stylesUrl = "https://unpkg.com/mochawesome-report-generator@".concat(pkg.version, "/dist/app.css");
  } // In DEV mode, return local urls


  if (dev) {
    assets.scriptsUrl = 'http://localhost:8080/app.js';
    assets.stylesUrl = 'http://localhost:8080/app.css';
  } // Copy the assets if needed


  if (!dev && !cdn && !inlineAssets) {
    copyAssets(reportOptions);
  }

  return assets;
}
/**
 * Prepare options, assets, and html for saving
 *
 * @param {object} reportData JSON test data
 * @param {object} opts Report options
 *
 * @return {object} Prepared data for saving
 */


function prepare(reportData, opts) {
  // Get the options
  var reportOptions = getOptions(opts, reportData); // Stop here if we're not generating an HTML report

  if (!reportOptions.saveHtml) {
    return {
      reportOptions: reportOptions
    };
  } // Get the assets


  var assets = getAssets(reportOptions); // Render basic template to string

  var renderedHtml = renderMainHTML(_objectSpread({
    data: JSON.stringify(reportData),
    options: reportOptions,
    title: reportOptions.reportPageTitle,
    useInlineAssets: reportOptions.inlineAssets && !reportOptions.cdn
  }, assets));
  var html = "<!doctype html>\n".concat(renderedHtml);
  return {
    html: html,
    reportOptions: reportOptions
  };
}
/**
 * Create the report
 *
 * @param {object} data JSON test data
 * @param {object} opts Report options
 *
 * @return {Promise} Resolves if report was created successfully
 */


function create(data, opts) {
  var _prepare = prepare(data, opts),
      html = _prepare.html,
      reportOptions = _prepare.reportOptions;

  var saveJson = reportOptions.saveJson,
      saveHtml = reportOptions.saveHtml,
      autoOpen = reportOptions.autoOpen,
      overwrite = reportOptions.overwrite,
      jsonFile = reportOptions.jsonFile,
      htmlFile = reportOptions.htmlFile;
  var savePromises = [];
  savePromises.push(saveHtml !== false ? saveFile(htmlFile, html, overwrite).then(function (savedHtml) {
    return autoOpen && openFile(savedHtml) || savedHtml;
  }) : null);
  savePromises.push(saveJson ? saveFile(jsonFile, // Preserve `undefined` values as `null` when stringifying
  JSON.stringify(data, function (k, v) {
    return v === undefined ? null : v;
  }, 2), overwrite) : null);
  return Promise.all(savePromises);
}
/**
 * Create the report synchronously
 *
 * @param {object} data JSON test data
 * @param {object} opts Report options
 *
 */


function createSync(data, opts) {
  var _prepare2 = prepare(data, opts),
      html = _prepare2.html,
      reportOptions = _prepare2.reportOptions;

  var autoOpen = reportOptions.autoOpen,
      htmlFile = reportOptions.htmlFile;
  fs.outputFileSync(htmlFile, html);
  if (autoOpen) opener(htmlFile);
}
/**
 * Expose functions
 *
 */


module.exports = {
  create: create,
  createSync: createSync
};