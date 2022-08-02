const tv4 = require('tv4');
const semver = require('semver');
const schema = require('./installLogsCollector.schema.json');
const CtrError = require('./CtrError');
const LOG_TYPE = require('./constants').LOG_TYPES;
const tv4ErrorTransformer = require('./tv4ErrorTransformer');

const LogCollectBrowserConsole = require("./collector/LogCollectBrowserConsole");
const LogCollectCypressCommand = require("./collector/LogCollectCypressCommand");
const LogCollectCypressRequest = require("./collector/LogCollectCypressRequest");
const LogCollectCypressRoute = require("./collector/LogCollectCypressRoute");
const LogCollectCypressIntercept = require("./collector/LogCollectCypressIntercept");
const LogCollectCypressXhr = require("./collector/LogCollectCypressXhr");
const LogCollectCypressFetch = require("./collector/LogCollectCypressFetch");
const LogCollectCypressLog = require("./collector/LogCollectCypressLog");

const LogCollectorState = require("./collector/LogCollectorState");
const LogCollectExtendedControl = require("./collector/LogCollectExtendedControl");
const LogCollectSimpleControl = require("./collector/LogCollectSimpleControl");

/**
 * Installs the logs collector for cypress.
 *
 * Needs to be added to support file.
 *
 * @see ./installLogsCollector.d.ts
 */
function installLogsCollector(config = {}) {
  validateConfig(config);

  config.collectTypes = config.collectTypes || Object.values(LOG_TYPE);
  config.collectRequestData = config.xhr && config.xhr.printRequestData;
  config.collectHeaderData = config.xhr && config.xhr.printHeaderData;

  let logCollectorState = new LogCollectorState(config);
  registerLogCollectorTypes(logCollectorState, config);

  if (config.enableExtendedCollector) {
    (new LogCollectExtendedControl(logCollectorState, config)).register();
  } else {
    (new LogCollectSimpleControl(logCollectorState, config)).register();
  }
}

function registerLogCollectorTypes(logCollectorState, config) {
  (new LogCollectBrowserConsole(logCollectorState, config)).register()

  if (config.collectTypes.includes(LOG_TYPE.CYPRESS_LOG)) {
    (new LogCollectCypressLog(logCollectorState, config)).register();
  }
  if (config.collectTypes.includes(LOG_TYPE.CYPRESS_XHR)) {
    (new LogCollectCypressXhr(logCollectorState, config)).register();
  }
  if (config.collectTypes.includes(LOG_TYPE.CYPRESS_FETCH)) {
    (new LogCollectCypressFetch(logCollectorState, config)).register();
  }
  if (config.collectTypes.includes(LOG_TYPE.CYPRESS_REQUEST)) {
    (new LogCollectCypressRequest(logCollectorState, config)).register();
  }
  if (config.collectTypes.includes(LOG_TYPE.CYPRESS_ROUTE)) {
    (new LogCollectCypressRoute(logCollectorState, config)).register();
  }
  if (config.collectTypes.includes(LOG_TYPE.CYPRESS_COMMAND)) {
    (new LogCollectCypressCommand(logCollectorState, config)).register();
  }
  if (config.collectTypes.includes(LOG_TYPE.CYPRESS_INTERCEPT) && semver.gte(Cypress.version, '6.0.0')) {
    (new LogCollectCypressIntercept(logCollectorState, config)).register();
  }
}

function validateConfig(config) {
  const result = tv4.validateMultiple(config, schema);

  if (!result.valid) {
    throw new CtrError(
      `Invalid plugin install options: ${tv4ErrorTransformer.toReadableString(
        result.errors
      )}`
    );
  }

  if (config.filterLog && typeof config.filterLog !== 'function') {
    throw new CtrError(`Filter log option expected to be a function.`);
  }
  if (config.processLog && typeof config.processLog !== 'function') {
    throw new CtrError(`Process log option expected to be a function.`);
  }
  if (config.collectTestLogs && typeof config.collectTestLogs !== 'function') {
    throw new CtrError(`Collect test logs option expected to be a function.`);
  }
}


module.exports = installLogsCollector;
