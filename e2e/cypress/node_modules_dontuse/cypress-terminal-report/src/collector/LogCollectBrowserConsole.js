const LOG_TYPE = require('../constants').LOG_TYPES;
const CONSTANTS = require('../constants');

module.exports = class LogCollectBrowserConsole {

  constructor(collectorState, config) {
    this.config = config;
    this.collectorState = collectorState;
  }

  register() {
    const oldConsoleMethods = {};

    Cypress.on('window:before:load', () => {
      const docIframe = window.parent.document.querySelector("[id*='Your project: ']") ||
        window.parent.document.querySelector("[id*='Your App']");
      const appWindow = docIframe.contentWindow;

      const processArg = (arg) => {
        if (['string', 'number', 'undefined', 'function'].includes(typeof arg)) {
          return arg ? arg.toString() : arg === undefined ? 'undefined' : '';
        }

        if (
          (arg instanceof appWindow.Error || arg instanceof Error) &&
          typeof arg.stack === 'string'
        ) {
          let stack = arg.stack;
          if (stack.indexOf(arg.message) !== -1) {
            stack = stack.slice(stack.indexOf(arg.message) + arg.message.length + 1);
          }
          return arg.toString() + '\n' + stack;
        }

        let json = '';
        try {
          json = JSON.stringify(arg, null, 2);
        } catch (e) {
          return '[unprocessable=' + arg + ']';
        }

        if (typeof json === 'undefined') {
          return 'undefined';
        }

        return json;
      };

      const createWrapper = (method, logType, type = CONSTANTS.SEVERITY.SUCCESS) => {
        oldConsoleMethods[method] = appWindow.console[method];

        appWindow.console[method] = (...args) => {
          this.collectorState.addLog([logType, args.map(processArg).join(`,\n`), type]);
          oldConsoleMethods[method](...args);
        };
      };

      if (this.config.collectTypes.includes(LOG_TYPE.BROWSER_CONSOLE_WARN)) {
        createWrapper('warn', LOG_TYPE.BROWSER_CONSOLE_WARN, CONSTANTS.SEVERITY.WARNING);
      }
      if (this.config.collectTypes.includes(LOG_TYPE.BROWSER_CONSOLE_ERROR)) {
        createWrapper('error', LOG_TYPE.BROWSER_CONSOLE_ERROR, CONSTANTS.SEVERITY.ERROR);
      }
      if (this.config.collectTypes.includes(LOG_TYPE.BROWSER_CONSOLE_INFO)) {
        createWrapper('info', LOG_TYPE.BROWSER_CONSOLE_INFO);
      }
      if (this.config.collectTypes.includes(LOG_TYPE.BROWSER_CONSOLE_DEBUG)) {
        createWrapper('debug', LOG_TYPE.BROWSER_CONSOLE_DEBUG);
      }
      if (this.config.collectTypes.includes(LOG_TYPE.BROWSER_CONSOLE_LOG)) {
        createWrapper('log', LOG_TYPE.BROWSER_CONSOLE_LOG);
      }
    });
  }

}
