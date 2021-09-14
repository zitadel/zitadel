module.exports = (on, config) => {

  config.env.domain = "zitadel.ch"
  config.env.newEmail = "demo@caos.ch"
  config.env.newUserName = "demo"
  config.env.fullUserName = `demo@caos-demo.${config.env.domain}`
  config.env.newFirstName = "demofirstname"
  config.env.newLastName = "demolastname"
  config.env.newPhonenumber = "+41 123456789"

  config.env.newMachineUserName = "machineusername"
  config.env.newMachineName = "name"
  config.env.newMachineDesription = "description"

  // IMPORTANT return the updated config object
  return config

}
