// NEEDS TO BE DISABLED!!!!!! this is just for testing
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
        return false
    }
})
// ###############################

describe('login username password', () => {

    it('should log a user in with username and password', ()=> {
        cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    })    

    it('should show personal information', () => {        
        //click on user information 
        cy.get('a[href*="users/me"').eq(0).click()
        cy.url().should('contain', '/users/me')
    })
})


