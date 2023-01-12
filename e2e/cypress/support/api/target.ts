import { ZITADELTarget } from 'support/commands';
import { sessionAsPredefinedUser, User } from 'support/login/session';
import { ensureOrgExists, removeOrg } from './orgs';

export function newTarget(orgName: string, cleanOrg?: boolean): Cypress.Chainable<ZITADELTarget> {
  sessionAsPredefinedUser(User.IAMAdminUser);
  return cy
    .getAllSessionStorage()
    .then((storage) => {
      const baseUrlParts = Cypress.config('baseUrl').split('/');
      const origin = `${baseUrlParts[0]}//${baseUrlParts[2]}`;
      const token = <string>storage[origin]['zitadel:access_token'];
      const prunedToken = token.replace('"', '').replace('"', '');
      return prunedToken;
    })
    .then((token) => {
      return cy
        .wrap({
          headers: {
            Authorization: `Bearer ${token}`,
            'x-zitadel-orgid': undefined,
          },
          mgmtBaseURL: `${Cypress.env('BACKEND_URL')}/management/v1`,
          adminBaseURL: `${Cypress.env('BACKEND_URL')}/admin/v1`,
        })
        .then((tmpTarget) => {
          return ensureOrgExists(tmpTarget, orgName).then((dirtyOrgTarget) => {
            if (!cleanOrg) {
              return cy.wrap(dirtyOrgTarget);
            }

            return removeOrg(dirtyOrgTarget).then(() => {
              return ensureOrgExists(dirtyOrgTarget, orgName);
            });
          });
        });
    });
}

export function newOrgTarget(target: ZITADELTarget, id: number): ZITADELTarget {
  return {
    ...target,
    headers: {
      ...target.headers,
      'x-zitadel-orgid': id.toString(),
    },
  };
}
