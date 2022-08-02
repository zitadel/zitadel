const LOG_TYPE = require('../constants').LOG_TYPES;
const HTTP_METHODS = require('../constants').HTTP_METHODS;
const CONSTANTS = require('../constants');
const LogFormat = require("./LogFormat");

Object.defineProperty(RegExp.prototype, "toJSON", {
  value: RegExp.prototype.toString
});

module.exports = class LogCollectCypressIntercept {

  constructor(collectorState, config) {
    this.config = config;
    this.collectorState = collectorState;

    this.format = new LogFormat(config);
  }

  register() {
    Cypress.Commands.overwrite('intercept', (originalFn, ...args) => {
      let message = '';

      if (typeof args[0] === "string" && HTTP_METHODS.includes(args[0].toUpperCase())) {
        message += `Method: ${args[0]}\nMatcher: ${JSON.stringify(args[1])}`;
        if (args[2]) {
          message += `\nMocked Response: ${typeof args[2] === 'object' ? JSON.stringify(args[2]) : args[2]}`;
        }
      } else {
        message += `Matcher: ${JSON.stringify(args[0])}`;
        if (args[1]) {
          message += `\nMocked Response: ${typeof args[1] === 'object' ? JSON.stringify(args[1]) : args[1]}`;
        }
      }

      this.collectorState.addLog([LOG_TYPE.CYPRESS_INTERCEPT, message, CONSTANTS.SEVERITY.SUCCESS]);
      return originalFn(...args);
    });
  }

}
