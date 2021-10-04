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
             * @example cy.ssoLogin('hodor', 'hodor1234')
             */
             ssoLogin(user: User): void

            /**
             * Custom command that authenticates a user.
             *
             * @example cy.consolelogin('hodor', 'hodor1234')
             */
            consolelogin(username: string, password: string): void
            
            /**
             * Custom command that returns a valid Authorization header value.
             *
             * @example cy.apiAuthHeader().then(header) => cy.request({ headers: {Authorization: header }})
             */
            apiAuthHeader(): Chainable<apiCallProperties>
        }
    }
}

Cypress.Commands.add('ssoLogin', { prevSubject: false }, (user: User) => {

    let creds = credentials(user)

    cy.session(creds.username, () => {

        const accountsHost = `accounts.${Cypress.env('apiCallsDomain')}`

        const cookies = new Map<string, string>()

        cy.intercept({
            method: 'GET',
            hostname: accountsHost,
            url: '/login*',
            times: 1
        }, (req) => {
            req.headers['cookie'] = requestCookies(cookies)
            req.continue((res) => {
                updateCookies(res.headers['set-cookie'] as string[], cookies)
            })
        }).as('login')

        cy.intercept({
            method: 'POST',
            hostname: accountsHost,
            url: '/loginname*',
            times: 1
        }, (req) => {
            req.headers['cookie'] = requestCookies(cookies)
            req.continue((res) => {
                updateCookies(res.headers['set-cookie'] as string[], cookies)
            })
        }).as('loginName')

        cy.intercept({
            method: 'POST',
            hostname: accountsHost,
            url: '/password*',
            times: 1
        }, (req) => {
            req.headers['cookie'] = requestCookies(cookies)
            req.continue((res) => {
                updateCookies(res.headers['set-cookie'] as string[], cookies)
            })
        }).as('password')

        cy.intercept({
            method: 'GET',
            hostname: accountsHost,
            url: '/login/success*',
            times: 1
        }, (req) => {
            req.headers['cookie'] = requestCookies(cookies)
            req.continue((res) => {
                updateCookies(res.headers['set-cookie'] as string[], cookies)
            })
        }).as('success') 

        cy.intercept({
            method: 'GET',
            hostname: accountsHost,
            url: '/oauth/v2/authorize/callback*',
            times: 1
        }, (req) => {
            req.headers['cookie'] = requestCookies(cookies)
            req.continue((res) => {
                updateCookies(res.headers['set-cookie'] as string[], cookies)
            })
        }).as('callback')    
        
        cy.intercept({
            method: 'GET',
            url: `https://${accountsHost}/oauth/v2/authorize*`,
            hostname: accountsHost,
            times: 1,
        }, (req) => {
            req.continue((res) => {
                updateCookies(res.headers['set-cookie'] as string[], cookies)
            })
        })

        cy.visit(Cypress.env('consoleUrl'));

        cy.wait('@login')
        cy.get('#loginName').type(creds.username)
        cy.get('#submit-button').click()

        cy.wait('@loginName')
        cy.get('#password').type(creds.password) 
        cy.get('#submit-button').click()

        cy.wait('@callback')

        cy.location('pathname', {timeout: 5 * 1000}).should('eq', '/');

    }, {
        validate: () => {
            cy.visit(`${Cypress.env('consoleUrl')}/users/me`)
        }        
    })

})
    
Cypress.Commands.add('consolelogin', { prevSubject: false }, (username: string, password: string) => {

    window.sessionStorage.removeItem("zitadel:access_token")
    cy.visit(Cypress.env('consoleUrl')).then(() => {
        // fill the fields and push button
        cy.get('#loginName').type(username, { log: false })
        cy.get('#submit-button').click()
        cy.get('#password').type(password, { log: false })
        cy.get('#submit-button').click()
        cy.location('pathname', {timeout: 5 * 1000}).should('eq', '/');
    })
})


Cypress.Commands.add('apiAuthHeader', { prevSubject: false }, () => {

    const apiDomain = Cypress.env('apiCallsDomain')
    const apiBaseURL = `https://api.${apiDomain}`

    // TODO: Why can't I just receive the correct value with Cypress.env('zitadelProjectResourceId')???
    var zitadelProjectResourceID = apiDomain == 'zitadel.ch' ? '69234237810729019' : '70669147545070419'

    var key = Cypress.env("serviceAccountKey")

    var now = new Date().getTime()
    var iat = Math.floor(now / 1000)
    var exp = Math.floor(new Date(now + 1000 * 60 * 55).getTime() / 1000) // 55 minutes
    var bearerToken = sign({
        iss: key.userId,
        sub: key.userId,
        aud: `https://issuer.${apiDomain}`,
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
        url: `${apiBaseURL}/oauth/v2/token`,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: {
            'grant_type': 'urn:ietf:params:oauth:grant-type:jwt-bearer',
            scope: `openid urn:zitadel:iam:org:project:id:${zitadelProjectResourceID}:aud`,
            assertion: bearerToken,
        }
    }).its('body.access_token').then(token => {
        return <apiCallProperties>{
            authHeader: `Bearer ${token}`,
            baseURL: apiBaseURL,
        }
    })
})

function updateCookies(newCookies: string[], currentCookies: Map<string, string>) {
    newCookies.forEach(cs => {
        cs.split('; ').forEach(cookie => {
            const idx = cookie.indexOf('=')
            currentCookies.set(cookie.substring(0,idx), cookie.substring(idx+1))
        })
    })
}

function requestCookies(currentCookies: Map<string, string>): string[] {
    let list = []
    currentCookies.forEach((val, key) => {
        list.push(key+"="+val)
    })
    return list
}

export enum User {
    OrgOwner = 'org_owner',
    OrgOwnerViewer = 'org_owner_viewer',
    OrgProjectCreator = 'org_project_creator',
}

function credentials(user: User) {
    
    return {
        username: `${user}_user_name@caos-demo.${Cypress.env('apiCallsDomain')}`,
        password: Cypress.env(`${user}_password`)
    }
}