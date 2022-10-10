import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API } from './types';

export function ensureProjectExists(api: API, projectName: string, orgId?: number): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgntBaseURL}projects/_search`,
    (project: any) => project.name === projectName,
    `${api.mgntBaseURL}projects`,
    { name: projectName },
    orgId,
  );
}

export function ensureProjectDoesntExist(api: API, projectName: string, orgId?: number): Cypress.Chainable<null> {
  return ensureItemDoesntExist(
    api,
    `${api.mgntBaseURL}projects/_search`,
    (project: any) => project.name === projectName,
    (project) => `${api.mgntBaseURL}projects/${project.id}`,
    orgId,
  );
}

class ResourceType {
  constructor(public resourcePath: string, public compareProperty: string, public identifierProperty: string) {}
}

export const Apps = new ResourceType('apps', 'name', 'id');
export const Roles = new ResourceType('roles', 'key', 'key');
//export const Grants = new ResourceType('apps', 'name')

export function ensureProjectResourceDoesntExist(
  api: API,
  projectId: number,
  resourceType: ResourceType,
  resourceName: string,
  orgId?: number,
): Cypress.Chainable<null> {
  return ensureItemDoesntExist(
    api,
    `${api.mgntBaseURL}projects/${projectId}/${resourceType.resourcePath}/_search`,
    (resource: any) => resource[resourceType.compareProperty] === resourceName,
    (resource) =>
      `${api.mgntBaseURL}projects/${projectId}/${resourceType.resourcePath}/${resource[resourceType.identifierProperty]}`,
    orgId,
  );
}

export function ensureApplicationExists(api: API, projectId: number, appName: string): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgntBaseURL}projects/${projectId}/${Apps.resourcePath}/_search`,
    (resource: any) => resource.name === appName,
    `${api.mgntBaseURL}projects/${projectId}/${Apps.resourcePath}/oidc`,
    {
      name: appName,
      redirectUris: ['https://e2eredirecturl.org'],
      responseTypes: ['OIDC_RESPONSE_TYPE_CODE'],
      grantTypes: ['OIDC_GRANT_TYPE_AUTHORIZATION_CODE'],
      authMethodType: 'OIDC_AUTH_METHOD_TYPE_NONE',
      postLogoutRedirectUris: ['https://e2elogoutredirecturl.org'],
    },
  );
}
