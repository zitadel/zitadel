const LOG_TYPE = require('../constants').LOG_TYPES;
const CONSTANTS = require('../constants');
const LogFormat = require("./LogFormat");

module.exports = class LogCollectCypressFetch {

  constructor(collectorState, config) {
    this.config = config;
    this.collectorState = collectorState;

    this.format = new LogFormat(config);
  }

  register() {
    const formatFetch = (options) => (options.alias !== undefined ? '(' + options.alias + ') ' : '') +
      (options.consoleProps["Request went to origin?"] !== 'yes' ? 'STUBBED ' : '') +
      options.consoleProps.Method + ' ' + options.consoleProps.URL;

    const formatDuration = (durationInMs) =>
      durationInMs < 1000 ? `${durationInMs} ms` : `${durationInMs / 1000} s`;

    Cypress.on('log:added', (options) => {
      if (options.instrument === 'command' && options.name === 'request' && options.displayName === 'fetch') {
        const log = formatFetch(options);
        const severity = options.state === 'failed' ? CONSTANTS.SEVERITY.WARNING : '';
        this.collectorState.addLog([LOG_TYPE.CYPRESS_FETCH, log, severity], options.id);
      }
    });

    Cypress.on('log:changed', async (options) => {
      if (
        options.instrument === 'command' && options.name === 'request' && options.displayName === 'fetch' &&
        options.state !== 'pending'
      ) {
        let statusCode;

        statusCode = options.consoleProps["Response Status Code"];

        const isSuccess = statusCode && (statusCode + '')[0] === '2';
        const severity = isSuccess ? CONSTANTS.SEVERITY.SUCCESS : CONSTANTS.SEVERITY.WARNING;
        let log = formatFetch(options);

        if (options.consoleProps.Duration) {
          log += ` (${formatDuration(options.consoleProps.Duration)})`;
        }
        if (statusCode) {
          log += `\nStatus: ${statusCode}`;
        }
        if (options.err && options.err.message) {
          log += ' - ' + options.err.message;
        }

        if (
          !isSuccess &&
          options.consoleProps["Response Body"]
        ) {
          log += `\nResponse body: ${await this.format.formatXhrBody(options.consoleProps["Response Body"])}`;
        }

        this.collectorState.updateLog(log, severity, options.id);
      }
    });
  }

}
