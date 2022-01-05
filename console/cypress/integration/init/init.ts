//import { apiAuth } from "../../support/api/apiauth";
//import { ensureMachineUserExists, ensureUserDoesntExist } from "../../support/api/users";
import cypress = require("cypress");
import { login, User } from "../../support/login/users";


describe('initialize organisation', () => {
    it('initializes', () => {
        const adminPw = 'Password1!'
        const consoleUrl: string = Cypress.env('consoleUrl') 
 
        login(User.IAMAdminUser, false, adminPw, null, null, () => {
            // Login as zitadel admin for the first time
            cy.contains('button', 'skip').click()
            cy.get('#change-old-password').type(adminPw)
            cy.get('#change-new-password').type(adminPw)
            cy.get('#change-password-confirmation').type(adminPw)
            cy.contains('button', 'next').click()
            cy.contains('button', 'next').click()
        })

        // Create org
        cy.visit(`${consoleUrl}/org/create`)
        cy.contains('Use your personal account as organisation owner').click({ force: true })
        cy.get('[formcontrolname="name"]').type('caos-demo')
        cy.contains('button', 'CREATE').click()
        cy.contains('button', 'Global').click()
        cy.contains('button', 'caos-demo').click()
  
        // Create sa
        cy.visit(`${consoleUrl}/users/create-machine`)
        cy.get('[formcontrolname="userName"]').type("e2e")
        cy.get('[formcontrolname="name"]').type("e2e")
        cy.get('[formcontrolname="description"]').type("User who calls the ZITADEL API for preparing end-to-end tests")
        cy.contains('button', 'Create').click()

        // Create and download sa key
        cy.contains('div .card', 'Keys').contains('a', 'New').click()     
        cy.contains('button', 'Add').click()
        cy.contains('button', 'Download').click()

        // Create e2e users
        // tmp
/*        cy.visit(`${consoleUrl}/users/me`)
        cy.contains('button', 'Global').click()
        cy.contains('button', 'caos-demo').click()*/

        ;[User.OrgOwner,  User.OrgOwnerViewer, User.OrgProjectCreator, User.LoginPolicyUser, User.PasswordComplexityUser].forEach(user => {
            cy.visit(`${consoleUrl}/users/create`)
            cy.get('[formcontrolname="email"]').type("dummy@example.com")
            cy.get('[formcontrolname="userName"]').type(`${user}_user_name`, { force: true })
            cy.get('[formcontrolname="firstName"]').type(`${user}_first_name`)
            cy.get('[formcontrolname="lastName"]').type(`${user}_last_name`)
            cy.get('mat-checkbox').click()
            cy.contains('button', 'Create').click()
            cy.url().should(url => {
                expect(url).to.match(/http\:\/\/localhost\:4200\/users\/[0-9]+/)
            })
            cy.url().then(url => {
                cy.visit(`${url}/password`)
                const pw = Cypress.env(`${user}_password`)
                cy.get('[formcontrolname="password"]').type(pw)
                cy.get('[formcontrolname="confirmPassword"]').type(pw)
                cy.contains('button', 'Set New Password').click()
            })
        })
    })
})