import { apiAuth } from "../../support/api/apiauth";
import {
  ensureMachineUserExists,
  ensureUserDoesntExist,
} from "../../support/api/users";
import { loginname } from "../../support/login/users";

describe("machines", () => {
  const machinesPath = `/users?type=machine`;
  const testMachineUserNameAdd = "e2emachineusernameadd";
  const testMachineUserNameRemove = "e2emachineusernameremove";

  describe("add", () => {
    before(`ensure it doesn't exist already`, () => {
      apiAuth().then((apiCallProperties) => {
        ensureUserDoesntExist(apiCallProperties, testMachineUserNameAdd).then(
          () => {
            cy.visit(machinesPath);
          }
        );
      });
    });

    it("should add a machine", () => {
      cy.get('[data-e2e="action-key-add"]')
        .parents('[data-e2e="create-user-button"]')
        .click();
      cy.url().should("contain", "users/create-machine");
      //force needed due to the prefilled username prefix
      cy.get('[formcontrolname="userName"]')
        .focus()
        .type(testMachineUserNameAdd);
      cy.get('[formcontrolname="name"]')
      .focus()
      .type("e2emachinename");
      cy.get('[formcontrolname="description"]')
      .focus()
      .type("e2emachinedescription");
      cy.get('[data-e2e="create-button"]').click();
      cy.get(".data-e2e-success");
      cy.wait(200);
      cy.get(".data-e2e-failure", { timeout: 0 }).should("not.exist");
    });
  });

  describe("remove", () => {
    before("ensure it exists", () => {
      apiAuth().then((api) => {
        ensureMachineUserExists(api, testMachineUserNameRemove).then(() => {
          cy.visit(machinesPath);
        });
      });
    });

    it("should delete a machine", () => {
      cy.contains("tr", testMachineUserNameRemove, { timeout: 1000 })
        // doesn't work, need to force click.
        // .trigger('mouseover')
        .find('[data-e2e="enabled-delete-button"]')
        .click({force: true});
      cy.get('[data-e2e="confirm-dialog-input"]')
        .focus()
        .type(loginname(testMachineUserNameRemove, Cypress.env("ORGANIZATION")));
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.get(".data-e2e-success");
      cy.wait(200);
      cy.get(".data-e2e-failure", { timeout: 0 }).should("not.exist");
    });
  });
});
