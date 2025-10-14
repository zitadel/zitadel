import { stub } from "../support/e2e";

describe("verify invite", () => {
  beforeEach(() => {
    stub("zitadel.org.v2.OrganizationService", "ListOrganizations", {
      data: {
        details: {
          totalResult: 1,
        },
        result: [{ id: "256088834543534543" }],
      },
    });

    stub("zitadel.user.v2.UserService", "ListAuthenticationMethodTypes", {
      data: {
        authMethodTypes: [], // user with no auth methods was invited
      },
    });

    stub("zitadel.user.v2.UserService", "GetUserByID", {
      data: {
        user: {
          userId: "221394658884845598",
          state: 1,
          username: "john@example.com",
          loginNames: ["john@example.com"],
          preferredLoginName: "john@example.com",
          human: {
            userId: "221394658884845598",
            state: 1,
            username: "john@example.com",
            loginNames: ["john@example.com"],
            preferredLoginName: "john@example.com",
            profile: {
              givenName: "John",
              familyName: "Doe",
              avatarUrl: "https://example.com/avatar.jpg",
            },
            email: {
              email: "john@example.com",
              isVerified: false,
            },
          },
        },
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

    stub("zitadel.settings.v2.SettingsService", "GetLoginSettings", {
      data: {
        settings: {
          passkeysType: 1,
          allowUsernamePassword: true,
        },
      },
    });
  });

  it("shows authenticators after successful invite verification", () => {
    stub("zitadel.user.v2.UserService", "VerifyInviteCode");

    cy.visit("/verify?userId=221394658884845598&code=abc&invite=true");
    cy.url().should("include", Cypress.config().baseUrl + "/authenticator/set");
  });

  it("shows an error if invite code validation failed", () => {
    stub("zitadel.user.v2.UserService", "VerifyInviteCode", {
      code: 3,
      error: "error validating code",
    });

    // TODO: Avoid uncaught exception in application
    cy.once("uncaught:exception", () => false);
    cy.visit("/verify?userId=221394658884845598&code=abc&invite=true");
    cy.contains("Could not verify invite", { timeout: 10_000 });
  });
});
