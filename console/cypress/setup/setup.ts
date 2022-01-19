//import { apiAuth } from "../../support/api/apiauth";
//import { ensureMachineUserExists, ensureUserDoesntExist } from "../../support/api/users";
import { login as commonLogin, User } from "../support/login/users";

describe('setup e2e test data in ZITADEL', () => {
    const consoleUrl: string = Cypress.env('consoleUrl') 
    it('wait for 30 seconds until ng serve is REALLY ready', () => {
        cy.wait(30_000)
    })

    it('logs the admin user in', () => {
        login(User.IAMAdminUser, 'Password1!')
    })

    it('creates the caos-demo org', () => {
        cy.visit(`${consoleUrl}/org/create`)
        cy.contains('Use your personal account as organisation owner').click({ force: true })
        cy.get('[formcontrolname="name"]').type('caos-demo', { force: true, })
        cy.contains('button', 'CREATE').click({ force: true })
        cy.contains('button', 'Global').click({ force: true })
        cy.contains('button', 'caos-demo').click({ force: true })
    })

    it('creates the service account used by the test suite to call the ZITADEL API', () => {
        cy.visit(`${consoleUrl}/users/create-machine`)
        cy.get('[formcontrolname="userName"]').type("e2e", { force: true })
        cy.get('[formcontrolname="name"]').type("e2e", { force: true })
        cy.get('[formcontrolname="description"]').type("User who calls the ZITADEL API for preparing end-to-end tests")
        cy.contains('button', 'Create').click({ force: true })
    })

    it('gives the service account to the ORG_OWNER role', () => {
        addOrganisationRole('ORG_OWNER')
    })

    it('creates and downloads the service account key', () => {
        cy.contains('div .card', 'Keys').contains('a', 'New').click({ force: true })     
        cy.contains('button', 'Add').click({ force: true })
        cy.contains('button', 'Download').click({ force: true })
    })

    it('enables all features for the demo-org', () => {
        cy.visit(`${consoleUrl}/org/features`)
        cy.get('label.mat-slide-toggle-label').click({ force: true, multiple: true })
        cy.contains('button', 'Save').click({ force: true })
    })

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

        it(`creates the test user ${user.user} and gives it the role ${user.role}`, () => {
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