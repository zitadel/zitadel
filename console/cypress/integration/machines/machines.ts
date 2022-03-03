import { apiAuth } from "../../support/api/apiauth";
import { ensureMachineUserExists, ensureUserDoesntExist } from "../../support/api/users";
import { login, User, username } from "../../support/login/users";

describe('machines', () => {

    const machinesPath = `${Cypress.env('consoleUrl')}/users/list/machines`
    const testMachineUserName = 'e2emachineusername'
    
    ;[User.OrgOwner].forEach(user => {

        describe(`as user "${user}"`, () => {

            beforeEach(()=> {
                login(user)
                cy.visit(machinesPath)
                cy.get('[data-cy=timestamp]')
            })

            describe('add', () => {

                before(`ensure it doesn't exist already`, () => {
                    apiAuth().then(apiCallProperties => {
                        ensureUserDoesntExist(apiCallProperties, testMachineUserName)
                    })
                })

                it('should add a machine', () => {
                    cy.get('a[href="/users/create-machine"]').click()
                    cy.url().should('contain', 'users/create-machine')
                    //force needed due to the prefilled username prefix
                    cy.get('[formcontrolname="userName"]').type(testMachineUserName,{force: true})
                    cy.get('[formcontrolname="name"]').type('e2emachinename')
                    cy.get('[formcontrolname="description"]').type('e2emachinedescription')
                    cy.get('[type="submit"]').should('be.visible').click()
                    cy.get('.data-e2e-success')
                    cy.wait(1000)
                    cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                })
            })

            describe.only('remove', () => {
                before('ensure it exists', () => {
                    apiAuth().then(api => {
                        ensureMachineUserExists(api, testMachineUserName)
                    })
                })

                it('should delete a machine', () => {
                    cy.contains("tr", testMachineUserName, { timeout: 1000 })
                        .find('button')
                        //force due to angular hidden buttons
                        .click({force: true})
                    cy.get('mat-dialog-container input').type(username(testMachineUserName, Cypress.env('org')))
                    cy.get('[e2e-data="confirm-dialog-button"]').click()
                    cy.get('.data-e2e-success')
                    cy.wait(1000)
                    cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                })
            })
        })
    })
})
