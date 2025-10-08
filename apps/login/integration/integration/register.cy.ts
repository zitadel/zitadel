import { stub } from "../support/e2e";

describe("register", () => {
  beforeEach(() => {
    stub("zitadel.org.v2.OrganizationService", "ListOrganizations", {
      data: {
        details: {
          totalResult: 1,
        },
        result: [{ id: "256088834543534543" }],
      },
    });
    stub("zitadel.settings.v2.SettingsService", "GetLoginSettings", {
      data: {
        settings: {
          passkeysType: 1,
          allowRegister: true,
          allowUsernamePassword: true,
          defaultRedirectUri: "",
        },
      },
    });
    stub("zitadel.user.v2.UserService", "AddHumanUser", {
      data: {
        userId: "221394658884845598",
      },
    });
    stub("zitadel.session.v2.SessionService", "CreateSession", {
      data: {
        details: {
          sequence: 859,
          changeDate: new Date("2024-04-04T09:40:55.577Z"),
          resourceOwner: "220516472055706145",
        },
        sessionId: "221394658884845598",
        sessionToken: "SDMc7DlYXPgwRJ-Tb5NlLqynysHjEae3csWsKzoZWLplRji0AYY3HgAkrUEBqtLCvOayLJPMd0ax4Q",
        challenges: undefined,
      },
    });

    stub("zitadel.session.v2.SessionService", "GetSession", {
      data: {
        session: {
          id: "221394658884845598",
          creationDate: new Date("2024-04-04T09:40:55.577Z"),
          changeDate: new Date("2024-04-04T09:40:55.577Z"),
          sequence: 859,
          factors: {
            user: {
              id: "221394658884845598",
              loginName: "john@example.com",
            },
            password: undefined,
            webAuthN: undefined,
            intent: undefined,
          },
          metadata: {},
        },
      },
    });
  });

  it("should redirect a user who selects passwordless on register to /passkey/set", () => {
    cy.visit("/register");
    cy.get('input[data-testid="firstname-text-input"]').focus().type("John");
    cy.get('input[data-testid="lastname-text-input"]').focus().type("Doe");
    cy.get('input[data-testid="email-text-input"]').focus().type("john@example.com");
    cy.get('input[type="checkbox"][value="privacypolicy"]').check();
    cy.get('input[type="checkbox"][value="tos"]').check();
    cy.get('button[type="submit"]').click();
    cy.url().should("include", Cypress.config().baseUrl + "/passkey/set");
  });
});
