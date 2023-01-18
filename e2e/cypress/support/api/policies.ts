import { ZITADELTarget } from 'support/commands';

export function ensureDomainPolicy(
  target: ZITADELTarget,
  userLoginMustBeDomain: boolean,
  validateOrgDomains: boolean,
  smtpSenderAddressMatchesInstanceDomain: boolean,
): Cypress.Chainable<null> {
  resetDomainPolicy(target);
  setDomainPolicy(target, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain);

  for (let i = 0; i < 10; i++) {
    getDomainPolicy(target).should((res) => {
      res.body.userLoginMustBeDomain == userLoginMustBeDomain &&
        res.body.validateOrgDomains == validateOrgDomains &&
        res.body.smtpSenderAddressMatchesInstanceDomain == smtpSenderAddressMatchesInstanceDomain;
    });
  }

  return null;
}

function resetDomainPolicy(target: ZITADELTarget) {
  return cy
    .request({
      method: 'DELETE',
      url: `${target.adminBaseURL}/orgs/${target.orgId}/policies/domain`,
      headers: target.headers,
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(404);
      }
      return res;
    });
}

function setDomainPolicy(
  target: ZITADELTarget,
  userLoginMustBeDomain: boolean,
  validateOrgDomains: boolean,
  smtpSenderAddressMatchesInstanceDomain: boolean,
): Cypress.Chainable<Cypress.Response<any>> {
  return cy
    .request({
      method: 'POST',
      url: `${target.adminBaseURL}/orgs/${target.orgId}/policies/domain`,
      body: {
        userLoginMustBeDomain: userLoginMustBeDomain,
        validateOrgDomains: validateOrgDomains,
        smtpSenderAddressMatchesInstanceDomain: smtpSenderAddressMatchesInstanceDomain,
      },
      headers: target.headers,
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(409);
      }
      return res;
    });
}

function getDomainPolicy(target: ZITADELTarget) {
  return cy.request({
    method: 'GET',
    url: `${target.adminBaseURL}/orgs/${target.orgId}/policies/domain`,
    headers: target.headers,
  });
}
