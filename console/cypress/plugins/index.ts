import { readFileSync } from "fs";

module.exports = (on, config) => {

  require('cypress-terminal-report/src/installLogsPrinter')(on);

  config.defaultCommandTimeout = 10_000

  config.env.parsedServiceAccountKey = config.env.serviceAccountKey
  if (config.env.serviceAccountKeyPath) {
    config.env.parsedServiceAccountKey = JSON.parse(readFileSync(config.env.serviceAccountKeyPath, 'utf-8'))
  }

  return config
}
