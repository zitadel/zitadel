module.exports = (on, config) => {
  // modify the config values
  config.defaultCommandTimeout = 60000

  //config.env.consoleUrl = "https://console.zitadel.ch"

  config.env.newEmail = "demo@caos.ch"
  config.env.newUserName = "demo"
  config.env.newFirstName = "demofirstname"
  config.env.newLastName = "demolastname"
  config.env.newPhonenumber = "+41 123456789"

  config.env.newMachineUserName = "machineusername"
  config.env.newMachineName = "name"
  config.env.newMachineDesription = "description"


  // IMPORTANT return the updated config object
  return config

}
