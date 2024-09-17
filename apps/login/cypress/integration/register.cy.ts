import { stub } from "../support/mock";

describe("register", () => {
  beforeEach(() => {
    stub("zitadel.user.v2.UserService", "AddHumanUser", {
      data: {
        userId: "123",
      },
    });
  });

  it("should redirect a user who selects passwordless on register to /passkeys/add", () => {
    cy.visit("/register");
    cy.get('input[autocomplete="firstname"]').focus().type("John");
    cy.get('input[autocomplete="lastname"]').focus().type("Doe");
    cy.get('input[autocomplete="email"]').focus().type("john@zitadel.com");
    cy.get('input[type="checkbox"][value="privacypolicy"]').check();
    cy.get('input[type="checkbox"][value="tos"]').check();
    cy.get('button[type="submit"]').click();
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/passkey/set");
  });
});
