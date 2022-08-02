"use strict";

const fs = require('fs-extra');

const path = require('path');

const chalk = require('chalk');

const t = require('tcomb-validation');

const dateFormat = require('dateformat');

const report = require('../lib/main');

const types = require('./types');

const logger = require('./logger');

const JsonErrRegex = /Unexpected token/;
const JsonFileRegex = /\.json{1}$/;
const htmlJsonExtRegex = /\.(?:html|json)$/;

const mapJsonErrors = errors => errors.map(e => `  ${e.message}`).join('\n');

const ERRORS = {
  NOT_FOUND: '  File not found.',
  GENERIC: '  There was a problem loading mochawesome data.',
  INVALID_JSON: errMsgs => mapJsonErrors(errMsgs)
};
let validFiles;
/**
 * Validate the data file
 *
 * @typedef {Object} File
 * @property {string} filename Name of the file
 * @property {object} data JSON test data
 * @property {object} err Error object
 *
 * @param {string} file File to load/validate
 *
 * @return {File} Validated file with test data, `err` will be null if valid
 */

function validateFile(file) {
  let data;
  let err = null; // Try to read and parse the file

  try {
    data = JSON.parse(fs.readFileSync(file, 'utf-8'));
  } catch (e) {
    if (e.code === 'ENOENT') {
      err = ERRORS.NOT_FOUND;
    } else if (JsonErrRegex.test(e.message)) {
      err = ERRORS.INVALID_JSON([e]);
    } else {
      err = ERRORS.GENERIC;
    }
  } // If the file was loaded successfully,
  // validate the json against the TestReport schema


  if (data) {
    const validationResult = t.validate(data, types.TestReport, {
      strict: true
    });

    if (!validationResult.isValid()) {
      err = ERRORS.INVALID_JSON(validationResult.errors);
    } else {
      validFiles += 1;
    }
  }

  return {
    filename: file,
    data,
    err
  };
}
/**
 * Set exit code and throw caught errors
 *
 * @param {Object|string} err Error object or error message
 *
 */


function handleError(err) {
  process.exitCode = 1;
  throw new Error(err);
}
/**
 * Loop through resolved promises to log the appropriate messages
 *
 * @param {Array} resolvedValues Result of promise.all
 *
 * @return {Array} Array of resolved promise values
 */


function handleResolved(resolvedValues) {
  const saved = [];
  const errors = [];
  resolvedValues.forEach(value => {
    if (value.err) {
      errors.push(value);
    } else {
      saved.push(value[0]);
    }
  });

  if (saved.length) {
    logger.info(chalk.green('\n✓ Reports saved:'));
    logger.info(saved.map(savedFile => `${chalk.underline(savedFile)}`).join('\n'));
  }

  if (errors.length) {
    logger.info(chalk.red('\n✘ Some files could not be processed:'));
    logger.info(errors.map(e => `${chalk.underline(e.filename)}\n${chalk.dim(e.err)}`).join('\n\n'));
    process.exitCode = 1;
  }

  if (!validFiles && !errors.length) {
    logger.info(chalk.yellow('\nDid not find any JSON files to process.'));
    process.exitCode = 1;
  }

  return resolvedValues;
}
/**
 * Get the dateformat format string based on the timestamp option
 *
 * @param {string|boolean} ts Timestamp option value
 *
 * @return {string} Valid dateformat format string
 */


function getTimestampFormat(ts) {
  return ts === undefined || ts === true || ts === 'true' || ts === false || ts === 'false' ? 'isoDateTime' : ts;
}
/**
 * Get the reportFilename option to be passed to `report.create`
 *
 * Returns the `reportFilename` option if provided otherwise
 * it returns the base filename stripped of path and extension
 *
 * @param {Object} file File object
 * @param {string} file.filename Name of file to be processed
 * @param {Object} file.data JSON test data
 * @param {Object} args CLI process arguments
 *
 * @return {string} Filename
 */


function getReportFilename({
  filename,
  data
}, {
  reportFilename,
  timestamp
}) {
  const DEFAULT_FILENAME = filename.split(path.sep).pop().replace(JsonFileRegex, '');
  const NAME_REPLACE = '[name]';
  const STATUS_REPLACE = '[status]';
  const DATETIME_REPLACE = '[datetime]';
  const STATUSES = {
    Pass: 'pass',
    Fail: 'fail'
  };
  let outFilename = reportFilename || DEFAULT_FILENAME;
  const hasDatetimeReplacement = outFilename.includes(DATETIME_REPLACE);
  const tsFormat = getTimestampFormat(timestamp);
  const ts = dateFormat(new Date(), tsFormat) // replace commas, spaces or comma-space combinations with underscores
  .replace(/(,\s*)|,|\s+/g, '_') // replace forward and back slashes with hyphens
  .replace(/\\|\//g, '-') // remove colons
  .replace(/:/g, '');

  if (timestamp) {
    if (!hasDatetimeReplacement) {
      outFilename = `${outFilename}_${DATETIME_REPLACE}`;
    }
  } // Special handling of replacement tokens


  const status = data.stats.failures > 0 ? STATUSES.Fail : STATUSES.Pass;
  outFilename = outFilename.replace(NAME_REPLACE, DEFAULT_FILENAME).replace(STATUS_REPLACE, status).replace(DATETIME_REPLACE, ts).replace(htmlJsonExtRegex, '');
  return outFilename;
}
/**
 * Process arguments, recursing through any directories,
 * to find and validate JSON files
 *
 * @param {array} args Array of paths
 * @param {array} files Array to populate
 *
 * @return {array} File objects to be processed
 */


function processArgs(args, files = []) {
  return args.reduce((acc, arg) => {
    let stats;

    try {
      stats = fs.statSync(arg);
    } catch (err) {} // Do nothing
    // If argument is a directory, process the files inside


    if (stats && stats.isDirectory()) {
      return processArgs(fs.readdirSync(arg).map(file => path.join(arg, file)), files);
    } // If `statSync` failed, validating will handle the error
    // If the argument is a file, check if its a JSON file before validating


    if (!stats || JsonFileRegex.test(arg)) {
      acc.push(validateFile(arg));
    }

    return acc;
  }, files);
}
/**
 * Main CLI Program
 *
 * @param {Object} args CLI arguments
 *
 * @return {Promise} Resolved promises with saved files or errors
 */


function marge(args) {
  // Reset valid files count
  validFiles = 0;
  const newArgs = Object.assign({}, args); // Get the array of JSON files to process

  const files = processArgs(args._); // When there are multiple valid files OR the timestamp option is set
  // we must force `overwrite` to `false` to ensure all reports are created

  /* istanbul ignore else */

  if (validFiles > 1 || args.timestamp !== false) {
    newArgs.overwrite = false;
  }

  const promises = files.map(file => {
    // Files with errors we just resolve
    if (file.err) {
      return Promise.resolve(file);
    } // Valid files get created but first we need to pass correct filename option
    // Default value is name of file
    // If a filename option was provided, all files get that name


    const reportFilename = getReportFilename(file, newArgs);
    return report.create(file.data, Object.assign({}, newArgs, {
      reportFilename
    }));
  });
  return Promise.all(promises).then(handleResolved).catch(handleError);
}

module.exports = marge;