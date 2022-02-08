import { apiAuth } from "../../support/api/apiauth";
import { ensureMachineUserExists, ensureUserDoesntExist } from "../../support/api/users";
import { login, User, username } from "../../support/login/users";

// NEEDS TO BE DISABLED!!!!!! this is just for testing
/*
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
        return false
    }
})
 */
// ###############################

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
                    cy.contains('a', 'New').click()
                    cy.url().should('contain', 'users/create-machine')
                    //force needed due to the prefilled username prefix
                    cy.get('[formcontrolname^=userName]').type(testMachineUserName,{force: true})
                    cy.get('[formcontrolname^=name]').type('e2emachinename')
                    cy.get('[formcontrolname^=description]').type('e2emachinedescription')
                    cy.get('button').filter(':contains("Create")').should('be.visible').click()
                    cy.contains('User created successfully')
                    cy.visit(machinesPath);
                    cy.wait(10_000) // TODO: eventual consistency ftw
                    cy.contains('button', 'refresh').click()
                    cy.contains("tr", testMachineUserName)
                })
            })

            describe('remove', () => {
                before('ensure it exists', () => {
                    apiAuth().then(api => {
                        ensureMachineUserExists(api, testMachineUserName)
                    })
                })

                it('should delete a machine', () => {
                    cy.get('h1')
                        .contains('Service Users')
                        .parent()
                        .contains("tr", testMachineUserName, { timeout: 1000 })
                        .find('button')
                        //force due to angular hidden buttons
                        .click({force: true})
                    cy.get('span.title')
                        .contains('Delete User')
                        .parent()
                        .find('button')
                        .contains('Delete')
                        .click()
                    cy.contains('mat-dialog-container', 'Delete User').find('input').type(username(testMachineUserName, Cypress.env('org')))
                    cy.contains('mat-dialog-container button', 'Delete').click()    
                    cy.contains('User deleted successfully')
                    cy.wait(10_000) // TODO: eventual consistency ftw
                    cy.contains('button', 'refresh').click()
                    cy.get(`[text*=${testMachineUserName}]`).should('not.exist');
                })
            })
        })
    })
})
