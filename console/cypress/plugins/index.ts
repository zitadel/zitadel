module.exports = (on, config) => {

  require('cypress-terminal-report/src/installLogsPrinter')(on);

  config.defaultCommandTimeout = 10_000

  return config
}
