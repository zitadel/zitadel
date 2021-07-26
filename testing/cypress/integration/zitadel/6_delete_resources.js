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


describe('CLEANUP: delete User', () => {
    it('CLEANUP: delete User', () => {
        //click on org to clear screen
        cy.visit(Cypress.env('consoleUrl') + '/org')
        cy.wait(1000)
        cy.visit(Cypress.env('consoleUrl') + '/users/list/humans')
        cy.url().should('contain', 'users/list/humans')
        cy.wait(10000)
        //force due to angular hidden buttons
        cy.get('tr').filter(':contains("demofirst")').find('button', { timeout: 30000 }).click({force: true})
        cy.get('button').filter(':contains("Delete")').click()
    })
})

describe('MACHINES: delete Machine', () => {
    it('MACHINES: delete Machine', () => {
        //click on org to clear screen
        cy.visit(Cypress.env('consoleUrl') + '/org')
        cy.wait(1000)
        cy.visit(Cypress.env('consoleUrl') + '/users/list/machines')
        cy.url().should('contain', 'users/list/machines')
        cy.wait(10000)
        //force due to angular hidden buttons
        cy.get('tr').filter(':contains("demomachineusername")').find('button', { timeout: 30000 }).click({force: true})
        cy.get('button').filter(':contains("Delete")').click()
    })
})

describe('CLEANUP: delete Project ', () => {
    it('CLEANUP: delete Project ', () => {
        cy.log(`PROJECT: delete project`);
        //click on org to clear screen
        cy.visit(Cypress.env('consoleUrl') + '/org')
        //click on Projects 
        cy.visit(Cypress.env('consoleUrl') + '/projects')
        cy.url().should('contain', '/projects')
        cy.wait(10000)
        //TODO variable for regex
        cy.get('.card').filter(':contains("newProjectToTest")', { timeout: 30000 }).find('button.delete-button').click()
        cy.get('button').filter(':contains("Delete")').click()
    })
})

