import { E2E_API_CALLS_DOMAIN, E2E_SERVICEACCOUNT_KEY } from "../playwright.config"
import { sign } from 'jsonwebtoken'
import fetch from 'node-fetch'
import { checkStatus } from "../fetch/status";

export interface APICallProperties {
    authHeader: string
    baseURL: string
}

export async function prepareAPICalls(): Promise<APICallProperties> {

    const apiBaseURL = `https://api.${E2E_API_CALLS_DOMAIN}`

    // TODO: Why can't I just receive the correct value with process.env.ZITADEL_PROJECT_RESOURCE_ID
    // zitadelProjectResourceID = ZITADEL_PROJECT_RESOURCE_ID
    var zitadelProjectResourceID = E2E_API_CALLS_DOMAIN == 'zitadel.ch' ? '69234237810729019' : '70669147545070419'

    var key = JSON.parse(E2E_SERVICEACCOUNT_KEY)

    var now = new Date().getTime()
    var iat = Math.floor(now / 1000)
    var exp = Math.floor(new Date(now + 1000 * 60 * 55).getTime() / 1000) // 55 minutes
    var bearerToken = sign({
        iss: key.userId,
        sub: key.userId,
        aud: `https://issuer.${E2E_API_CALLS_DOMAIN}`,
        iat: iat,
        exp: exp
    }, key.key, {
        header: {
            alg: "RS256",
            kid: key.keyId
        }
    })

    const payloadContent = {
        'grant_type': 'urn:ietf:params:oauth:grant-type:jwt-bearer',
        scope: `openid urn:zitadel:iam:org:project:id:${zitadelProjectResourceID}:aud`,
        assertion: bearerToken,
    }

    var payloadForm = [];
    for (var property in payloadContent) {
      var encodedKey = encodeURIComponent(property);
      var encodedValue = encodeURIComponent(payloadContent[property]);
      payloadForm.push(encodedKey + "=" + encodedValue);
    }

    const response = await fetch(`${apiBaseURL}/oauth/v2/token`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: payloadForm.join("&")
    })

    checkStatus(response)

    const body = await response.json()

    return {
        baseURL: apiBaseURL,
        authHeader: `Bearer ${body['access_token']}`
    }
}
