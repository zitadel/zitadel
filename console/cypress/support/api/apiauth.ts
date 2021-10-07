import { sign } from 'jsonwebtoken'

export interface apiCallProperties {
    authHeader: string
    mgntBaseURL: string
}

export function apiAuth(): Cypress.Chainable<apiCallProperties> {
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

    return cy.request({
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
            mgntBaseURL: `${apiBaseURL}/management/v1/`,
        }
    })
}