import {
  Apps,
  ensureProjectExists,
  ensureProjectResourceDoesntExist,
} from "../../support/api/projects";
import { apiAuth } from "../../support/api/apiauth";

describe("applications", () => {
  const testProjectName = "e2eprojectapplication";
  const testAppName = "e2eappundertest";

  beforeEach(() => {
    apiAuth()
      .as("api")
      .then(api => {
        ensureProjectExists(api, testProjectName)
          .as("projectId")
      })
  })

  describe("add app", function() {
    beforeEach(`ensure it doesn't exist already`, function() {
      ensureProjectResourceDoesntExist(
        this.api,
        this.projectId,
        Apps,
        testAppName
      )
      cy.visit(`/projects/${this.projectId}`);
    });

    it("add app", () => {
      cy.get('[data-e2e="app-card-add"]')
        .should("be.visible")
        .click();
      cy.get('[formcontrolname="name"]')
        .type(testAppName);
      cy.get('[for="WEB"]')
        .click();
      cy.get('[data-e2e="continue-button-nameandtype"]')
        .click();
      cy.get('[for="PKCE"]')
        .click();
      cy.get('[data-e2e="continue-button-authmethod"]')
        .click();
      cy.get("cnsl-redirect-uris")
        .eq(0)
        .type("http://localhost:3000/api/auth/callback/zitadel");
      cy.get("cnsl-redirect-uris")
        .eq(1)
        .type("http://localhost:3000");
      cy.get('[data-e2e="continue-button-redirecturis"]')
        .click();
      cy.get('[data-e2e="create-button"]')
        .click()
        .then(() => {
          cy.get("[id*=overlay]")
            .should("exist");
        });
      cy.get(".data-e2e-success");
      cy.wait(200);
      cy.get(".data-e2e-failure", { timeout: 0 })
        .should("not.exist");
      //TODO: check client ID/Secret
    });
  })

  describe("edit app", () => {
    it("should configure an application to enable dev mode")
    it("should configure an application to put user roles and info inside id token")
  })
});
