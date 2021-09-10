import { sign } from 'jsonwebtoken'

// NEEDS TO BE DISABLED!!!!!! this is just for testing
Cypress.on('uncaught:exception', (err, runnable) => {
    // returning false here prevents Cypress from
    if (err.message.includes('addEventListener')) {
        return false
    }
})
// ###############################

describe("machines", ()=> {

    before(()=> {
//        cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    })

    it('should show machines', () => {
        cy.visit(Cypress.env('consoleUrl') + '/users/list/machines')
        cy.url().should('contain', 'users/list/machines')
    })

    describe('add', () => {

        it.only('should cleanup', () => {
            
            var key = Cypress.env("serviceAccountKey")

            var now = new Date().getTime()
            var iat = Math.floor(now / 1000)
            var exp = Math.floor(new Date(now + 1000 * 60 * 55).getTime() / 1000) // 55 minutes
            var bearerToken = sign({
                iss: key.userId,
                sub: key.userId,
                aud: "https://issuer.zitadel.ch",
                iat: iat,
                exp: exp
            }, key.key, {
                header: {
                    alg: "RS256", 
                    kid: key.keyId
                }
            })

            cy.request({
                method: 'POST',
                url: 'https://api.zitadel.ch/oauth/v2/token',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
                body: {
                    'grant_type': 'urn:ietf:params:oauth:grant-type:jwt-bearer',
                    scope: 'openid urn:zitadel:iam:org:project:id:69234237810729019:aud',
                    assertion: bearerToken,
                }
            }).then(res => {
                cy.request({
                    method: 'POST',
                    url: 'https://api.zitadel.ch/management/v1/users/_search',
                    headers: {
                        Authorization: "Bearer " + res.body['access_token']
                    },
                    qs: {
                        'user_name_query': "cypress"
                    }
                }).then(res => {
                    debugger;
                    cy.log('all users', res)
                })
            })
        })

        before('cleanup', () => {
            //click on org to clear screen
/*            cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
                cy.url().should('contain', '/org');
            })
            cy.visit(Cypress.env('consoleUrl') + '/users/list/machines').then(() => {
                cy.url().should('contain', '/users/list/machines');
                cy.get('h1')
                    .contains('Service Users')
                    .parent()
                    .find("tr", { timeout: 50 })
                    .filter(':contains("machineusername")', {timeout: 200 })
                    .find('button')
                    //force due to angular hidden buttons
                    .click({force: true})
                    .then(() => {
                        cy.get('span.title')
                            .contains('Delete User')
                            .parent()
                            .find('button')
                            .contains('Delete')
                            .click()
                            .then(() => {
                                cy.wait(3000)
                                cy.visit(Cypress.env('consoleUrl') + '/users/list/machines');
                                cy.get('[text*=machineusername]').should('not.exist');
                            })
                    })
            })*/
        })

        it('should add a machine', () => {
            //click on org to clear screen
            cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
                cy.url().should('contain', '/org');
            })
            cy.visit(Cypress.env('consoleUrl') + '/users/list/machines')
            cy.url().should('contain', 'users/list/machines')
            cy.visit(Cypress.env('consoleUrl') + '/users/create-machine')
            cy.url().should('contain', 'users/create-machine')
            //force needed due to the prefilled username prefix
            cy.get('[formcontrolname^=userName]').type(Cypress.env('newMachineUserName'),{force: true})
            cy.get('[formcontrolname^=name]').type(Cypress.env('newMachineName'))
            cy.get('[formcontrolname^=description]').type(Cypress.env('newMachineDesription'))
            cy.get('button').filter(':contains("Create")').should('be.visible').click().then(() => {
                cy.wait(3000)
                cy.visit(Cypress.env('consoleUrl') + '/users/list/machines');
                cy.get('tr').should('contain.text', "machineusername").and('exist');
            })
        })
    })
})
