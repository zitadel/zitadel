import { login, User } from 'support/login/users'

export interface apiCallProperties {
    authHeader: string
    mgntBaseURL: string
}

export function apiAuth(): Cypress.Chainable<apiCallProperties> {
    return login(User.IAMAdminUser, 'Password1!', false, true).then(token => {
        return <apiCallProperties>{
            authHeader: `Bearer ${token}`,
            mgntBaseURL: `/management/v1/`,
        }
    })
}
