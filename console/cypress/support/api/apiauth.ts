import { sign } from 'jsonwebtoken'

export interface apiCallProperties {
    authHeader: string
    mgntBaseURL: string
}

export function apiAuth(): Cypress.Chainable<apiCallProperties> {
    const apiUrl = Cypress.env('apiUrl')
    const issuerUrl = Cypress.env('issuerUrl')
    const zitadelProjectResourceID = (<string>Cypress.env('zitadelProjectResourceId')).replace('bignumber-', '')

    debugger
    const key = Cypress.env("serviceAccountKey")

    const now = new Date().getTime()
    const iat = Math.floor(now / 1000)
    const exp = Math.floor(new Date(now + 1000 * 60 * 55).getTime() / 1000) // 55 minutes
    const bearerToken = sign({
        iss: key.userId,
        sub: key.userId,
        aud: `${issuerUrl}/oauth/v2`,
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
        url: `${apiUrl}/oauth/v2/token`,
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
            mgntBaseURL: `${apiUrl}/management/v1/`,
        }
    })
}