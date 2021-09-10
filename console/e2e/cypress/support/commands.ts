
declare namespace Cypress {
    interface Chainable {
        /**
         * Custom command that authenticates a user.
         * 
         * @example cy.consolelogin('hodor', 'hodor1234', https://console.zitadel.ch })
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
    })