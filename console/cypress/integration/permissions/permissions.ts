import { apiAuth, apiCallProperties } from "../../support/api/apiauth";
import { ensureProjectExists, ensureProjectResourceDoesntExist, Roles } from "../../support/api/projects";
import { login, User } from "../../support/login/users";

describe.only('permissions', () => {

    const testProjectName = 'e2eprojectpermission'
    const testAppName = 'e2eapppermission'
    const testRoleName = 'e2eroleundertestname'
    const testRoleDisplay = 'e2eroleundertestdisplay'
    const testRoleGroup = 'e2eroleundertestgroup'
    const testGrantName = 'e2egrantundertest'

    ;[User.OrgOwner].forEach(user => {

        describe(`as user "${user}"`, () => {

            var api: apiCallProperties
            var projectId: number

            beforeEach(() => {
                login(user)
                apiAuth().then(apiCalls => {
                    api = apiCalls
                    ensureProjectExists(apiCalls, testProjectName).then(projId => {
                        projectId = projId
                        cy.visit(`${Cypress.env('consoleUrl')}/projects/${projId}`)
                    })
                })
            })


            describe('add role', () => {
                beforeEach(()=> {
                    ensureProjectResourceDoesntExist(api, projectId, Roles, testRoleName)
                })

                it.only('should add a role', () => {
                    cy.contains('[data-e2e=app-card]', 'Roles').within(() => {
                        cy.contains('a', 'New').click({ force: true })
                    })
                    cy.get('[formcontrolname^=key]').type(testRoleName)
                    cy.get('[formcontrolname^=displayName]').type(testRoleDisplay)
                    cy.get('[formcontrolname^=group]').type(testRoleGroup)
                    cy.contains('button', 'Save').should('be.visible').click()
                })
            })
        })
    })
})
/*

describe('permissions', () => {

    before(()=> {
//        cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    })

    it('should show projects ', () => {
        cy.visit(Cypress.env('consoleUrl') + '/projects')
        cy.url().should('contain', '/projects')
    })

    it('should add a role', () => {
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url().should('contain', '/projects');
            cy.get('.card').should('contain.text', "newProjectToTest")
        })
        cy.get('.card').filter(':contains("newProjectToTest")').click()
        cy.get('.app-container').filter(':contains("newAppToTest")').should('be.visible').click()
        let projectID
        cy.url().then(url => {
            cy.log(url.split('/')[4])
            projectID = url.split('/')[4]
        });
        
        cy.then(() => cy.visit(Cypress.env('consoleUrl') + '/projects/' + projectID +'/roles/create'))
        cy.get('[formcontrolname^=key]').type("newdemorole")
        cy.get('[formcontrolname^=displayName]').type("newdemodisplayname")
        cy.get('[formcontrolname^=group]').type("newdemogroupname")
        cy.get('button').filter(':contains("Save")').should('be.visible').click()
        //let the Role get processed
        cy.wait(5000)
    })

    it('should add a grant', () => {
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url().should('contain', '/projects');
            cy.get('.card').should('contain.text', "newProjectToTest")
        })
        cy.get('.card').filter(':contains("newProjectToTest")').click()
        cy.get('.app-container').filter(':contains("newAppToTest")').should('be.visible').click()
        let projectID
        cy.url().then(url => {
            cy.log(url.split('/')[4])
            projectID = url.split('/')[4]
        });
        
        cy.then(() => cy.visit(Cypress.env('consoleUrl') + '/grant-create/project/' + projectID ))
        cy.get('input').type("demo")
        cy.get('[role^=listbox]').filter(`:contains("${Cypress.env("fullUserName")}")`).should('be.visible').click()
        cy.wait(5000)
        //cy.get('.button').contains('Continue').click()
        cy.get('button').filter(':contains("Continue")').click()
        cy.wait(5000)
        cy.get('tr').filter(':contains("demo")').find('label').click()
        cy.get('button').filter(':contains("Save")').should('be.visible').click()
        //let the grant get processed
        cy.wait(5000)
    })
})

*/