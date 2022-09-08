import { apiAuth } from "../../support/api/apiauth";
import { ensureProjectExists, ensureProjectResourceDoesntExist, Roles } from "../../support/api/projects";

describe('permissions', () => {

    describe("management", ()=> {

        describe("organizations", () => {
            it("should add an organization manager")
            it("should remove an organization manager")
        })

        describe("projects", () => {

            const testProjectName = 'e2eprojectpermission'
            const testAppName = 'e2eapppermission'
            const testRoleName = 'e2eroleundertestname'
            const testRoleDisplay = 'e2eroleundertestdisplay'
            const testRoleGroup = 'e2eroleundertestgroup'
            const testGrantName = 'e2egrantundertest'

            var projectId: number

            beforeEach(() => {
                apiAuth().then(apiCalls => {
                    ensureProjectExists(apiCalls, testProjectName).then(projId => {
                        projectId = projId
                    })
                })
            })

            describe("managers", () => {
                it("should add a project manager")
                it("should remove a project manager")
            })

            describe("authorizations", () => {
                it('should add an authorization')
                it('should remove an authorization')
            })

            describe("owned projects", () => {

                describe('roles', () => {
                    beforeEach(()=> {
                        apiAuth().then((api)=> {
                            ensureProjectResourceDoesntExist(api, projectId, Roles, testRoleName)
                            cy.visit(`/projects/${projectId}?id=roles`)
                        })
                    })

                    it('should add a role',  () => {
                        cy.get('[data-e2e="add-new-role"]').click()
                        cy.get('[formcontrolname="key"]')
                            .type(testRoleName)
                        cy.get('[formcontrolname="displayName"]')
                            .type(testRoleDisplay)
                        cy.get('[formcontrolname="group"]')
                            .type(testRoleGroup)
                        cy.get('[data-e2e="save-button"]')
                            .click()
                        cy.get('.data-e2e-success')
                        cy.wait(200)
                        cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                    })
                    it('should remove a role')
                })

                describe('grants', () => {
                    it('should add a grant')
                    it('should remove a grant')
                })
            })
        })
    })

    describe('validations', () => {

        describe("owned projects", () => {
            describe("no ownership", () => {
                it("a user without project global ownership can ...")
                it("a user without project global ownership can not ...")
            })
            describe("project owner viewer global", () => {
                it("a project owner viewer global additionally can ...")
                it("a project owner viewer global still can not ...")
            })
            describe("project owner global", () => {
                it("a project owner global additionally can ...")
                it("a project owner global still can not ...")
            })
        })

        describe("granted projects", () => {
            describe("no ownership", () => {
                it("a user without project grant ownership can ...")
                it("a user without project grant ownership can not ...")
            })
            describe("project grant owner viewer", () => {
                it("a project grant owner viewer additionally can ...")
                it("a project grant owner viewer still can not ...")
            })
            describe("project grant owner", () => {
                it("a project grant owner additionally can ...")
                it("a project grant owner still can not ...")
            })
        })
        describe("organization", () => {
            describe("org owner", () => {
                it("a project owner global can ...")
                it("a project owner global can not ...")
            })
        })
    })
})