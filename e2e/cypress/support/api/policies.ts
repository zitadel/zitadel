import { ZITADELTarget } from 'support/commands';

export function ensureDomainPolicy(
  target: ZITADELTarget,
  userLoginMustBeDomain: boolean,
  validateOrgDomains: boolean,
  smtpSenderAddressMatchesInstanceDomain: boolean,
): Cypress.Chainable<null> {
  resetDomainPolicy(target);
  setDomainPolicy(target, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain);

  return null;
}

function resetDomainPolicy(target: ZITADELTarget) {
  return cy
    .request({
      method: 'DELETE',
      url: `${target.adminBaseURL}/orgs/${target.headers['x-zitadel-orgid']}/policies/domain`,
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
      url: `${target.adminBaseURL}/orgs/${target.headers['x-zitadel-orgid']}/policies/domain`,
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
