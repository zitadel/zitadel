// NEEDS TO BE DISABLED!!!!!! this is just for testing
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
        return false
    }
})
// ###############################

describe("projects", ()=> {

    before(()=> {
        cy.consolelogin(Cypress.env('username'), Cypress.env('password'))
    })

    it('should show projects', () => {
        cy.visit(Cypress.env('consoleUrl') + '/projects')
        cy.url().should('contain', '/projects')
    })

    describe('add', () => {

        before('cleanup', () => {
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

        it('should add a project', () => {
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

        it('should create an app', () => {
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
})
