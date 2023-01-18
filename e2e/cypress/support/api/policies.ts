import { ZITADELTarget } from 'support/commands';
import { standardCreate, standardRemove } from './standard';

export function ensureDomainPolicy(
  target: ZITADELTarget,
  userLoginMustBeDomain: boolean,
  validateOrgDomains: boolean,
  smtpSenderAddressMatchesInstanceDomain: boolean,
) {
  resetDomainPolicy(target);
  setDomainPolicy(target, userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain);
  return getDomainPolicy(target).should(
    (res) =>
      res.body.userLoginMustBeDomain == userLoginMustBeDomain &&
      res.body.validateOrgDomains == validateOrgDomains &&
      res.body.smtpSenderAddressMatchesInstanceDomain == smtpSenderAddressMatchesInstanceDomain,
  );
}

function resetDomainPolicy(target: ZITADELTarget) {
  return standardRemove(target, `${target.adminBaseURL}/orgs/${target.orgId}/policies/domain`);
}

function setDomainPolicy(
  target: ZITADELTarget,
  userLoginMustBeDomain: boolean,
  validateOrgDomains: boolean,
  smtpSenderAddressMatchesInstanceDomain: boolean,
) {
  return standardCreate(
    target,
    `${target.adminBaseURL}/orgs/${target.orgId}/policies/domain`,
    {
      userLoginMustBeDomain: userLoginMustBeDomain,
      validateOrgDomains: validateOrgDomains,
      smtpSenderAddressMatchesInstanceDomain: smtpSenderAddressMatchesInstanceDomain,
    },
    'no id',
  );
}

function getDomainPolicy(target: ZITADELTarget) {
  return cy.request({
    method: 'GET',
    url: `${target.adminBaseURL}/orgs/${target.orgId}/policies/domain`,
    headers: target.headers,
  });
}
