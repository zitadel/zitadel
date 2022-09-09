import { apiAuth } from "../../support/api/apiauth";
import {
  ensureProjectDoesntExist,
  ensureProjectExists,
} from "../../support/api/projects";

describe("projects", () => {

  beforeEach(() => {
    apiAuth().as("api")
  })

  const testProjectNameCreate = "e2eprojectcreate";
  const testProjectNameDeleteList = "e2eprojectdeletelist";
  const testProjectNameDeleteGrid = "e2eprojectdeletegrid";

  describe("add project", () => {
    beforeEach(`ensure it doesn't exist already`, function() {
      ensureProjectDoesntExist(this.api, testProjectNameCreate);
      cy.visit(`/projects`);
    });

    it("should add a project", () => {
      cy.get(".add-project-button")
        .click({ force: true });
      cy.get("input")
        .type(testProjectNameCreate);
      cy.get('[data-e2e="continue-button"]')
        .click();
      cy.get(".data-e2e-success");
      cy.wait(200);
      cy.get(".data-e2e-failure", { timeout: 0 })
        .should("not.exist");
    });

    it("should configure a project to assert roles on authentication")
  });

  describe("edit project", () => {
    beforeEach("ensure it exists", function() {
      ensureProjectExists(this.api, testProjectNameDeleteList);
      cy.visit(`/projects`);
    });

    describe("remove project", () => {

      beforeEach("ensure it exists", function() {
        ensureProjectExists(this.api, testProjectNameDeleteGrid);
        cy.visit(`/projects`);
      });

      it("removes the project from list view", () => {
        cy.get('[data-e2e="toggle-grid"]')
          .click();
        cy.get('[data-e2e="timestamp"]');
        cy.contains("tr", testProjectNameDeleteList, { timeout: 1000 })
          .find('[data-e2e="delete-project-button"]')
          .click({ force: true });
        cy.get('[data-e2e="confirm-dialog-input"]')
          .focus()
          .type(testProjectNameDeleteList);
        cy.get('[data-e2e="confirm-dialog-button"]')
          .click();
        cy.get(".data-e2e-success");
        cy.wait(200);
        cy.get(".data-e2e-failure", { timeout: 0 })
          .should("not.exist");
      });

      it("removes the project from grid view", () => {
        cy.contains('[data-e2e="grid-card"]', testProjectNameDeleteGrid)
          .find('[data-e2e="delete-project-button"]')
          .click({force: true});
        cy.get('[data-e2e="confirm-dialog-input"]')
          .focus()
          .type(testProjectNameDeleteGrid);
        cy.get('[data-e2e="confirm-dialog-button"]')
          .click();
        cy.get(".data-e2e-success");
        cy.wait(200);
        cy.get(".data-e2e-failure", { timeout: 0 })
          .should("not.exist");
      });
      });

      it("should add a project manager")
      it("should remove a project manager")
    })
});
