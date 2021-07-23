module.exports = (on, config) => {
  // modify the config values
  config.defaultCommandTimeout = 10000

  config.env.newEmail = "demo@caos.ch"
  config.env.newUserName = "demo"
  config.env.newFirstName = "demofirstname"
  config.env.newLastName = "demolastname"
  config.env.newPhonenumber = "+41 123456789"




  // IMPORTANT return the updated config object
  return config

}
