import { ZITADELTarget } from 'support/commands';

export function ensureHumanIsOrgMember(target: ZITADELTarget, userId: number, roles: string[]) {
  return cy.request({
    method: 'POST',
    url: `${target.mgmtBaseURL}/orgs/me/members`,
    body: {
      userId: userId,
      roles: roles,
    },
    headers:target.headers,
    failOnStatusCode: false,
  }).then(res => {
    if (!res.isOkStatusCode){
      expect(res.status).to.equal(409)
    }
  })
}

export function ensureHumanIsNotOrgMember(target: ZITADELTarget, userId: number) {
  return cy.request({
    method: 'DELETE',
    url: `${target.mgmtBaseURL}/orgs/me/members/${userId}`,
    headers:target.headers,
    failOnStatusCode: false,
  }).then(res => {
    if (!res.isOkStatusCode){
      expect(res.status).to.equal(404)
    }
  })
}

export function ensureHumanIsProjectMember(target: ZITADELTarget, projectId: number, userId: number, roles: string[], grantId?: number) {
  return cy.request({
    method: 'POST',
    url: `${target.mgmtBaseURL}/projects/${projectId}${grantId ? `/grants/${grantId}` : ''}/members`,
    body: {
      userId: userId,
      roles: roles,
    },
    headers:target.headers,
    failOnStatusCode: false,
  }).then(res => {
    if (!res.isOkStatusCode){
      expect(res.status).to.equal(409)
    }
  })
}

export function ensureHumanIsNotProjectMember(target: ZITADELTarget, projectId: number, userId: number, grantId?: number) {
  return cy.request({
    method: 'DELETE',
    url: `${target.mgmtBaseURL}/projects/${projectId}${grantId ? `grants/${grantId}/` : ''}/members/${userId}`,
    headers:target.headers,
    failOnStatusCode: false,
  }).then(res => {
    if (!res.isOkStatusCode){
      expect(res.status).to.equal(404)
    }
  })
}
