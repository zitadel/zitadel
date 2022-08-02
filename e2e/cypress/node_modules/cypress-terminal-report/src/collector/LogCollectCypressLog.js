const LOG_TYPE = require('../constants').LOG_TYPES;

module.exports = class LogCollectCypressLog {

  constructor(collectorState, config) {
    this.config = config;
    this.collectorState = collectorState;
  }

  register() {
    Cypress.Commands.overwrite('log', (subject, ...args) => {
      this.collectorState.addLog([LOG_TYPE.CYPRESS_LOG, args.join(' ')]);
      subject(...args);
    });
  }

}
