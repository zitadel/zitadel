import { stub } from "../support/e2e";

const IDP_URL = "https://example.com/idp/url";

describe("register idps", () => {
  beforeEach(() => {
    stub("zitadel.user.v2.UserService", "StartIdentityProviderIntent", {
      data: {
        authUrl: IDP_URL,
      },
    });
  });

  it("should redirect the user to the correct url", () => {
    cy.visit("/idp");
    cy.get('button[e2e="google"]').click();
    cy.origin(IDP_URL, { args: IDP_URL }, (url) => {
      cy.location("href").should("eq", url);
    });
  });
});
