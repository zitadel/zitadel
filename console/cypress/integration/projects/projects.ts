import { apiAuth } from "../../support/api/apiauth";
import { ensureProjectDoesntExist, ensureProjectExists } from "../../support/api/projects";
import { login, User } from "../../support/login/users";

describe("projects", ()=> {

    const testProjectName = 'e2eproject'

    ;[User.OrgOwner].forEach(user => {

        describe(`as user "${user}"`, () => {

            beforeEach(()=> {
                login(user)
                cy.visit(`${Cypress.env('consoleUrl')}/projects`)
            })

            describe('add project', () => {
                before(`ensure it doesn't exist already`, () => {
                    apiAuth().then(api => {
                        ensureProjectDoesntExist(api, testProjectName)
                    })
                })

                it('should add a project', () => {
                    cy.get('.add-project-button').click({ force: true })
                    cy.get('input').type(testProjectName)
                    cy.get('[type^=submit]').click()
                    cy.get('h1').should('contain', `Project ${testProjectName}`)
                    cy.get('a').contains('arrow_back').click()
                    cy.get('[data-e2e=grid-card]').contains(testProjectName)
                    cy.get('[data-e2e=toggle-grid]').click()
                    cy.contains("tr", testProjectName)
                })
            })

            describe('remove project', () => {
                beforeEach('ensure it exists', () => {
                    apiAuth().then(api => {
                        ensureProjectExists(api, testProjectName)
                    })
                })

                afterEach('project should be deleted', () => {
                    cy.get('span.title')
                        .contains('Delete Project')
                        .parent()
                        .find('button')
                        .contains('Delete')
                        .click()
                    cy.contains('Deleted Project')
                    cy.get(`[text*=${testProjectName}]`).should('not.exist');
                    cy.get('[data-e2e=toggle-grid]').click()
                    cy.get(`[text*=${testProjectName}]`).should('not.exist');
                })

                it.skip('via list view', () => {
                    cy.get('[data-e2e=toggle-grid]').click()
                    cy.get('[data-cy=timestamp]')
                    cy.contains('h1', 'Projects')
                        .parent()
                        .contains("tr", testProjectName, { timeout: 1000 })
                        .find('[data-e2e=delete-project-button]')
                        .click({force: true})
                })

                it('via grid view', () => {
                    cy.contains('[data-e2e=grid-card]', testProjectName)
                        .find('[data-e2e=delete-project-button]')
                        .trigger('mouseover')
                        .click()
                })
            })
        })
    })
})
