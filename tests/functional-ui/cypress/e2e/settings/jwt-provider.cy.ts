import { apiAuth } from '../../support/api/apiauth';

describe('jwt provider', () => {
  const jwtProviderPath = `/instance/provider/jwt/create`;

  beforeEach(`visit the jwt provider form`, () => {
    apiAuth().then(() => {
      cy.visit(jwtProviderPath);
    });
  });

  it(`should warn about the implication of leaving the audience empty`, () => {
    cy.get('[formcontrolname="audience"]').should('have.value', '');
    cy.get('[data-e2e="audience-empty-warning"]').should('be.visible');
  });

  it(`should hide the warning once an audience is set`, () => {
    cy.get('[data-e2e="audience-empty-warning"]').should('be.visible');
    cy.get('[formcontrolname="audience"]').should('be.enabled').type('my-client');
    cy.get('[data-e2e="audience-empty-warning"]').should('not.exist');
  });
});
