import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API } from './types';

export function ensureProjectExists(api: API, projectName: string, orgId?: number): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/projects/_search`,
    (project: any) => project.name === projectName,
    `${api.mgmtBaseURL}/projects`,
    { name: projectName },
    orgId,
  );
}

export function ensureProjectDoesntExist(api: API, projectName: string, orgId?: number): Cypress.Chainable<null> {
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
  projectId: number,
  resourceType: ResourceType,
  resourceName: string,
  orgId?: number,
): Cypress.Chainable<null> {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/projects/${projectId}/${resourceType.resourcePath}/_search`,
    (resource: any) => resource[resourceType.compareProperty] === resourceName,
    (resource) =>
      `${api.mgmtBaseURL}/projects/${projectId}/${resourceType.resourcePath}/${resource[resourceType.identifierProperty]}`,
    orgId,
  );
}

export function ensureApplicationExists(api: API, projectId: number, appName: string): Cypress.Chainable<number> {
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

export function ensureRoleExists(api: API, projectId: number, roleKey: string): Cypress.Chainable<null> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/projects/${projectId}/roles`,
      body: {
        roleKey: roleKey,
        displayName: roleKey,
      },
      auth: {
        bearer: api.token,
      },
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
      }
      return null;
    });
}
