import { apiCallProperties } from "./apiauth"
import { ensureSomethingDoesntExist, ensureSomethingExists } from "./ensure"

export function ensureHumanUserExists(api: apiCallProperties, username: string): Cypress.Chainable<number> {

    return ensureSomethingExists(
        api,
        'users/_search',
        (user: any) => user.userName === username,
        'users/human',
        {
            user_name: username,
            profile: {
                first_name: 'e2efirstName',
                last_name: 'e2elastName',
            },
            email: { 
                email: 'e2e@email.ch',
            },
            phone: {
                phone: '+41 123456789',
        },
    })
}

export function ensureMachineUserExists(api: apiCallProperties, username: string): Cypress.Chainable<number> {
    
    return ensureSomethingExists(
        api,
        'users/_search',
        (user: any) => user.userName === username,
        'users/machine',
        {
            user_name: username,
            name: 'e2emachinename',
            description: 'e2emachinedescription',
        },
    )
}

export function ensureUserDoesntExist(api: apiCallProperties, username: string): Cypress.Chainable<null> {

    return ensureSomethingDoesntExist(
        api,
        'users/_search',
        (user: any) => user.userName === username,
        (user) => `users/${user.id}`
    )
}
