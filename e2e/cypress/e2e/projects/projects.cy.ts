import { apiAuth } from "../../support/api/apiauth";
import {
  ensureProjectDoesntExist,
  ensureProjectExists,
} from "../../support/api/projects";

describe("projects", () => {
  const testProjectNameCreate = "e2eprojectcreate";
  const testProjectNameDeleteList = "e2eprojectdeletelist";
  const testProjectNameDeleteGrid = "e2eprojectdeletegrid";

  describe("add project", () => {
    beforeEach(`ensure it doesn't exist already`, () => {
      apiAuth().then((api) => {
        ensureProjectDoesntExist(api, testProjectNameCreate);
      });
      cy.visit(`/projects`);
    });

    it("should add a project", () => {
      cy.get(".add-project-button").click({ force: true });
      cy.get("input").type(testProjectNameCreate);
      cy.get('[data-e2e="continue-button"]').click();
      cy.get(".data-e2e-success");
      cy.wait(200);
      cy.get(".data-e2e-failure", { timeout: 0 }).should("not.exist");
    });
  });

  describe.skip("remove project", () => {
    describe("list view", () => {
      beforeEach("ensure it exists", () => {
        apiAuth().then((api) => {
          ensureProjectExists(api, testProjectNameDeleteList);
        });
        cy.visit(`/projects`);
      });

      it("removes the project", () => {
        cy.get('[data-e2e="toggle-grid"]').click();
        cy.get('[data-e2e="timestamp"]');
        cy.contains("tr", testProjectNameDeleteList, { timeout: 1000 })
          .find('[data-e2e="delete-project-button"]')
          .click({ force: true });
        cy.get('[data-e2e="confirm-dialog-input"]').type(
          testProjectNameDeleteList
        );
        cy.get('[data-e2e="confirm-dialog-button"]').click();
        cy.get(".data-e2e-success");
        cy.wait(200);
        cy.get(".data-e2e-failure", { timeout: 0 }).should("not.exist");
      });
    });

    describe("grid view", () => {
      beforeEach("ensure it exists", () => {
        apiAuth().then((api) => {
          ensureProjectExists(api, testProjectNameDeleteGrid);
        });
        cy.visit(`/projects`);
      });

      it("removes the project", () => {
        cy.contains('[data-e2e="grid-card"]', testProjectNameDeleteGrid)
          .find('[data-e2e="delete-project-button"]')
          .trigger("mouseover")
          .click();
        cy.get('[data-e2e="confirm-dialog-input"]').type(
          testProjectNameDeleteGrid
        );
        cy.get('[data-e2e="confirm-dialog-button"]').click();
        cy.get(".data-e2e-success");
        cy.wait(200);
        cy.get(".data-e2e-failure", { timeout: 0 }).should("not.exist");
      });
    });
  });
});
