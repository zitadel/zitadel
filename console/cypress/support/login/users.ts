export enum User {
    OrgOwner = 'org_owner',
    OrgOwnerViewer = 'org_owner_viewer',
    OrgProjectCreator = 'org_project_creator',
}

export function login(user:User): void {
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
}

function credentials(user: User) {
    return {
        username: `${user}_user_name@caos-demo.${Cypress.env('apiCallsDomain')}`,
        password: Cypress.env(`${user}_password`)
    }
}

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
