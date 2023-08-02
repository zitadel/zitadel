import { stub } from "../support/mock";

const IDP_URL = "https://google.com";

describe("register idps", () => {
  beforeEach(() => {
    stub("zitadel.user.v2alpha.UserService", "StartIdentityProviderFlow", {
      data: {
        authUrl: IDP_URL,
      },
    });
  });

  it("should redirect the user to the correct url", () => {
    cy.visit("/register/idp");
    const button = cy.get('button[e2e="google"]');
    button.click();
    cy.location("href", { timeout: 10_000 }).should("eq", IDP_URL);
  });
});
