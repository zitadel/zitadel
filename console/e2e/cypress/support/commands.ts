import { sign } from 'jsonwebtoken'

interface apiCallProperties {
    authHeader: string
    baseURL: string
}

declare global {
    namespace Cypress {
        interface Chainable {
            /**
             * Custom command that authenticates a user.
             * 
             * @example cy.consolelogin('hodor', 'hodor1234')
             */
            consolelogin(username: string, password: string): Chainable<void>

            /**
             * Custom command that returns a valid Authorization header value.
             * 
             * @example cy.apiAuthHeader().then(header) => cy.request({ headers: {Authorization: header }})
             */
            apiAuthHeader(): Chainable<apiCallProperties>
        }
    }
}

Cypress.Commands.add('consolelogin', { prevSubject: false }, (username: string, password: string) => {
    window.sessionStorage.removeItem("zitadel:access_token")
    cy.visit(Cypress.env('consoleUrl')) 
    // fill the fields and push button
    cy.get('#loginName').type(username, { log: false })
    cy.get('#submit-button').click()
    cy.get('#password').type(password, { log: false })
    cy.get('#submit-button').click()
    cy.location('pathname', {timeout: 5 * 1000}).should('eq', '/');
})

Cypress.Commands.add('apiAuthHeader', { prevSubject: false }, () => {
    var key = Cypress.env("serviceAccountKey")

    var now = new Date().getTime()
    var iat = Math.floor(now / 1000)
    var exp = Math.floor(new Date(now + 1000 * 60 * 55).getTime() / 1000) // 55 minutes
    var bearerToken = sign({
        iss: key.userId,
        sub: key.userId,
        aud: `https://issuer.${Cypress.env('domain')}`,
        iat: iat,
        exp: exp
    }, key.key, {
        header: {
            alg: "RS256", 
            kid: key.keyId
        }
    })

    const baseURL = `https://api.${Cypress.env('domain')}`

    cy.request({
        method: 'POST',
        url: `${baseURL}/oauth/v2/token`,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: {
            'grant_type': 'urn:ietf:params:oauth:grant-type:jwt-bearer',
            scope: 'openid urn:zitadel:iam:org:project:id:69234237810729019:aud',
            assertion: bearerToken,
        }
    }).its('body.access_token').then(token => {
        return <apiCallProperties>{
            authHeader: `Bearer ${token}`,
            baseURL: baseURL,
        }
    })
})