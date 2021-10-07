import { apiCallProperties } from "./apiauth"

export function ensureAppDoesntExist(apiCallProperties: apiCallProperties, projectId: number, appName: string): Cypress.Chainable<null> {
    return cy.request({
        method: 'POST',
        url: `${apiCallProperties.mgntBaseURL}projects/${projectId}/apps/_search`,
        headers: {
            Authorization: apiCallProperties.authHeader
        },
    }).then(appsRes => {
        const appId = appsRes.body.result?.find(app => app.name === appName).id
        if (appId) {
            return cy.request({
                method: 'DELETE',
                url: `${apiCallProperties.mgntBaseURL}projects/${projectId}/apps/${appId}`,
                headers: {
                    Authorization: apiCallProperties.authHeader
                },
            }).then(()=> null)
        }
        return null
    })
}