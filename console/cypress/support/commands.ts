/*
namespace Cypress {
    interface Chainable {
*/
        /**
         * Custom command that authenticates a user.
         *
         * @example cy.consolelogin('hodor', 'hodor1234')
         */
/*        consolelogin(username: string, password: string): void            
    }
}

Cypress.Commands.add('consolelogin', { prevSubject: false }, (username: string, password: string) => {

    window.sessionStorage.removeItem("zitadel:access_token")
    cy.visit(Cypress.env('consoleUrl')).then(() => {
        // fill the fields and push button
        cy.get('#loginName').type(username, { log: false })
        cy.get('#submit-button').click()
        cy.get('#password').type(password, { log: false })
        cy.get('#submit-button').click()
        cy.location('pathname', {timeout: 5 * 1000}).should('eq', '/');
    })
})
*/