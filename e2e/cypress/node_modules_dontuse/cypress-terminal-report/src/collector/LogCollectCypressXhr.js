const LOG_TYPE = require('../constants').LOG_TYPES;
const CONSTANTS = require('../constants');
const LogFormat = require("./LogFormat");

module.exports = class LogCollectCypressXhr {

  constructor(collectorState, config) {
    this.config = config;
    this.collectorState = collectorState;

    this.format = new LogFormat(config);
  }

  register() {
    const formatXhr = (options) => options.message +
      (options.consoleProps.Stubbed === 'Yes' ? 'STUBBED ' : '') +
      options.consoleProps.Method + ' ' + options.consoleProps.URL;

    const formatDuration = (durationInMs) =>
      durationInMs < 1000 ? `${durationInMs} ms` : `${durationInMs / 1000} s`;

    Cypress.on('log:added', (options) => {
      if (
        options.instrument === 'command' &&
        options.consoleProps &&
        options.name === 'xhr' &&
        // Prevent duplicated xhr logs in case of cy.intercept.
        options.displayName !== 'req'
      ) {
        const log = formatXhr(options);
        const severity = options.state === 'failed' ? CONSTANTS.SEVERITY.WARNING : '';
        this.collectorState.addLog([LOG_TYPE.CYPRESS_XHR, log, severity], options.id);
      }
    });

    Cypress.on('log:changed', async (options) => {
      if (
        options.instrument === 'command' &&
        ['request', 'xhr'].includes(options.name) &&
        options.consoleProps &&
        options.state !== 'pending'
      ) {
        let statusCode, statusText;

        if (!options.consoleProps.XHR) {
          [, statusCode, statusText] = /^(\d{3})\s\((.+)\)$/.exec(options.consoleProps.Status) || [];
        } else {
          statusCode = options.consoleProps.XHR.status;
          statusText = options.consoleProps.XHR.statusText;
        }

        const isSuccess = statusCode && (statusCode + '')[0] === '2';
        const severity = isSuccess ? CONSTANTS.SEVERITY.SUCCESS : CONSTANTS.SEVERITY.WARNING;
        let log = formatXhr(options);

        if (options.consoleProps.Duration) {
          log += ` (${formatDuration(options.consoleProps.Duration)})`;
        }
        if (statusCode && statusText) {
          log += `\nStatus: ${statusCode} - ${statusText}`;
        }
        if (options.err && options.err.message.match(/abort/)) {
          log += ' - ABORTED';
        }
        if (
          !isSuccess &&
          options.consoleProps.XHR &&
          options.consoleProps.XHR.response &&
          options.consoleProps.XHR.response.size &&
          !this.collectorState.hasXhrResponseBeenLogged(options.consoleProps.XHR.id)
        ) {
          log += `\nResponse body: ${await this.format.formatXhrBody(options.consoleProps.XHR.response)}`;
        }

        this.collectorState.updateLog(log, severity, options.id);
      }
    });
  }

}
