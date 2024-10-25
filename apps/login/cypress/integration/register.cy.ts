import { stub } from "../support/mock";

describe("register", () => {
  beforeEach(() => {
    stub("zitadel.org.v2.OrganizationService", "ListOrganizations", {
      data: {
        result: [{ id: "123" }],
      },
    });
    stub("zitadel.user.v2.UserService", "AddHumanUser", {
      data: {
        userId: "123",
        email: {
          email: "john@zitadel.com",
        },
        profile: {
          givenName: "John",
          familyName: "Doe",
        },
      },
    });
  });

  it("should redirect a user who selects passwordless on register to /passkey/set", () => {
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
