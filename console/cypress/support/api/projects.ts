import { apiCallProperties } from "./apiauth"
import { ensureSomethingDoesntExist, ensureSomethingExists } from "./ensure"

export function ensureProjectExists(api: apiCallProperties, projectName: string): Cypress.Chainable<number> {
    
    return ensureSomethingExists(
        api,
        `projects/_search`,
        (project: any) => project.name === projectName,
        'projects',
        { name: projectName },
    )
}

export function ensureProjectDoesntExist(api: apiCallProperties, projectName: string): Cypress.Chainable<null> {
    
    return ensureSomethingDoesntExist(
        api,
        `projects/_search`,
        (project: any) => project.name === projectName,
        (project) => `projects/${project.id}`,
    )
}

class ResourceType {
    constructor(
        public resourcePath: string,
        public compareProperty: string,
        public identifierProperty: string,
    ){}
}

export const Apps = new ResourceType('apps', 'name', 'id')
export const Roles = new ResourceType('roles', 'key', 'key')
//export const Grants = new ResourceType('apps', 'name')


export function ensureProjectResourceDoesntExist(api: apiCallProperties, projectId: number, resourceType: ResourceType, resourceName: string): Cypress.Chainable<null> {
    return ensureSomethingDoesntExist(
        api,
        `projects/${projectId}/${resourceType.resourcePath}/_search`,
        (resource: any) => {
            return resource[resourceType.compareProperty] === resourceName
        },
        (resource) => {
            return `projects/${projectId}/${resourceType.resourcePath}/${resource[resourceType.identifierProperty]}`
        }
    )
}

export function ensureApplicationExists(api: apiCallProperties, projectId: number, appName: string): Cypress.Chainable<number> {

    return ensureSomethingExists(
        api,
        `projects/${projectId}/${Apps.resourcePath}/_search`,
        (resource: any) => resource.name === appName,
        `projects/${projectId}/${Apps.resourcePath}/oidc`,
        {
            name: appName,
            redirectUris: [
                'https://e2eredirecturl.org'
            ],
            responseTypes: [
                "OIDC_RESPONSE_TYPE_CODE"
            ],
            grantTypes: [
                "OIDC_GRANT_TYPE_AUTHORIZATION_CODE"
            ],
            authMethodType: "OIDC_AUTH_METHOD_TYPE_NONE",
            postLogoutRedirectUris: [
                'https://e2elogoutredirecturl.org'
            ],
/*            "clientId": "129383004379407963@e2eprojectpermission",
            "clockSkew": "0s",
            "allowedOrigins": [
                "https://testurl.org"
            ]*/
        },
    )
}
