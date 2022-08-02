const LOG_TYPE = require('../constants').LOG_TYPES;
const HTTP_METHODS = require('../constants').HTTP_METHODS;
const CONSTANTS = require('../constants');
const LogFormat = require("./LogFormat");

module.exports = class LogCollectCypressRequest {

  constructor(collectorState, config) {
    this.config = config;
    this.collectorState = collectorState;

    this.format = new LogFormat(config);
  }

  register() {
    const isValidHttpMethod = (str) =>
      typeof str === 'string' && HTTP_METHODS.some((s) => str.toUpperCase().includes(s));

    const isNetworkError = (e) =>
      e.message && e.message.startsWith('`cy.request()` failed trying to load:');

    const isStatusCodeFailure = (e) => e.message && e.message.startsWith('`cy.request()` failed on:');

    const parseRequestStatusCodeFailureMessage = (message) => {
      const responseStart = '\n\nThe response we got was:\n\n';
      const statusStart = 'Status: ';
      const headersStart = '\nHeaders: ';
      const bodyStart = '\nBody: ';
      if (
        message.indexOf(responseStart) === -1 ||
        message.indexOf(statusStart) === -1 ||
        message.indexOf(headersStart) === -1 ||
        message.indexOf(bodyStart) === -1
      ) {
        return {status: 'Cannot parse cy.request status code failure message!'};
      }
      const response = message.substr(message.indexOf(responseStart) + responseStart.length);
      const statusStr = response.substr(
        response.indexOf(statusStart) + statusStart.length,
        response.indexOf(headersStart) - (response.indexOf(statusStart) + statusStart.length)
      );
      const headersStr = response.substr(
        response.indexOf(headersStart) + headersStart.length,
        response.indexOf(bodyStart) - (response.indexOf(headersStart) + headersStart.length)
      );
      const bodyStr = response.substr(response.indexOf(bodyStart) + bodyStart.length);
      return {status: statusStr, headers: headersStr, body: bodyStr.trimEnd()};
    };

    const parseRequestNetworkError = (message) => {
      const errorPartStart = 'We received this error at the network level:\n\n  > ';
      const errorPrefix = 'Error: ';
      if (message.indexOf(errorPartStart) === -1) {
        return {status: 'Cannot parse cy.request network error message!'};
      }
      let fromError = message.substr(message.indexOf(errorPartStart) + errorPartStart.length);
      let errorPart = fromError.substr(0, fromError.indexOf('\n'));
      if (errorPart.startsWith(errorPrefix)) {
        return errorPart.substr(errorPrefix.length).trim();
      }
      return errorPart.trim();
    };

    Cypress.Commands.overwrite('request', (originalFn, ...args) => {
      if (typeof args === 'object' && args !== null && args[0]['log'] === false){
        return originalFn(...args);
      }

      let log;
      let requestBody;
      let requestHeaders;

      // args can have up to 3 arguments
      // https://docs.cypress.io/api/commands/request.html#Syntax
      if (args[0].method) {
        log = `${args[0].method} ${args[0].url ? `${args[0].url}` : args[0]}`;
        requestBody = args[0].body;
        requestHeaders = args[0].headers;
      } else if (isValidHttpMethod(args[0])) {
        log = `${args[0]} ${args[1]}`;
        requestBody = args[3];
      } else {
        log = `${args[0]}`;
      }

      return Promise.all([
        this.format.formatXhrBody(requestHeaders),
        this.format.formatXhrBody(requestBody)
      ])
        .then(([formattedRequestHeaders, formattedRequestBody]) => {
          const requestData = {
            headers: formattedRequestHeaders,
            body: formattedRequestBody,
          };

          return originalFn(...args).catch(async (e) => {
            if (isNetworkError(e)) {
              log +=
                `\n` +
                this.format.formatXhrLog({
                  request: requestData,
                  networkError: parseRequestNetworkError(e.message),
                });
            } else if (isStatusCodeFailure(e)) {
              const xhr = parseRequestStatusCodeFailureMessage(e.message);
              log +=
                `\n` +
                this.format.formatXhrLog({
                  request: requestData,
                  response: {
                    status: xhr.status,
                    headers: await this.format.formatXhrBody(xhr.headers),
                    body: await this.format.formatXhrBody(xhr.body),
                  },
                });
            } else if (e.message.match(/timed out/)) {
              log += `\n` + 'Timed out!';
            } else {
              log += `\n` + 'Cannot parse cy.request error content!';
            }

            this.collectorState.addLog([LOG_TYPE.CYPRESS_REQUEST, log, CONSTANTS.SEVERITY.ERROR]);
            throw e;
          })
            .then(response => {
              return Promise.all([
                this.format.formatXhrBody(response.headers),
                this.format.formatXhrBody(response.body)
              ])
                .then(([formattedResponseHeaders, formattedResponseBody]) => {
                  log +=
                    `\n` +
                    this.format.formatXhrLog({
                      request: requestData,
                      response: {
                        status: response.status,
                        headers: formattedResponseHeaders,
                        body: formattedResponseBody,
                      },
                    });

                  this.collectorState.addLog([LOG_TYPE.CYPRESS_REQUEST, log]);
                  return response;
                });
            });
        });
    });
  }

}
