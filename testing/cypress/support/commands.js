Cypress.Commands.add('consolelogin',
    (username, password, consoleUrl) => {

        cy.visit(consoleUrl)
        // fill the fields and push button
        cy.get('#loginName').type(username, { log: false })
        cy.get('#submit-button').click()
        cy.get('#password').type(password, { log: false })
        cy.get('#submit-button').click()
        //TODO: MFA

    })