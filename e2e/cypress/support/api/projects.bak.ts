import { ZITADELTarget } from 'support/commands';

export function ensureProjectExists(api: ZITADELTarget, projectName: string, orgId?: number): Cypress.Chainable<number> {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/projects/_search`,
    (project: any) => project.name === projectName,
    `${api.mgmtBaseURL}/projects`,
    { name: projectName },
    orgId,
  );
}

export function ensureProjectDoesntExist(api: ZITADELTarget, projectName: string, orgId?: number): Cypress.Chainable<null> {
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
  api: ZITADELTarget,
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

export function ensureRoleExists(api: ZITADELTarget, projectId: number, roleKey: string): Cypress.Chainable<null> {
  return cy
    .request({
      method: 'POST',
      url: `${api.mgmtBaseURL}/projects/${projectId}/roles`,
      body: {
        roleKey: roleKey,
        displayName: roleKey,
      },
      headers: api.headers,
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
      }
      return null;
    });
}

function search(target: ZITADELTarget, name: string): Cypress.Chainable<any> {
  return cy
    .request({
      method: 'POST',
      url: `${target.mgmtBaseURL}/projects/_search`,
      headers: target.headers,
    })
    .then((res) => {
      return res.body?.result?.find((entity) => entity.name == name) || cy.wrap(null);
    });
}

function create(target: ZITADELTarget, name: string): Cypress.Chainable<any> {
  return cy
    .request({
      method: 'POST',
      url: `${target.mgmtBaseURL}/projects`,
      body: {
        name: name,
        allowedToFail: false,
        timeout: '10s',
      },
      failOnStatusCode: false,
      headers: target.headers,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
        return null;
      }
      return res.body;
    });
}

function remove(target: ZITADELTarget, id: number) {
  return cy
    .request({
      method: 'DELETE',
      url: `${target.mgmtBaseURL}/projects/${id}`,
      failOnStatusCode: false,
      headers: target.headers,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(404);
      }
    });
}
