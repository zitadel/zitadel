import { login, User } from "../../support/login/users";
import { Apps, ensureProjectExists, ensureProjectResourceDoesntExist } from "../../support/api/projects";
import { apiAuth } from "../../support/api/apiauth";

describe('applications', () => {

    const testProjectName = 'e2eprojectapplication'
    const testAppName = 'e2eappundertest'

    ;[User.OrgOwner].forEach(user => {

        describe(`as user "${user}"`, () => {

            beforeEach(`ensure it doesn't exist already`, () => {
                login(user)
                apiAuth().then(api => {
                    ensureProjectExists(api, testProjectName).then(projectID => {
                        ensureProjectResourceDoesntExist(api, projectID, Apps, testAppName).then(() => {
                            cy.visit(`${Cypress.env('consoleUrl')}/projects/${projectID}`)
                        })
                    })
                })
            })

            it('add app', () => {
                cy.get('mat-spinner')
                cy.get('mat-spinner').should('not.exist')
                cy.get('[data-e2e="app-card-add"]').should('be.visible').click()
                // select webapp
                cy.get('[formcontrolname="name"]').type(testAppName)
                cy.get('[for="WEB"]').click()
                cy.get('[data-e2e="continue-button-nameandtype"]').click()
                //select authentication
                cy.get('[for="PKCE"]').click()
                cy.get('[data-e2e="continue-button-authmethod"]').click()
                //enter URL
                cy.get('cnsl-redirect-uris').eq(0).type("https://testurl.org")
                cy.get('cnsl-redirect-uris').eq(1).type("https://testlogouturl.org")
                cy.get('[data-e2e="continue-button-redirecturis"]').click()
                cy.get('[data-e2e="create-button"]').click().then(() => {
                    cy.get('[id*=overlay]').should('exist')
                })
                cy.get('.data-e2e-success')
                cy.wait(200)
                cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                //TODO: check client ID/Secret
            })
        })
    }) 
})