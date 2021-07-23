// NEEDS TO BE DISABLED!!!!!! this is just for testing
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
        return false
    }
})
// ###############################


it('LOGIN: Fill in credentials and login', () => {
    //console login
    cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    //wait for console to load
    cy.wait(5000)
})

describe('USER: show personal information', () => {
    // user interaction
    it('USER: show personal information', () => {
        cy.log(`USER: show personal information`);
        //click on user information 
        cy.get('a[href*="users/me"').eq(0).click()
        cy.url().should('contain', '/users/me')
    })
})


describe('USERS: show Users ', () => {
    it('PROJECT: show Projects ', () => {
        cy.visit('https://console.zitadel.ch/users/list/humans')
        cy.url().should('contain', 'users/list/humans')
    })
})

describe('USERS: add User', () => {
    it('USERS: add User', () => {
        //click on org to clear screen
        cy.visit('https://console.zitadel.ch/org')
        cy.wait(1000)
        cy.visit('https://console.zitadel.ch/users/list/humans')
        cy.url().should('contain', 'users/list/humans')
        cy.visit('https://console.zitadel.ch/users/create')
        cy.url().should('contain', 'users/create')
        cy.get('[formcontrolname^=email]').type(Cypress.env('newEmail'))
        //force needed due to the prefilled username prefix
        cy.get('[formcontrolname^=userName]').type(Cypress.env('newUserName'),{force: true})
        cy.get('[formcontrolname^=firstName]').type(Cypress.env('newFirstName'))
        cy.get('[formcontrolname^=lastName]').type(Cypress.env('newLastName'))
        cy.get('[formcontrolname^=phone]').type(Cypress.env('newPhonenumber'))
        cy.get('button').filter(':contains("Create")').should('be.visible').click()

    })
})

describe('USERS: delete User', () => {
    it('USERS: delete User', () => {
        //click on org to clear screen
        cy.visit('https://console.zitadel.ch/org')
        cy.wait(1000)
        cy.visit('https://console.zitadel.ch/users/list/humans')
        cy.url().should('contain', 'users/list/humans')
        cy.wait(6000)
        //force due to angular hidden buttons
        cy.get('tr').filter(':contains("demofirst")').find('button').click({force: true})
        cy.get('button').filter(':contains("Delete")').click()
    })
})
