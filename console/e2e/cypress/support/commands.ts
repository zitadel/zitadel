
declare namespace Cypress {
    interface Chainable {
        /**
         * Overwritten command that unlike the original cy.exec throws an error including the commands full stdout and stderr.
         * see https://github.com/cypress-io/cypress/issues/5470#issuecomment-569627930
         * 
         * @example cy.execZero('whoami', result => { expect(result.stdout).to.be('hodor') })
         */
            consolelogin(username: string, password: string, consoleUrl: string): void
    }
}

Cypress.Commands.add('consolelogin',
    (username: string, password: string, consoleUrl: string) => {

        cy.visit(consoleUrl)
        // fill the fields and push button
        cy.get('#loginName').type(username, { log: false })
        cy.get('#submit-button').click()
        cy.get('#password').type(password, { log: false })
        cy.get('#submit-button').click()
        cy.location('pathname', {timeout: 5 * 1000}).should('eq', '/');
        //TODO: MFA
        // skip MFA for now

    })