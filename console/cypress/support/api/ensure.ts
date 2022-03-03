import { apiCallProperties } from "./apiauth"

export function ensureSomethingExists(api: apiCallProperties, searchPath: string, find: (resource: any) => boolean, createPath: string, body: any): Cypress.Chainable<number> {

    return searchSomething(api, searchPath, find).then(sth => {
        if (sth) {
            return cy.wrap(sth.id)
        }
        return cy.request({
            method: 'POST',
            url: `${api.mgntBaseURL}${createPath}`,
            headers: {
                Authorization: api.authHeader
            },
            body: body,
            failOnStatusCode: false,
            followRedirect: false,
        }).then(res => {
            expect(res.status).to.equal(200)
            return res.body.id
        })
    }).then((id: number) => {
        awaitDesired(15, (sth) => !!sth , api, searchPath, find)
        return cy.wrap(id)
    })
}

export function ensureSomethingDoesntExist(api: apiCallProperties, searchPath: string, find: (resource: any) => boolean, deletePath: (resource: any) => string): Cypress.Chainable<null> {

    return searchSomething(api, searchPath, find).then(sth => {
        if (!sth) {
            return cy.wrap(null)
        }
        return cy.request({
            method: 'DELETE',
            url: `${api.mgntBaseURL}${deletePath(sth)}`,
            headers: {
                Authorization: api.authHeader
            },
            failOnStatusCode: false
        }).then((res) => {
            expect(res.status).to.equal(200)
        })
    }).then(() => {
        awaitDesired(15, (sth) => !sth , api, searchPath, find)
        return null
    })
}

function searchSomething(api: apiCallProperties, searchPath: string, find: (resource: any) => boolean) {

    return cy.request({
        method: 'POST',
        url: `${api.mgntBaseURL}${searchPath}`,
        headers: {
            Authorization: api.authHeader
        },
    }).then(res => res.body.result?.find(find) || null)
}

function awaitDesired(trials: number, expectSth: (sth: any) => boolean, api: apiCallProperties, searchPath: string, find: (resource: any) => boolean) {
    searchSomething(api, searchPath, find).then(sth => {
        if (!expectSth(sth)) {
            expect(trials, `trying ${trials} more times`).to.be.greaterThan(0);
            cy.wait(1000)
            awaitDesired(trials - 1, expectSth, api, searchPath, find)
        }            
    })
}