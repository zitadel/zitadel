import { User } from "../../support/commands"

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
    const testMachineDescription = 'e2emachinedescription'
    const testMachineName = 'e2emachinename'
    
    ;[User.OrgOwner].forEach(user => {

        describe(`as user "${user}"`, () => {

            beforeEach(()=> {
                cy.ssoLogin(user)
                cy.visit(machinesPath)
                cy.get('[data-cy=timestamp]')
            })

            describe('add', () => {

                before(`ensure it doesn't exist already`, () => {
                    cy.apiAuthHeader().then(apiCallProperties => {
                        cy.request({
                            method: 'POST',
                            url: `${apiCallProperties.baseURL}/management/v1/users/_search`,
                            headers: {
                                Authorization: apiCallProperties.authHeader
                            },
                        }).then(usersRes => {
                            var machineUser = usersRes.body.result.find(user => user.userName === testMachineUserName)
                            if (machineUser) {
                                cy.request({
                                    method: 'DELETE',
                                    url: `${apiCallProperties.baseURL}/management/v1/users/${machineUser.id}`,
                                    headers: {
                                        Authorization: apiCallProperties.authHeader
                                    },
                                })
                            }
                        })
                    })
                })

                it('should add a machine', () => {
                    cy.contains('a', 'New').click()
                    cy.url().should('contain', 'users/create-machine')
                    //force needed due to the prefilled username prefix
                    cy.get('[formcontrolname^=userName]').type(testMachineUserName,{force: true})
                    cy.get('[formcontrolname^=name]').type(testMachineName)
                    cy.get('[formcontrolname^=description]').type(testMachineDescription)
                    cy.get('button').filter(':contains("Create")').should('be.visible').click()
                    cy.contains('User created successfully')
                    cy.visit(machinesPath);
                    cy.contains("tr", testMachineUserName)
                })
            })

            describe('remove', () => {
                before('ensure it exists', () => {
                    cy.apiAuthHeader().then(apiCallProperties => {
                        cy.request({
                            method: 'POST',
                            url: `${apiCallProperties.baseURL}/management/v1/users/machine`,
                            headers: {
                                Authorization: apiCallProperties.authHeader
                            },
                            body: {
                                user_name: testMachineUserName,
                                name: testMachineName,
                                description: testMachineDescription,
                            },
                            failOnStatusCode: false,
                            followRedirect: false
                        }).then(res => {
                            expect(res.status).to.be.oneOf([200,409])
                        })
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
                    cy.contains('User deleted successfully')
                    cy.get(`[text*=${testMachineUserName}]`).should('not.exist');
                })
            })
        })
    })
})
