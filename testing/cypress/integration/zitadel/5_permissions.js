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


describe('PERMISSIONS: show Projects ', () => {
    it('PROJECT: show Projects ', () => {
        cy.visit(Cypress.env('consoleUrl') + '/projects')
        cy.url().should('contain', '/projects')
    })
})

describe('PERMISSIONS: add Role ', () => {

    it('PERMISSIONS: add Role ', () => {
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url({ timeout: 30000 }).should('contain', '/projects');
            cy.get('.card', { timeout: 30000 }).should('contain.text', "newProjectToTest")
        })
        cy.get('.card').filter(':contains("newProjectToTest")', { timeout: 30000 }).click()
        cy.get('.app-container').filter(':contains("newAppToTest")').should('be.visible').click()
        let projectID
        cy.url().then(url => {
            cy.log(url.split('/')[4])
            projectID = url.split('/')[4]
          });
        
        cy.then(() => cy.visit(Cypress.env('consoleUrl') + '/projects/' + projectID +'/roles/create'))
        cy.get('[formcontrolname^=key]').type("newdemorole")
        cy.get('[formcontrolname^=displayName]').type("newdemodisplayname")
        cy.get('[formcontrolname^=group]').type("newdemogroupname")
        cy.get('button').filter(':contains("Save")').should('be.visible').click()
        //let the Role get processed
        cy.wait(5000)
    })
})

describe('PERMISSIONS: add Grant ', () => {

    it('PERMISSIONS: add Grant ', () => {
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url({ timeout: 30000 }).should('contain', '/projects');
            cy.get('.card', { timeout: 30000 }).should('contain.text', "newProjectToTest")
        })
        cy.get('.card').filter(':contains("newProjectToTest")', { timeout: 30000 }).click()
        cy.get('.app-container').filter(':contains("newAppToTest")').should('be.visible').click()
        let projectID
        cy.url().then(url => {
            cy.log(url.split('/')[4])
            projectID = url.split('/')[4]
          });
        
        cy.then(() => cy.visit(Cypress.env('consoleUrl') + '/grant-create/project/' + projectID ))
        cy.get('input').type("demo")
        cy.get('[role^=listbox]').filter(':contains("demo@caos-demo.zitadel.ch")' ).should('be.visible').click()
        cy.wait(5000)
        //cy.get('.button').contains('Continue').click()
        cy.get('button').filter(':contains("Continue")', { timeout: 30000 }).click()
        cy.wait(5000)
        cy.get('tr').filter(':contains("demo")').find('label').click()
        cy.get('button').filter(':contains("Save")').should('be.visible').click()
        //let the grant get processed
        cy.wait(5000)
    })
})


