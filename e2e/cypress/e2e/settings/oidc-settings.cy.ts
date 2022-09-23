import {apiAuth} from '../../support/api/apiauth';
import {ensureOIDCSettingsSet} from "../../support/api/oidc-settings";

describe('oidc settings', () => {
    const oidcSettingsPath = `/settings?id=oidc`;

    before(`ensure it is set`, () => {
        apiAuth().then((apiCallProperties) => {
            ensureOIDCSettingsSet(apiCallProperties, 1 * 60 * 60, 2 * 60 * 60, 2 * 24 * 60 * 60, 7 * 24 * 60 * 60);
            cy.visit(oidcSettingsPath);
        });
    });

    it(`should update oidc settings`, () => {
        cy.get('[formcontrolname="accessTokenLifetime"]').should('value', 1);
        cy.get('[formcontrolname="accessTokenLifetime"]').clear().type("1");
        cy.get('[formcontrolname="idTokenLifetime"]').clear().type("24");
        cy.get('[formcontrolname="refreshTokenExpiration"]').clear().type("30");
        cy.get('[formcontrolname="refreshTokenIdleExpiration"]').clear().type("7");
        cy.get('[data-e2e="save-button"]').click();
        cy.get('.data-e2e-success');
    });
});
