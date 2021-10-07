import { apiCallProperties } from "./apiauth"

export function ensureProjectExists(apiCallProperties: apiCallProperties, projectName: string): Cypress.Chainable<number> {
    return getProjectID(apiCallProperties, projectName).then(projectId => {
        if (projectId) {
            return cy.wrap(projectId)
        }

        return cy.request({
            method: 'POST',
            url: `${apiCallProperties.mgntBaseURL}projects`,
            headers: {
                Authorization: apiCallProperties.authHeader
            }, 
            body: {
                name: projectName,
            },
            failOnStatusCode: false,
            followRedirect: false
        }).then(res => {
            expect(res.status).to.equal(200)
            return res.body.id
        })
    })
}

export function ensureProjectDoesntExist(apiCallProperties: apiCallProperties, projectName: string): Cypress.Chainable<null> {
    return getProjectID(apiCallProperties, projectName).then(projectId => {
        debugger
        if (!projectId) {
            return cy.wrap(null)
        }

        return cy.request({
            method: 'DELETE',
            url: `${apiCallProperties.mgntBaseURL}projects/${projectId}`,
            headers: {
                Authorization: apiCallProperties.authHeader
            },
        }).then(() => null)
    })
}


export function getProjectID(apiCallProperties: apiCallProperties, projectName: string): Cypress.Chainable<number> {
    return cy.request({
        method: 'POST',
        url: `${apiCallProperties.mgntBaseURL}projects/_search`,
        headers: {
            Authorization: apiCallProperties.authHeader
        },
    }).then(projectsRes => {
        return projectsRes.body.result?.find(project => project.name === projectName)?.id || null
    })
}