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

describe('MACHINES: show Machines ', () => {
    it('MACHINES: show Machines ', () => {
        cy.visit(Cypress.env('consoleUrl') + '/users/list/machines')
        cy.url().should('contain', 'users/list/machines')
    })
})

describe('MACHINES: add Machine', () => {
    it('MACHINES: add Machine', () => {
        //click on org to clear screen
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/users/list/machines')
        cy.url().should('contain', 'users/list/machines')
        cy.visit(Cypress.env('consoleUrl') + '/users/create-machine')
        cy.url().should('contain', 'users/create-machine')
        //force needed due to the prefilled username prefix
        cy.get('[formcontrolname^=userName]').type(Cypress.env('newMachineUserName'),{force: true})
        cy.get('[formcontrolname^=name]').type(Cypress.env('newMachineName'))
        cy.get('[formcontrolname^=description]').type(Cypress.env('newMachineDesription'))
        cy.get('button').filter(':contains("Create")').should('be.visible').click().then(() => {
            cy.wait(3000)
            cy.visit(Cypress.env('consoleUrl') + '/users/list/machines');
            cy.get('tr', { timeout: 30000 }).should('contain.text', "machineusername").and('exist');
        })
    })
})

