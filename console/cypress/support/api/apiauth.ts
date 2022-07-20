import { sign } from 'jsonwebtoken'

export interface apiCallProperties {
    authHeader: string
    mgntBaseURL: string
}

export function apiAuth(): Cypress.Chainable<apiCallProperties> {
    const baseUrl = Cypress.env('baseUrl')
    const issuerUrl = `${baseUrl}/oauth/v2`
    const zitadelProjectResourceID = (<string>Cypress.env('zitadelProjectResourceId')).replace('bignumber-', '')

    const key = Cypress.env("parsedServiceAccountKey")

    const now = new Date().getTime()
    const iat = Math.floor(now / 1000)
    const exp = Math.floor(new Date(now + 1000 * 60 * 55).getTime() / 1000) // 55 minutes
    const bearerToken = sign({
        iss: key.userId,
        sub: key.userId,
        aud: `${baseUrl}`,
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
        url: `${issuerUrl}/token`,
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
            mgntBaseURL: `${baseUrl}/management/v1/`,
        }
    })
}