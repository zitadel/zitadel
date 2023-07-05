import { addStub, removeStub } from "../support/mock";

describe("login", () => {
  beforeEach(() => {
    removeStub("zitadel.session.v2alpha.SessionService", "CreateSession");
    addStub("zitadel.session.v2alpha.SessionService", "CreateSession", {
      data: {
        details: {
          sequence: 859,
          changeDate: "2023-07-04T07:58:20.126Z",
          resourceOwner: "220516472055706145",
        },
        sessionId: "221394658884845598",
        sessionToken:
          "SDMc7DlYXPgwRJ-Tb5NlLqynysHjEae3csWsKzoZWLplRji0AYY3HgAkrUEBqtLCvOayLJPMd0ax4Q",
        challenges: undefined,
      },
    });

    removeStub("zitadel.session.v2alpha.SessionService", "GetSession");
    addStub("zitadel.session.v2alpha.SessionService", "GetSession", {
      data: {
        session: {
          id: "221394658884845598",
          creationDate: "2023-07-04T07:58:20.026Z",
          changeDate: "2023-07-04T07:58:20.126Z",
          sequence: 859,
          factors: {
            user: {
              id: "123",
              loginName: "john@zitadel.com",
            },
            password: undefined,
            passkey: undefined,
            intent: undefined,
          },
          metadata: {},
          domain: "localhost",
        },
      },
    });

    removeStub("zitadel.settings.v2alpha.SessionService", "GetLoginSettings");
    addStub("zitadel.settings.v2alpha.SettingsService", "GetLoginSettings", {
      data: {
        settings: {
          passkeysType: 1,
        },
      },
    });
  });
  describe("password login", () => {
    beforeEach(() => {
      removeStub(
        "zitadel.user.v2alpha.UserService",
        "ListAuthenticationMethodTypes"
      );
      addStub(
        "zitadel.user.v2alpha.UserService",
        "ListAuthenticationMethodTypes",
        {
          data: {
            authMethodTypes: [1], // 1 for password authentication
          },
        }
      );
    });
    it("should redirect a user with password authentication to /password", () => {
      cy.visit("/loginname?loginName=johndoe%40zitadel.com&submit=true");
      cy.location("pathname", { timeout: 10_000 }).should("eq", "/password");
    });
  });

  describe("passkey login", () => {
    beforeEach(() => {
      removeStub(
        "zitadel.user.v2alpha.UserService",
        "ListAuthenticationMethodTypes"
      );
      addStub(
        "zitadel.user.v2alpha.UserService",
        "ListAuthenticationMethodTypes",
        {
          data: {
            authMethodTypes: [2], // 2 for passwordless authentication
          },
        }
      );
    });
    it("should redirect a user with passwordless authentication to /passkey/login", () => {
      cy.visit("/loginname?loginName=johndoe%40zitadel.com&submit=true");
      cy.location("pathname", { timeout: 10_000 }).should(
        "eq",
        "/passkey/login"
      );
    });
  });

  describe("password login with passkey prompt", () => {
    beforeEach(() => {
      removeStub(
        "zitadel.user.v2alpha.UserService",
        "ListAuthenticationMethodTypes"
      );
      addStub(
        "zitadel.user.v2alpha.UserService",
        "ListAuthenticationMethodTypes",
        {
          data: {
            authMethodTypes: [1], // 1 for password authentication
          },
        }
      );

      removeStub("zitadel.session.v2alpha.SessionService", "SetSession");
      addStub("zitadel.session.v2alpha.SessionService", "SetSession", {
        data: {
          details: {
            sequence: 859,
            changeDate: "2023-07-04T07:58:20.126Z",
            resourceOwner: "220516472055706145",
          },
          sessionToken:
            "SDMc7DlYXPgwRJ-Tb5NlLqynysHjEae3csWsKzoZWLplRji0AYY3HgAkrUEBqtLCvOayLJPMd0ax4Q",
          challenges: undefined,
        },
      });
    });
    it("should prompt a user to setup passwordless authentication if passkey is allowed in the login settings", () => {
      cy.visit("/loginname?loginName=john%40zitadel.com&submit=true");
      cy.location("pathname", { timeout: 10_000 }).should("eq", "/password");
      cy.get('input[type="password"]').focus().type("MyStrongPassword!1");
      cy.get('button[type="submit"]').click();
      cy.location("pathname", { timeout: 10_000 }).should("eq", "/passkey/add");
    });
  });
});
