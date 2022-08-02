const Base = require('mocha/lib/reporters/base');
const mochaPkg = require('mocha/package.json');
const uuid = require('uuid');
const marge = require('mochawesome-report-generator');
const margePkg = require('mochawesome-report-generator/package.json');
const conf = require('./config');
const utils = require('./utils');
const pkg = require('../package.json');
const Mocha = require('mocha');
const { EVENT_SUITE_END } = Mocha.Runner.constants;

// Import the utility functions
const { log, mapSuites } = utils;

// Track the total number of tests registered/skipped
const testTotals = {
  registered: 0,
  skipped: 0,
};

/**
 * Done function gets called before mocha exits
 *
 * Creates and saves the report HTML and JSON files
 *
 * @param {Object} output    Final report object
 * @param {Object} options   Options to pass to report generator
 * @param {Object} config    Reporter config object
 * @param {Number} failures  Number of reported failures
 * @param {Function} exit
 *
 * @return {Promise} Resolves with successful report creation
 */
function done(output, options, config, failures, exit) {
  return marge
    .create(output, options)
    .then(([htmlFile, jsonFile]) => {
      if (!htmlFile && !jsonFile) {
        log('No files were generated', 'warn', config);
      } else {
        jsonFile && log(`Report JSON saved to ${jsonFile}`, null, config);
        htmlFile && log(`Report HTML saved to ${htmlFile}`, null, config);
      }
    })
    .catch(err => {
      log(err, 'error', config);
    })
    .then(() => {
      exit && exit(failures > 0 ? 1 : 0);
    });
}

/**
 * Get the class of the configured console reporter. This reporter outputs
 * test results to the console while mocha is running, and before
 * mochawesome generates its own report.
 *
 * Defaults to 'spec'.
 *
 * @param {String} reporter   Name of reporter to use for console output
 *
 * @return {Object} Reporter class object
 */
function consoleReporter(reporter) {
  if (reporter) {
    try {
      return require(`mocha/lib/reporters/${reporter}`);
    } catch (e) {
      log(`Unknown console reporter '${reporter}', defaulting to spec`);
    }
  }

  return require('mocha/lib/reporters/spec');
}

/**
 * Initialize a new reporter.
 *
 * @param {Runner} runner
 * @api public
 */
function Mochawesome(runner, options) {
  // Set the config options
  this.config = conf(options);

  // Ensure stats collector has been initialized
  if (!runner.stats) {
    const createStatsCollector = require('mocha/lib/stats-collector');
    createStatsCollector(runner);
  }

  // Reporter options
  const reporterOptions = {
    ...options.reporterOptions,
    reportFilename: this.config.reportFilename,
    saveHtml: this.config.saveHtml,
    saveJson: this.config.saveJson,
  };

  // Done function will be called before mocha exits
  // This is where we will save JSON and generate the HTML report
  this.done = (failures, exit) =>
    done(this.output, reporterOptions, this.config, failures, exit);

  // Reset total tests counters
  testTotals.registered = 0;
  testTotals.skipped = 0;

  // Call the Base mocha reporter
  Base.call(this, runner);

  const reporterName = reporterOptions.consoleReporter;
  if (reporterName !== 'none') {
    const ConsoleReporter = consoleReporter(reporterName);
    new ConsoleReporter(runner); // eslint-disable-line
  }

  let endCalled = false;

  // Add a unique identifier to each suite/test/hook
  ['suite', 'test', 'hook', 'pending'].forEach(type => {
    runner.on(type, item => {
      item.uuid = uuid.v4();
    });
  });

  // Handle events from workers in parallel mode
  if (runner.constructor.name === 'ParallelBufferedRunner') {
    const setSuiteDefaults = suite => {
      [
        'suites',
        'tests',
        '_beforeAll',
        '_beforeEach',
        '_afterEach',
        '_afterAll',
      ].forEach(field => {
        suite[field] = suite[field] || [];
      });
      suite.suites.forEach(it => setSuiteDefaults(it));
    };

    runner.on(EVENT_SUITE_END, function (suite) {
      if (suite.root) {
        setSuiteDefaults(suite);
        runner.suite.suites.push(...suite.suites);
      }
    });
  }

  // Process the full suite
  runner.on('end', () => {
    try {
      /* istanbul ignore else */
      if (!endCalled) {
        // end gets called more than once for some reason
        // so we ensure the suite is processed only once
        endCalled = true;

        const rootSuite = mapSuites(this.runner.suite, testTotals, this.config);

        // Attempt to set a filename for the root suite to
        // support `reportFilename` [name] replacement token
        if (rootSuite) {
          if (rootSuite.suites.length === 1) {
            const firstSuite = rootSuite.suites[0];
            rootSuite.file = firstSuite.file || rootSuite.file;
            rootSuite.fullFile = firstSuite.fullFile || rootSuite.fullFile;
          } else if (!rootSuite.suites.length && rootSuite.tests.length) {
            const firstTest = this.runner.suite.tests[0];
            rootSuite.file = firstTest.file || rootSuite.file;
            rootSuite.fullFile = firstTest.fullFile || rootSuite.fullFile;
          }
        }

        const obj = {
          stats: this.stats,
          results: [rootSuite],
          meta: {
            mocha: {
              version: mochaPkg.version,
            },
            mochawesome: {
              options: this.config,
              version: pkg.version,
            },
            marge: {
              options: options.reporterOptions,
              version: margePkg.version,
            },
          },
        };

        obj.stats.testsRegistered = testTotals.registered;

        const { passes, failures, pending, tests, testsRegistered } = obj.stats;
        const passPercentage = (passes / (testsRegistered - pending)) * 100;
        const pendingPercentage = (pending / testsRegistered) * 100;

        obj.stats.passPercent = passPercentage;
        obj.stats.pendingPercent = pendingPercentage;
        obj.stats.other = passes + failures + pending - tests; // Failed hooks
        obj.stats.hasOther = obj.stats.other > 0;
        obj.stats.skipped = testTotals.skipped;
        obj.stats.hasSkipped = obj.stats.skipped > 0;
        obj.stats.failures -= obj.stats.other;

        // Save the final output to be used in the done function
        this.output = obj;
      }
    } catch (e) {
      // required because thrown errors are not handled directly in the
      // event emitter pattern and mocha does not have an "on error"
      /* istanbul ignore next */
      log(`Problem with mochawesome: ${e.stack}`, 'error');
    }
  });
}

module.exports = Mochawesome;
