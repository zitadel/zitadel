import {
  Apps,
  ensureProjectExists,
  ensureProjectResourceDoesntExist,
} from "../../support/api/projects";
import { apiAuth } from "../../support/api/apiauth";

describe("applications", () => {
  const testProjectName = "e2eprojectapplication";
  const testAppName = "e2eappundertest";

  beforeEach(`ensure it doesn't exist already`, () => {
    apiAuth().then((api) => {
      ensureProjectExists(api, testProjectName).then((projectID) => {
        ensureProjectResourceDoesntExist(
          api,
          projectID,
          Apps,
          testAppName
        ).then(() => {
          cy.visit(`/projects/${projectID}`);
        });
      });
    });
  });

  it("add app", () => {
    cy.get('[data-e2e="app-card-add"]').should("be.visible").click();
    // select webapp
    cy.get('[formcontrolname="name"]').type(testAppName);
    cy.get('[for="WEB"]').click();
    cy.get('[data-e2e="continue-button-nameandtype"]').click();
    //select authentication
    cy.get('[for="PKCE"]').click();
    cy.get('[data-e2e="continue-button-authmethod"]').click();
    //enter URL
    cy.get("cnsl-redirect-uris").eq(0).type("http://localhost:3000/api/auth/callback/zitadel");
    cy.get("cnsl-redirect-uris").eq(1).type("http://localhost:3000");
    cy.get('[data-e2e="continue-button-redirecturis"]').click();
    cy.get('[data-e2e="create-button"]')
      .click()
      .then(() => {
        cy.get("[id*=overlay]").should("exist");
      });
    cy.get(".data-e2e-success");
    cy.wait(200);
    cy.get(".data-e2e-failure", { timeout: 0 }).should("not.exist");
    //TODO: check client ID/Secret
  });

  it("should configure an application to enable dev mode")
  it("should configure an application to put user roles and info inside id token")
});
