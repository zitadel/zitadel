import { apiCallProperties } from "./apiauth"


export enum Policy {
    Label = "label"
}

export function resetPolicy(api: apiCallProperties, policy: Policy) {
    cy.request({
        method: 'DELETE',
        url: `${api.mgntBaseURL}/policies/${policy}`,
        headers: {
            Authorization: api.authHeader
        },
    }).then(res => {
        expect(res.status).to.equal(200)
        return null
    })    
}