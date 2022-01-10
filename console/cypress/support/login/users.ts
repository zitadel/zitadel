export enum User {
    OrgOwner = 'org_owner',
    OrgOwnerViewer = 'org_owner_viewer',
    OrgProjectCreator = 'org_project_creator',
    LoginPolicyUser = 'login_policy_user',
    PasswordComplexityUser = 'password_complexity_user',
    IAMAdminUser = "zitadel-admin"
}

export function login(user:User, force?: boolean, pw?: string, onUsernameScreen?: () => void, onPasswordScreen?: () => void, onAuthenticated?: () => void): void {
    let creds = credentials(user, pw)

    const apiUrl: string = Cypress.env('apiUrl')
    const accountsUrl: string = Cypress.env('accountsUrl')
    const consoleUrl: string = Cypress.env('consoleUrl') 
    const multipleDomains = host(apiUrl) != host(accountsUrl)

    cy.session(creds.username, () => {

        const cookies = new Map<string, string>()

        if (multipleDomains) {
            cy.intercept({
                method: 'GET',
                hostname: "localhost",
                url: '/login/login*',
                times: 1
            }, (req) => {
                req.headers['cookie'] = requestCookies(cookies)
                req.continue((res) => {
                    updateCookies(res.headers['set-cookie'] as string[], cookies)
                })
            }).as('login')

            cy.intercept({
                method: 'POST',
                hostname: "localhost",
                url: '/login/loginname*',
                times: 1
            }, (req) => {
                req.headers['cookie'] = requestCookies(cookies)
                req.continue((res) => {
                    updateCookies(res.headers['set-cookie'] as string[], cookies)
                })
            }).as('loginName')

            cy.intercept({
                method: 'POST',
                hostname: "localhost",
                url: '/login/password*',
                times: 1
            }, (req) => {
                req.headers['cookie'] = requestCookies(cookies)
                req.continue((res) => {
                    updateCookies(res.headers['set-cookie'] as string[], cookies)
                })
            }).as('password')

            cy.intercept({
                method: 'GET',
                hostname: "localhost",
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
                hostname: "localhost",
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
                url: `${accountsUrl}/oauth/v2/authorize*`,
                hostname: "localhost",
                times: 1,
            }, (req) => {
                req.continue((res) => {
                    updateCookies(res.headers['set-cookie'] as string[], cookies)
                })
            })
        }

        cy.visit(`${consoleUrl}/loginname`, { retryOnNetworkFailure: true });

        multipleDomains && cy.wait('@login')
        onUsernameScreen ? onUsernameScreen() : null
        cy.get('#loginName').type(creds.username)
        cy.get('#submit-button').click()

        multipleDomains && cy.wait('@loginName')
        onPasswordScreen ? onPasswordScreen() : null
        cy.get('#password').type(creds.password) 
        cy.get('#submit-button').click()

        onAuthenticated ? onAuthenticated() : null

        multipleDomains && cy.wait('@callback')

        cy.location('pathname', {timeout: 5 * 1000}).should('eq', '/');

    }, {        
        validate: () => {

            if (force) {
                throw new Error("clear session");
            }

            cy.visit(`${consoleUrl}/users/me`)
        }
    })
}



export function username(withoutDomain: string, project?: string): string {
    return `${withoutDomain}@${project ? `${project}.` : ''}${host(Cypress.env('apiUrl')).replace('api.', '')}`
}

function credentials(user: User, pw?: string) {
    const isAdmin = user == User.IAMAdminUser
    return {
        username: username(isAdmin ? user : `${user}_user_name`, isAdmin ? 'caos-ag' : 'caos-demo'),
        password: pw ? pw : Cypress.env(`${user}_password`)
    }
}

function updateCookies(newCookies: string[] | undefined, currentCookies: Map<string, string>) {
    if (newCookies === undefined) {
        return
    }
    newCookies.forEach(cs => {
        cs.split('; ').forEach(cookie => {
            const idx = cookie.indexOf('=')
            currentCookies.set(cookie.substring(0,idx), cookie.substring(idx+1))
        })
    })
}

function requestCookies(currentCookies: Map<string, string>): string[] {
    let list: Array<string> = []
    currentCookies.forEach((val, key) => {
        list.push(key+"="+val)
    })
    return list
}

export function host(url: string): string {
    return stripPort(stripProtocol(url))
}

function stripPort(s: string): string {
    const idx = s.indexOf(":")
    return idx === -1 ? s : s.substring(0,idx)
}

function stripProtocol(url: string): string {
    return url.replace('http://', '').replace('https://', '')
}

