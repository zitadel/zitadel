const LOG_TYPE = require('../constants').LOG_TYPES;
const CONSTANTS = require('../constants');

module.exports = class LogCollectCypressCommand {

  constructor(collectorState, config) {
    this.config = config;
    this.collectorState = collectorState;
  }

  register() {
    Cypress.on('log:added', (options) => {
      if (
        options.instrument === 'command' &&
        options.consoleProps &&
        !['xhr', 'log', 'request'].includes(options.name) &&
        !(options.name === 'task' && options.message.match(/ctrLogMessages/))
      ) {
        const log = options.name + '\t' + options.message;
        const severity = options.state === 'failed' ? CONSTANTS.SEVERITY.ERROR : '';
        this.collectorState.addLog([LOG_TYPE.CYPRESS_COMMAND, log, severity], options.id);
      }
    });
  }

}
