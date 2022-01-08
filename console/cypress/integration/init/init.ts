//import { apiAuth } from "../../support/api/apiauth";
//import { ensureMachineUserExists, ensureUserDoesntExist } from "../../support/api/users";
import cypress = require("cypress");
import { login as commonLogin, User } from "../../support/login/users";

describe('initialize organisation', () => {
    it('initializes', () => {
        const consoleUrl: string = Cypress.env('consoleUrl') 

        // Wait until ng serve is really ready
        cy.wait(30_000)

        login(User.IAMAdminUser, 'Password1!')

        // Create org
        cy.visit(`${consoleUrl}/org/create`)
        cy.contains('Use your personal account as organisation owner').click({ force: true })
        cy.get('[formcontrolname="name"]').type('caos-demo', { force: true, })
        cy.contains('button', 'CREATE').click({ force: true })
        cy.contains('button', 'Global').click({ force: true })
        cy.contains('button', 'caos-demo').click({ force: true })
  
        // Create sa
        cy.visit(`${consoleUrl}/users/create-machine`)
        cy.get('[formcontrolname="userName"]').type("e2e", { force: true })
        cy.get('[formcontrolname="name"]').type("e2e", { force: true })
        cy.get('[formcontrolname="description"]').type("User who calls the ZITADEL API for preparing end-to-end tests")
        cy.contains('button', 'Create').click({ force: true })


        addOrganisationRole('ORG_OWNER')

        // Create and download sa key
        cy.contains('div .card', 'Keys').contains('a', 'New').click({ force: true })     
        cy.contains('button', 'Add').click({ force: true })
        cy.contains('button', 'Download').click({ force: true })

/*        
        // Create e2e users
        // tmp
        cy.visit(`${consoleUrl}/users/me`)
        cy.contains('button', 'Global').click({ force: true })
        cy.contains('button', 'caos-demo').click({ force: true })
        

        //tmp
        cy.visit(`${consoleUrl}/users/list/machines`)
        cy.contains('tr', 'e2e').click({ force: true })
*/


        ;[{
            user: User.OrgOwner,
            role: 'ORG_OWNER'
        }, {
            user: User.OrgOwnerViewer,
            role: 'ORG_OWNER_VIEWER'
        }, {
            user: User.OrgProjectCreator,
            role: 'ORG_PROJECT_CREATOR'
        }, {
            user: User.LoginPolicyUser,
            role: null
        }, { 
            user: User.PasswordComplexityUser,
            role: null
        }].forEach(user => {
            login(User.IAMAdminUser, 'Password1!', true)
            cy.visit(`${consoleUrl}/users/create`)
            cy.contains('button', 'Global').click({ force: true })
            cy.contains('button', 'caos-demo').click({ force: true })
            cy.visit(`${consoleUrl}/users/create`)
            cy.get('[formcontrolname="email"]').type(`${user.user}@dummy.com`, { force: true })
            cy.get('[formcontrolname="userName"]').type(`${user.user}_user_name`, { force: true })
            cy.get('[formcontrolname="firstName"]').type(`${user.user}_first_name`, { force: true })
            cy.get('[formcontrolname="lastName"]').type(`${user.user}_last_name`, { force: true })
            cy.contains('label', 'Email Verified').click({ force: true })
            cy.contains('label', 'Set Initial Password').click({ force: true })
            const pw = Cypress.env(`${user.user}_password`)
            cy.get('[formcontrolname="password"]').type(pw, { force: true })
            cy.get('[formcontrolname="confirmPassword"]').type(pw, { force: true })

            cy.contains('button', 'Create').click({ force: true })
            cy.contains('h1', `${user.user}_first_name ${user.user}_last_name`)
            if (user.role) {
                addOrganisationRole(user.role)
            }
            login(user.user, pw)
        })
    })
})

function login(user: User, pw: string, relogin?: boolean) {
    commonLogin(user, false, pw, null, null, () => {
        if (relogin) {
            return
        }
        // Skip MFA and change password
        cy.contains('button', 'skip').click({ force: true })
        cy.get('#change-old-password').type(pw, { force: true })
        cy.get('#change-new-password').type(pw, { force: true })
        cy.get('#change-password-confirmation').type(pw, { force: true })
        cy.contains('button', 'next').click({ force: true })
        cy.contains('button', 'next').click({ force: true })
    })
}

function addOrganisationRole(role: string){
    cy.get('i.la-arrow-right').click({ force: true })
    cy.get('button[aria-label="add membership"]').click({ force: true })
    cy.contains('cnsl-label', 'Creation Type').next().click({ force: true })
    cy.contains('mat-option', 'Organisation').click({ force: true })
    cy.contains('cnsl-label', 'Role Name').next().click({ force: true })
    cy.contains('mat-option', role).click({ force: true })
    cy.contains('button', 'Add').click({ force: true })
}