const LOG_TYPE = require('../constants').LOG_TYPES;
const CONSTANTS = require('../constants');
const LogFormat = require("./LogFormat");

module.exports = class LogCollectCypressRoute {

  constructor(collectorState, config) {
    this.config = config;
    this.collectorState = collectorState;

    this.format = new LogFormat(config);
  }

  register() {
    Cypress.Commands.overwrite('server', (originalFn, options = {}) => {
      const prevCallback = options && options.onAnyResponse;
      options.onAnyResponse = async (route, xhr) => {
        if (prevCallback) {
          prevCallback(route, xhr);
        }

        if (!route) {
          return;
        }

        const severity = String(xhr.status).match(/^2[0-9]+$/) ? '' : CONSTANTS.SEVERITY.WARNING;
        let logMessage = `(${route.alias}) ${xhr.method} ${xhr.url}\n`;
        logMessage += this.format.formatXhrLog({
          request: {
            headers: await this.format.formatXhrBody(xhr.request.headers),
            body: await this.format.formatXhrBody(xhr.request.body),
          },
          response: {
            status: xhr.status,
            headers: await this.format.formatXhrBody(xhr.response.headers),
            body: await this.format.formatXhrBody(xhr.response.body),
          },
        });

        this.collectorState.addLog([LOG_TYPE.CYPRESS_ROUTE, logMessage, severity], null, xhr.id);
      };
      originalFn(options);
    });
  }

}
