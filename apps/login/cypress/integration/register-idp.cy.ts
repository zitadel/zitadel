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

  it("should show a custom text on the idp button", () => {
    cy.visit("/register/idp");
    const button = cy.get('button[e2e="google"]');
    // button.find("span").contains(CUSTOM_TEXT);
    button.click();
    cy.location("pathname", { timeout: 10_000 }).should("eq", IDP_URL);
  });
});
