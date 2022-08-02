module.exports = class CtrError extends Error {

  constructor(message) {
    super(`cypress-terminal-report: ${message}`);
  }

};
