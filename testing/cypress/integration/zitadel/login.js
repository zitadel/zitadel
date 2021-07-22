// NEEDS TO BE DISABLED!!!!!! this is just for testing
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
    return false
  }
})
// ###############################

describe('LOGIN: Basic User/PW console.zitadel.ch', { execTimeout: 9000 }, () => {
    // availability of zitadel
    //Environment Variables with CYPRESS_* are available automatically
    const username = Cypress.env('username')
    const password = Cypress.env('password')
    const consoleUrl = Cypress.env('consoleUrl')

    it('LOGIN: Fill in credentials and login', () => {
     cy.visit(consoleUrl)
     // fill the fields and push button
     cy.get('#loginName').type(username, {log: false}) 
     cy.get('#submit-button').click()
     cy.get('#password').type(password, {log: false})  
     cy.get('#submit-button').click()
     //TODO: MFA
    })
    it('LOGIN: Wait 5 seconds for console to load', () => {
        //wait for console to load into browser
        cy.wait(5000)
       })
    })

    describe('USER: show personal information', { execTimeout: 9000 }, () => {
        // user interaction
        it('USER: show personal information', () => {
        cy.log(`USER: show personal information`);
        //click on user information 
        cy.get('a[href*="users/me"').eq(0).click() 
        cy.url().should('contain', '/users/me')    
        })
        })

        describe('USER: show Projects ', { execTimeout: 9000 }, () => {
            it('USER: show Projects ', () => {
            cy.get('a[href*="projects"').eq(0).click()      
            cy.url().should('contain', '/projects')    
              
        })
            })
   
            describe('PROJECT: add Project ', { execTimeout: 9000 }, () => {

                it('PROJECT: add Project ', () => {
                cy.get('a[href*="projects"').eq(0).click()      
                cy.url().should('contain', '/projects')
                cy.get('.add-project-button', {timeout: 5000}).click() 
                cy.get('input').type("newProjectToTest") 
                cy.get('[type^=submit]', {timeout: 5000}).click() 

            })
                })
       
                describe('PROJECT: create app in Project ', { execTimeout: 9000 }, () => {

                    it('PROJECT: create app ', () => {
                    //click on org to clear screen
                    cy.get('a[href*="org"').eq(0).click()
                    cy.wait(1000)
                    cy.get('a[href*="projects"').eq(0).click()      
                    cy.url().should('contain', '/projects')
                    cy.wait(1000)
                    cy.get('.card').contains("newProjectToTest").click()
                    cy.get('.cnsl-app-card', {timeout: 5000}).filter(':contains("add")').click() 
                    cy.get('[formcontrolname^=name]').type("newAppToTest") 
                    // select webapp
                    cy.get('[for^=WEB]').click() 
                    cy.get('[type^=submit]', {timeout: 5000}).filter(':contains("Continue")').should('be.visible').eq(0).click() 
                    //select authentication
                    cy.get('[for^=PKCE]').click() 
                    cy.get('[type^=submit]', {timeout: 5000}).filter(':contains("Continue")').should('be.visible').eq(1).click() 
                    //enter URL
                    cy.get('cnsl-redirect-uris').eq(0).type("https://testurl.org")
                    cy.get('cnsl-redirect-uris').eq(1).type("https://testlogouturl.org")
                    cy.get('[type^=submit]', {timeout: 5000}).filter(':contains("Continue")').should('be.visible').eq(2).click() 
                    cy.get('button', {timeout: 5000}).filter(':contains("Create")').should('be.visible').click() 
                    //wait for application to be created
                    cy.wait(3000)
                    //TODO: check client ID/Secret
                    cy.get('button', {timeout: 5000}).filter(':contains("Close")').should('be.visible').click() 

                    
                })
                    })

                describe('PROJECT: delete Project ', { execTimeout: 9000 }, () => {
                    it('PROJECT: delete Project ', () => {
                    cy.log(`PROJECT: delete project`);
                    //click on org to clear screen
                    cy.get('a[href*="org"').eq(0).click()
                    //click on Projects 
                    cy.get('a[href*="projects"').eq(0).click()      
                    cy.url().should('contain', '/projects')
                    cy.wait(3000)
                    //TODO variable for regex
                    cy.get('.card').filter(':contains("newProjectToTest")').find('button.delete-button').click()
                    cy.get('button').filter(':contains("Delete")').click()

                })
                    })

