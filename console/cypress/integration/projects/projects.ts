import { apiAuth } from "../../support/api/apiauth";
import { ensureProjectDoesntExist, ensureProjectExists } from "../../support/api/projects";
import { login, User } from "../../support/login/users";

describe("projects", ()=> {

    const testProjectNameCreate = 'e2eprojectcreate'
    const testProjectNameDeleteList = 'e2eprojectdeletelist'
    const testProjectNameDeleteGrid = 'e2eprojectdeletegrid'

    ;[User.OrgOwner].forEach(user => {

        describe(`as user "${user}"`, () => {

            beforeEach(()=> {
                login(user)
            })

            describe('add project', () => {
                beforeEach(`ensure it doesn't exist already`, () => {
                    apiAuth().then(api => {
                        ensureProjectDoesntExist(api, testProjectNameCreate)
                    })
                    cy.visit(`${Cypress.env('consoleUrl')}/projects`)
                })

                it('should add a project', () => {
                    cy.get('.add-project-button').click({ force: true })
                    cy.get('input').type(testProjectNameCreate)
                    cy.get('[type^=submit]').click()
                    cy.get('.data-e2e-success')
                    cy.wait(200)
                    cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                })
            })

            describe('remove project', () => {

                describe('list view', () => {
                    beforeEach('ensure it exists', () => {
                        apiAuth().then(api => {
                            ensureProjectExists(api, testProjectNameDeleteList)
                        })
                        cy.visit(`${Cypress.env('consoleUrl')}/projects`)
                    })

                    it('removes the project', () => {
                        cy.get('[data-e2e=toggle-grid]').click()
                        cy.get('[data-cy=timestamp]')
                        cy.contains("tr", testProjectNameDeleteList, { timeout: 1000 })
                            .find('[data-e2e=delete-project-button]')
                            .click({force: true})
                        cy.get('[e2e-data="confirm-dialog-button"]').click()
                        cy.get('.data-e2e-success')
                        cy.wait(200)
                        cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                    })    
                })

                describe('grid view', () => {
                    beforeEach('ensure it exists', () => {
                        apiAuth().then(api => {
                            ensureProjectExists(api, testProjectNameDeleteGrid)
                        })
                        cy.visit(`${Cypress.env('consoleUrl')}/projects`)
                    })

                    it('removes the project', () => {
                        cy.contains('[data-e2e=grid-card]', testProjectNameDeleteGrid)
                            .find('[data-e2e=delete-project-button]')
                            .trigger('mouseover')
                            .click()
                        cy.get('[e2e-data="confirm-dialog-button"]').click()
                        cy.get('.data-e2e-success')
                        cy.wait(200)
                        cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                    })
                })
            })
        })
    })
})
