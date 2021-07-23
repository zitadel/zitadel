module.exports = (on, config) => {
  // modify the config values
  config.defaultCommandTimeout = 10000

  // IMPORTANT return the updated config object
  return config

}
