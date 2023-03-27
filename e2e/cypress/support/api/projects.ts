import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API, Entity } from './types';

export function ensureProjectExists(api: API, projectName: string, orgId?: string) {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/projects/_search`,
    (project: any) => project.name === projectName,
    `${api.mgmtBaseURL}/projects`,
    { name: projectName },
    orgId,
  );
}

export function ensureProjectDoesntExist(api: API, projectName: string, orgId?: string) {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/projects/_search`,
    (project: any) => project.name === projectName,
    (project) => `${api.mgmtBaseURL}/projects/${project.id}`,
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
  projectId: string,
  resourceType: ResourceType,
  resourceName: string,
  orgId?: string,
): Cypress.Chainable<null> {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/projects/${projectId}/${resourceType.resourcePath}/_search`,
    (resource: Entity) => resource[resourceType.compareProperty] === resourceName,
    (resource: Entity) =>
      `${api.mgmtBaseURL}/projects/${projectId}/${resourceType.resourcePath}/${resource[resourceType.identifierProperty]}`,
    orgId,
  );
}

export function ensureApplicationExists(api: API, projectId: number, appName: string) {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/projects/${projectId}/${Apps.resourcePath}/_search`,
    (resource: any) => resource.name === appName,
    `${api.mgmtBaseURL}/projects/${projectId}/${Apps.resourcePath}/oidc`,
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
