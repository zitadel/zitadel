import { sign } from 'jsonwebtoken'
import { getUnixTime } from 'date-fns';
import { debug } from 'console';

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
             ssoLogin(username: string, password: string): void

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


/**
 * Hit the local login endpoint in the application which will redirect to Auth0.
 */
 function startLogin() {
    return cy.request({
      url: 'http://localhost:8080/login',
      followRedirect: false
    });
  }
  
  /**
   * Universal Login
   * @param {*} user
   * @param {*} loginUrl
   */
  function followUniversalLogin(user, loginUrl) {
    return cy.task('LoginPuppeteer', {
      username: user.email,
      password: user.password,
      loginUrl,
      callbackUrl: 'http://localhost:8080/callback'
    });
  }

  function createCookie(cookie: any, domain: string) {
    return {
      name: cookie.name,
      value: cookie.value,
      options: {
        domain: `${cookie.domain}`,
        expiry: getFutureTime(15),
        httpOnly: cookie.httpOnly,
        path: cookie.path,
        sameSite: cookie.sameSite,
        secure: cookie.secure,
        session: cookie.session
      }
    };
  }

  function getFutureTime(minutesInFuture) {
    const time = new Date(new Date().getTime() + minutesInFuture * 60000);
    return getUnixTime(time);
  }

Cypress.Commands.add('ssoLogin', { prevSubject: false }, (username: string, password: string) => {
    cy.task('login', { username: username, password: password }, { timeout: 30000 }).then(({ cookies, callbackUrl }) => {
        cy.intercept(`https://accounts.${Cypress.env('apiCallsDomain')}/oauth/v2/authorize*`, (req) => {
            req.headers['x-custom-headers'] = 'added by cy.intercept'
//            req.url = "https://example.com"
//            req.continue()
        }).as('headers')
        var consoleUrl: string = Cypress.env('consoleUrl')
        consoleUrl = consoleUrl.substr(consoleUrl.indexOf("://")+3)

/*        cy.setCookie("elio", "dorro", {domain: "accounts.zitadel.dev"}).then(() => {
            cy.visit(callbackUrl);
        })*/

        cy.wait('@headers', { timeout: 30_000 }).then(req => {
            debugger
        })        


        const cookie0 = createCookie(cookies[0],consoleUrl)
        cy.setCookie(cookie0.name, cookie0.value, cookie0.options).then(() => {
            const cookie1 = createCookie(cookies[1],consoleUrl)
            cy.setCookie(cookie1.name, cookie1.value, cookie1.options).then(() => {
                const cookie2 = createCookie(cookies[2],consoleUrl)
                cy.setCookie(cookie2.name, cookie2.value, cookie2.options).then(() => {
                    cy.visit(callbackUrl);
                })
            })
        })
    })
})
    
Cypress.Commands.add('consolelogin', { prevSubject: false }, (username: string, password: string) => {

//    cy.setCookie('__Secure-caos.zitadel.useragent', 'MTYzMjIyMDc4MnxXQUdsTloyNjJhX0xKYkFQYUduUlo2cWZGbXBqSTUwMHFtZ2Rqa3JjSDJfbV9LM3p6TXVRMUdac1hWYUsxNzNFdEF0WDVQaUJoWExHMjZ4U3FGRkZKWVU2bWp2a19Gaz18FV5k9ZcbWfmw7VLpsIFCzR4EIeM9owJnWDc7OeBbOSA=', {secure: true})

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
