import { ORG_MANAGER } from "../shared/types"

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

describe('machine', () => {

    ;[ORG_MANAGER.org_owner].forEach(user => {

        describe(`impersonating an organization manager with permission "${user}"`, () => {

            before(()=> {
                cy.consolelogin(`${user.toLowerCase()}_user_name@caos-demo.${Cypress.env('apiCallsDomain')}`, Cypress.env(`${user.toLowerCase()}_password`))
//                cy.ssoLogin(`${user.toLowerCase()}_user_name@caos-demo.${Cypress.env('apiCallsDomain')}`, Cypress.env(`${user.toLowerCase()}_password`))
//                cy.visit(Cypress.env('consoleUrl') + '/users/list/machines')
//                cy.get("app-refresh-table")
            })

            it.only('debug', () => {

            })

            describe(`as user ${user}`, () => {

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
                                var machineUser = usersRes.body.result.find(user => user.userName === Cypress.env('newMachineUserName'))
                                if (machineUser) {
                                    cy.request({
                                        method: 'DELETE',
                                        url: `https://api.zitadel.ch/management/v1/users/${machineUser.id}`,
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
                        cy.get('[formcontrolname^=userName]').type(Cypress.env('newMachineUserName'),{force: true})
                        cy.get('[formcontrolname^=name]').type(Cypress.env('newMachineName'))
                        cy.get('[formcontrolname^=description]').type(Cypress.env('newMachineDesription'))
                        cy.get('button').filter(':contains("Create")').should('be.visible').click()
                        cy.contains('User created successfully')
                        cy.visit(Cypress.env('consoleUrl') + '/users/list/machines');
                        cy.contains("tr", Cypress.env('newMachineUserName'))
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
                                    user_name: Cypress.env('newMachineUserName'),
                                    name: Cypress.env('newMachineName'),
                                    description: Cypress.env('newMachineDesription'),
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
                            .contains("tr", Cypress.env('newMachineUserName'), { timeout: 1000 })
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
                        cy.get(`[text*=${Cypress.env('newMachineUserName')}]`).should('not.exist');
                    })
                })
            })
        })
    })
})
