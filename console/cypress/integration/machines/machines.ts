import { apiAuth } from "../../support/api/apiauth";
import { ensureMachineUserExists, ensureUserDoesntExist } from "../../support/api/users";
import { login, User, username } from "../../support/login/users";

describe('machines', () => {

    const machinesPath = `${Cypress.env('consoleUrl')}/users/list/machines`
    const testMachineUserNameAdd = 'e2emachineusernameadd'
    const testMachineUserNameRemove = 'e2emachineusernameremove'
    
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
                        ensureUserDoesntExist(apiCallProperties, testMachineUserNameAdd)
                    })
                })

                it('should add a machine', () => {
                    cy.get('a[href="/users/create-machine"]').click()
                    cy.url().should('contain', 'users/create-machine')
                    //force needed due to the prefilled username prefix
                    cy.get('[formcontrolname="userName"]').type(testMachineUserNameAdd,{force: true})
                    cy.get('[formcontrolname="name"]').type('e2emachinename')
                    cy.get('[formcontrolname="description"]').type('e2emachinedescription')
                    cy.get('[type="submit"]').should('be.visible').click()
                    cy.get('.data-e2e-success')
                    cy.wait(200)
                    cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                })
            })

            describe('remove', () => {
                before('ensure it exists', () => {
                    apiAuth().then(api => {
                        ensureMachineUserExists(api, testMachineUserNameRemove)
                    })
                })

                it('should delete a machine', () => {
                    cy.contains("tr", testMachineUserNameRemove, { timeout: 1000 })
                        .find('button')
                        //force due to angular hidden buttons
                        .click({force: true})
                    cy.get('[e2e-data="confirm-dialog-input"]').type(username(testMachineUserNameRemove, Cypress.env('org')))
                    cy.get('[e2e-data="confirm-dialog-button"]').click()
                    cy.get('.data-e2e-success')
                    cy.wait(200)
                    cy.get('.data-e2e-failure', { timeout: 0 }).should('not.exist')
                })
            })
        })
    })
})
