import { stub } from "../support/e2e";

describe("login", () => {
  beforeEach(() => {
    stub("zitadel.org.v2.OrganizationService", "ListOrganizations", {
      data: {
        details: {
          totalResult: 1,
        },
        result: [{ id: "256088834543534543" }],
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
  describe("password login", () => {
    beforeEach(() => {
      stub("zitadel.user.v2.UserService", "ListUsers", {
        data: {
          details: {
            totalResult: 1,
          },
          result: [
            {
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
                  isVerified: true,
                },
              },
            },
          ],
        },
      });
      stub("zitadel.user.v2.UserService", "ListAuthenticationMethodTypes", {
        data: {
          authMethodTypes: [1], // 1 for password authentication
        },
      });
    });
    it("should redirect a user with password authentication to /password", () => {
      cy.visit("/loginname?loginName=john%40example.com&submit=true");
      cy.url({ timeout: 5 * 60_000 }).should("include", Cypress.config().baseUrl + "/password");
    });
    describe("with passkey prompt", () => {
      beforeEach(() => {
        stub("zitadel.session.v2.SessionService", "SetSession", {
          data: {
            details: {
              sequence: 859,
              changeDate: "2023-07-04T07:58:20.126Z",
              resourceOwner: "220516472055706145",
            },
            sessionToken: "SDMc7DlYXPgwRJ-Tb5NlLqynysHjEae3csWsKzoZWLplRji0AYY3HgAkrUEBqtLCvOayLJPMd0ax4Q",
            challenges: undefined,
          },
        });
      });
      // it("should prompt a user to setup passwordless authentication if passkey is allowed in the login settings", () => {
      //   cy.visit("/loginname?loginName=john%40example.com&submit=true");
      //   cy.location("pathname", { timeout: 5 * 60_000 }).should("eq", "/password");
      //   cy.get('input[type="password"]').focus().type("MyStrongPassword!1");
      //   cy.get('button[type="submit"]').click();
      //   cy.location("pathname", { timeout: 10_000 }).should(
      //     "eq",
      //     "/passkey/set",
      //   );
      // });
    });
  });
  describe("passkey login", () => {
    beforeEach(() => {
      stub("zitadel.user.v2.UserService", "ListUsers", {
        data: {
          details: {
            totalResult: 1,
          },
          result: [
            {
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
                  isVerified: true,
                },
              },
            },
          ],
        },
      });
      stub("zitadel.user.v2.UserService", "ListAuthenticationMethodTypes", {
        data: {
          authMethodTypes: [2], // 2 for passwordless authentication
        },
      });
    });

    it("should redirect a user with passwordless authentication to /passkey", () => {
      cy.visit("/loginname?loginName=john%40example.com&submit=true");
      cy.url().should("include", Cypress.config().baseUrl + "/passkey");
    });
  });
});
