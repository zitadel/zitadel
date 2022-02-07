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
                cy.get('[data-e2e=app-card-add]').click()
                cy.get('[formcontrolname^=name]').type(testAppName)
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
                cy.contains('arrow_back').click({ force: true })
                cy.contains('[data-e2e=app-card]', testAppName)
            })
        })
    }) 
})