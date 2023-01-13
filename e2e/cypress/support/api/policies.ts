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

/*
export function legacyEnsureDomainPolicy(
  api: API,
  userLoginMustBeDomain: boolean,
  validateOrgDomains: boolean,
  smtpSenderAddressMatchesInstanceDomain: boolean,
): Cypress.Chainable<number> {
  return ensureSetting(
    api,
    `${api.adminBaseURL}/policies/domain`,
    (body: any) => {
      const result = {
        sequence: parseInt(<string>body.policy?.details?.sequence),
        id: body.policy?.details?.resourceOwner,
        entity: null,
      };
      if (
        body.policy &&
        (body.policy.userLoginMustBeDomain ? body.policy.userLoginMustBeDomain : false) == userLoginMustBeDomain &&
        (body.policy.validateOrgDomains ? body.policy.validateOrgDomains : false) == validateOrgDomains &&
        (body.policy.smtpSenderAddressMatchesInstanceDomain ? body.policy.smtpSenderAddressMatchesInstanceDomain : false) ==
          smtpSenderAddressMatchesInstanceDomain
      ) {
        return { ...result, entity: body.policy };
      }
      return result;
    },
    `${api.adminBaseURL}/policies/domain`,
    {
      userLoginMustBeDomain: userLoginMustBeDomain,
      validateOrgDomains: validateOrgDomains,
      smtpSenderAddressMatchesInstanceDomain: smtpSenderAddressMatchesInstanceDomain,
    },
  );
}
*/
