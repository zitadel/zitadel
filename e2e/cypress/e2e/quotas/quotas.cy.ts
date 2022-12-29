describe("quotas", () => {
    beforeEach(() => {
        cy.task('runSQL', "TRUNCATE logstore.access;")
    })
    it('passes', () => {

    })
})