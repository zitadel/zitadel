describe('/verify', () => {
    it('redirects after successful email verification', () => {
        cy.visit("/verify?userID=123&code=abc&submit=true")
        cy.location('pathname').should('eq', '/username')
    })
    it('shows an error if validation failed', () => {
        cy.visit("/verify?userID=123&code=abc&submit=true")
        cy.contains('error validating code')
    })
})