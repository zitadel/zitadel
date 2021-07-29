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


describe('PROJECT: show Projects ', () => {
    it('PROJECT: show Projects ', () => {
        cy.visit(Cypress.env('consoleUrl') + '/projects')
        cy.url().should('contain', '/projects')
    })
})

describe('PROJECT: add Project ', () => {

    it('PROJECT: add Project ', () => {
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url().should('contain', '/projects');
            cy.get('.add-project-button')
        })
        cy.get('.add-project-button').click({ force: true })
        cy.get('input').type("newProjectToTest")
        cy.get('[type^=submit]').click().then(() => {
            cy.get('h1').should('contain', "Project newProjectToTest")
        })
    })
})

describe('PROJECT: create app in Project ', () => {

    it('PROJECT: create app ', () => {
        //click on org to clear screen
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url().should('contain', '/projects');
            cy.get('.card').should('contain.text', "newProjectToTest")
        })
        cy.get('.card').contains("newProjectToTest").click()
        cy.get('.cnsl-app-card').filter(':contains("add")').click()
        cy.get('[formcontrolname^=name]').type("newAppToTest")
        // select webapp
        cy.get('[for^=WEB]').click()
        cy.get('[type^=submit]').filter(':contains("Continue")').should('be.visible').eq(0).click()
        //select authentication
        cy.get('[for^=PKCE]').click()
        cy.get('[type^=submit]').filter(':contains("Continue")').should('be.visible').eq(1).click()
        //enter URL
        cy.get('cnsl-redirect-uris').eq(0).type("https://testurl.org")
        cy.get('cnsl-redirect-uris').eq(1).type("https://testlogouturl.org")
        cy.get('[type^=submit]').filter(':contains("Continue")').should('be.visible').eq(2).click()
        cy.get('button').filter(':contains("Create")').should('be.visible').click().then(() => {
            cy.get('[id*=overlay]').should('exist')
        })
        //TODO: check client ID/Secret
        cy.get('button').filter(':contains("Close")').should('exist').click()
    })
})

