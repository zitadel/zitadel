import { stub } from "../support/mock";

const CUSTOM_TEXT = "Hubba Bubba";
const IDP_URL = "https://google.com";

describe("register idps", () => {
  beforeEach(() => {
    stub(
      "zitadel.settings.v2alpha.SettingsService",
      "GetActiveIdentityProviders",
      {
        data: {
          identityProviders: [
            {
              id: "123",
              name: CUSTOM_TEXT,
              type: 10,
            },
          ],
        },
      }
    );
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
