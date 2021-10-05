import { User } from "../../support/commands"

describe('humans', () => {

    const humansPath = `${Cypress.env('consoleUrl')}/users/list/humans`
    const testHumanUserName = 'e2ehumanusername'
    const testHumanFirstName = 'e2ehumanfirstname'
    const testHumanLastName = 'e2ehumanlastname'
    const testHumanEmail = `e2ehuman@${Cypress.env('apiCallsDomain')}`
    const testHumanPhone = '+41 123456789'    

    ;[User.OrgOwner].forEach(user => {

        describe(`as user "${user}"`, () => {

            beforeEach(()=> {
                cy.ssoLogin(user)
                cy.visit(humansPath)
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
                            var humanUser = usersRes.body.result.find(user => user.userName === testHumanUserName)
                            if (humanUser) {
                                cy.request({
                                    method: 'DELETE',
                                    url: `${apiCallProperties.baseURL}/management/v1/users/${humanUser.id}`,
                                    headers: {
                                        Authorization: apiCallProperties.authHeader
                                    },
                                })
                            }
                        })
                    })
                })

                it('should add a user', () => {
                    cy.contains('a', 'New').click()
                    cy.url().should('contain', 'users/create')
                    cy.get('[formcontrolname^=email]').type(testHumanEmail)
                    //force needed due to the prefilled username prefix
                    cy.get('[formcontrolname^=userName]').type(testHumanUserName, {force: true})
                    cy.get('[formcontrolname^=firstName]').type(testHumanFirstName)
                    cy.get('[formcontrolname^=lastName]').type(testHumanLastName)
                    cy.get('[formcontrolname^=phone]').type(testHumanPhone)
                    cy.get('button').filter(':contains("Create")').should('be.visible').click()
                    cy.contains('User created successfully')
                    cy.visit(humansPath);
                    cy.contains("tr", testHumanUserName)
                })        
            })
            
            describe('remove', () => {
                before('ensure it exists', () => {
                    cy.apiAuthHeader().then(apiCallProperties => {
                        cy.request({
                            method: 'POST',
                            url: `${apiCallProperties.baseURL}/management/v1/users/human`,
                            headers: {
                                Authorization: apiCallProperties.authHeader
                            },
                            body: {
                                user_name: testHumanUserName,
                                profile: {
                                    first_name: testHumanFirstName,
                                    last_name: testHumanLastName,
                                },
                                email: { 
                                    email: testHumanEmail,
                                },
                                phone: {
                                    phone: testHumanPhone,
                                },
                            },
                            failOnStatusCode: false,
                            followRedirect: false
                        }).then(res => {
                            expect(res.status).to.be.oneOf([200,409])
                        })
                    })                    
                })

                it('should delete a human user', () => {
                    cy.get('h1')
                        .contains('Users')
                        .parent()
                        .contains("tr", testHumanUserName, { timeout: 1000 })
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
                    cy.get(`[text*=${testHumanUserName}]`).should('not.exist');                    
                })
            })
        })
    })
})
/*
describe("users", ()=> {

    before(()=> {
        cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    })

    it('should show personal information', () => {
        cy.log(`USER: show personal information`);
        //click on user information 
        cy.get('a[href*="users/me"').eq(0).click()
        cy.url().should('contain', '/users/me')
    })

    it('should show users', () => {
        cy.visit(Cypress.env('consoleUrl') + '/users/list/humans')
        cy.url().should('contain', 'users/list/humans')
    })
})

*/