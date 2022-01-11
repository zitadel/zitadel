import { apiCallProperties } from "./apiauth"

export function ensureSomethingExists(api: apiCallProperties, searchPath: string, find: (resource: any) => boolean, createPath: string, body: any): Cypress.Chainable<number> {

    return searchSomething(api, searchPath, find).then(sth => {
        if (sth) {
            return sth.id
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
    })
}

export function ensureSomethingDoesntExist(api: apiCallProperties, searchPath: string, find: (resource: any) => boolean, deletePath: (resource: any) => string): Cypress.Chainable<null> {

    return searchSomething(api, searchPath, find).then(sth => {
        if (!sth) {
            return null
        }
        return cy.request({
            method: 'DELETE',
            url: `${api.mgntBaseURL}${deletePath(sth)}`,
            headers: {
                Authorization: api.authHeader
            },
        }).then(res => {
            expect(res.status).to.equal(200)
            return null
        })
    })
}

function searchSomething(api: apiCallProperties, searchPath: string, find: (resource: any) => boolean) {

    return cy.request({
        method: 'POST',
        url: `${api.mgntBaseURL}${searchPath}`,
        headers: {
            Authorization: api.authHeader
        },
    }).then(res => {
        return res.body.result?.find(find) || null
    })    
}