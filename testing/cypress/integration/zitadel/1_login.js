// NEEDS TO BE DISABLED!!!!!! this is just for testing
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
        return false
    }
})
// ###############################

describe('LOGIN: Basic User/PW console.zitadel.ch', () => {
    it('LOGIN: Fill in credentials and login', () => {


        //console login
        cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
        //wait for console to load
        cy.wait(5000)
    })
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


