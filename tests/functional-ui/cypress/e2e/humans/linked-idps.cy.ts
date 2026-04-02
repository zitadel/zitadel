import { Context } from 'support/commands';
import { createHumanUser } from '../../support/api/users';
import { addUserIDPLink, createJWTIDP, deleteIDP, waitForLinkedIDPCount } from '../../support/api/idps';

describe('human linked idps', () => {
  let userId = '';
  let idpId = '';

  const getEmailEditButton = () => {
    return cy.get('cnsl-contact .contact-method-row').first().find('.right button');
  };

  beforeEach(() => {
    const suffix = Date.now().toString();
    const username = `e2elinkedidp${suffix}`;
    const idpName = `e2e-jwt-idp-${suffix}`;

    cy.context().as('ctx');

    cy.get<Context>('@ctx').then((ctx) => {
      createJWTIDP(ctx.api, idpName).then((createdIDPId) => {
        idpId = createdIDPId;

        createHumanUser(ctx.api, username).then((response) => {
          expect(response.status).to.equal(200);
          userId = response.body.userId;

          addUserIDPLink(ctx.api, userId, idpId, `external-${suffix}`, `${username}@external.example`).then(() => {
            waitForLinkedIDPCount(ctx.api, userId, 1);
          });
        });
      });
    });
  });

  afterEach(() => {
    cy.get<Context>('@ctx').then((ctx) => {
      if (userId) {
        cy.request({
          method: 'DELETE',
          url: `${ctx.api.mgmtBaseURL}/users/${userId}`,
          headers: {
            Authorization: `Bearer ${ctx.api.token}`,
          },
          failOnStatusCode: false,
        });
      }

      if (idpId) {
        deleteIDP(ctx.api, idpId);
      }
    });
  });

  it('should re-enable email editing after removing the last linked idp', () => {
    cy.visit(`/users/${userId}?id=general`);

    cy.get('[data-e2e="sidenav-element-general"]').click();
    cy.get('cnsl-contact').should('be.visible');
    getEmailEditButton().should('be.disabled');

    cy.get('[data-e2e="sidenav-element-idp"]').click();
    cy.get('cnsl-external-idps table').should('be.visible');
    cy.get('cnsl-external-idps tbody tr')
      .first()
      .within(() => {
        cy.get('button[color="warn"]').click({ force: true });
      });
    cy.get('[data-e2e="confirm-dialog-button"]').click();

    cy.get('cnsl-external-idps .no-content-row', { timeout: 20_000 }).should('be.visible');

    cy.get('[data-e2e="sidenav-element-general"]').click();
    cy.get('cnsl-contact').should('be.visible');
    getEmailEditButton().should('be.enabled');
  });
});
