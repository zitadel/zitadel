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
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        // cy.visit(Cypress.env('consoleUrl') + '/users/list/humans')
        // cy.url().should('contain', 'users/list/humans')
        cy.visit(Cypress.env('consoleUrl') + '/users/list/humans').then(() => {
            cy.url().should('contain', '/users/list/humans');
            cy.get('tr').should('contain.text', "demofirst")
        })
        
        //force due to angular hidden buttons
        cy.get('tr').filter(':contains("demofirst")').find('button').click({force: true}).then(() => {
            cy.get('button').should('contain', 'Delete');
        })
        cy.get('button').filter(':contains("Delete")').click().then(() => {
            cy.wait(3000)
            cy.visit(Cypress.env('consoleUrl') + '/users/list/humans');
            cy.get('[text*=demofirst]').should('not.exist');
            })
    })
})

describe('CLEANUP: delete Machine', () => {
    it('CLEANUP: delete Machine', () => {
        //click on org to clear screen
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/users/list/machines').then(() => {
            cy.url().should('contain', '/users/list/machines');
            cy.get('tr').should('contain.text', "machineusername")
        })
        
        //force due to angular hidden buttons
        cy.get('tr').filter(':contains("machineusername")').find('button').click({force: true}).then(() => {
            cy.get('button').should('contain', 'Delete');
        })
        cy.get('button').filter(':contains("Delete")').click().then(() => {
            cy.wait(3000)
            cy.visit(Cypress.env('consoleUrl') + '/users/list/machines');
            cy.get('[text*=machineusername]').should('not.exist');
        })
    })
})

describe('CLEANUP: delete Project ', () => {
    it('CLEANUP: delete Project ', () => {
        cy.log(`PROJECT: delete project`);
        //click on org to clear screen
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        //click on Projects 
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url().should('contain', '/projects');
            cy.get('.card').should('contain.text', "newProjectToTest")
        })
        //TODO variable for regex
        cy.get('.card').filter(':contains("newProjectToTest")').find('button.delete-button').click()
        cy.get('button').filter(':contains("Delete")').click().then(() => {
            cy.wait(2000)
            cy.visit(Cypress.env('consoleUrl') + '/projects');
            cy.get('.card').contains("newProjectToTest").should('not.exist');
        })
    })
})

