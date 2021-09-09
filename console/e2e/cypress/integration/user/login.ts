// NEEDS TO BE DISABLED!!!!!! this is just for testing
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
        return false
    }
})
// ###############################

describe('USER: show personal information', () => {

    before(()=> {
        cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    })
    
    it('USER: show personal information', () => {        
        //click on user information 
        cy.get('a[href*="users/me"').eq(0).click()
        cy.url().should('contain', '/users/me')
    })
})


