import { addStub, removeStub } from "../support/mock";

describe("/passkey/login", () => {
  it("should redirect a user with password authentication to /password", () => {
    removeStub("zitadel.user.v2alpha.SessionService", "CreateSession");
    addStub("zitadel.user.v2alpha.SessionService", "CreateSession", {
      details: {
        sequence: 859,
        changeDate: new Date("2023-07-04T07:58:20.126Z"),
        resourceOwner: "220516472055706145",
      },
      sessionId: "221394658884845598",
      sessionToken:
        "SDMc7DlYXPgwRJ-Tb5NlLqynysHjEae3csWsKzoZWLplRji0AYY3HgAkrUEBqtLCvOayLJPMd0ax4Q",
      challenges: undefined,
    });

    removeStub("zitadel.user.v2alpha.SessionService", "GetSession");
    addStub("zitadel.user.v2alpha.SessionService", "GetSession", {
      session: {
        id: "221394658884845598",
        creationDate: new Date("2023-07-04T07:58:20.026Z"),
        changeDate: new Date("2023-07-04T07:58:20.126Z"),
        sequence: 859,
        factors: {
          user: {
            id: "123",
            loginName: "zitadel-admin@zitadel.localhost",
          },
          password: undefined,
          passkey: undefined,
          intent: undefined,
        },
        metadata: {},
        domain: "localhost",
      },
    });

    removeStub(
      "zitadel.user.v2alpha.UserService",
      "ListAuthenticationMethodTypes"
    );
    addStub(
      "zitadel.user.v2alpha.UserService",
      "ListAuthenticationMethodTypes",
      {
        authMethodTypes: [1], // 1 for password authentication
      }
    );

    cy.visit(
      "/loginname?loginName=zitadel-admin%40zitadel.localhost&submit=true"
    );
    cy.location("pathname", { timeout: 10_000 }).should(
      "eq",
      "/password?loginName=zitadel-admin%40zitadel.localhost&promptPasswordless=true"
    );
  });
  it("should redirect a user with passwordless authentication to /passkey/login", () => {
    removeStub("zitadel.user.v2alpha.SessionService", "CreateSession");
    addStub("zitadel.user.v2alpha.SessionService", "CreateSession", {
      details: {
        sequence: 859,
        changeDate: new Date("2023-07-04T07:58:20.126Z"),
        resourceOwner: "220516472055706145",
      },
      sessionId: "221394658884845598",
      sessionToken:
        "SDMc7DlYXPgwRJ-Tb5NlLqynysHjEae3csWsKzoZWLplRji0AYY3HgAkrUEBqtLCvOayLJPMd0ax4Q",
      challenges: undefined,
    });

    removeStub("zitadel.user.v2alpha.SessionService", "GetSession");
    addStub("zitadel.user.v2alpha.SessionService", "GetSession", {
      session: {
        id: "221394658884845598",
        creationDate: new Date("2023-07-04T07:58:20.026Z"),
        changeDate: new Date("2023-07-04T07:58:20.126Z"),
        sequence: 859,
        factors: {
          user: {
            id: "123",
            loginName: "johndoe@zitadel.com",
          },
          password: undefined,
          passkey: undefined,
          intent: undefined,
        },
        metadata: {},
        domain: "localhost",
      },
    });

    removeStub(
      "zitadel.user.v2alpha.UserService",
      "ListAuthenticationMethodTypes"
    );
    addStub(
      "zitadel.user.v2alpha.UserService",
      "ListAuthenticationMethodTypes",
      {
        authMethodTypes: [2], // 2 for passwordless authentication
      }
    );

    cy.visit("/loginname?loginName=johndoe%40zitadel.com&submit=true");
    cy.location("pathname", { timeout: 10_000 }).should(
      "eq",
      "/passkey/login?loginName=zitadel-admin%40zitadel.localhost"
    );
  });

  //   it("should prompt a user to setup passwordless authentication if passkey is allowed in the login settings", () => {
  //     removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
  //     addStub("zitadel.user.v2alpha.SessionService", "SetSession");
  //     cy.visit("/passkey/login?loginName=zitadel-admin%40zitadel.localhost");
  //     cy.location("pathname", { timeout: 10_000 }).should("eq", "/accounts");
  //   });
  //   it("redirects after successful login", () => {
  //     removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
  //     addStub("zitadel.user.v2alpha.SessionService", "SetSession");
  //     cy.visit("/passkey/login?loginName=zitadel-admin%40zitadel.localhost");
  //     cy.location("pathname", { timeout: 10_000 }).should("eq", "/accounts");
  //   });
});

// removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
// addStub("zitadel.user.v2alpha.SessionService", "SetSession", {
//   id: "221390781972217886",
//   creationDate: new Date("2023-07-04T07:19:49.178Z"),
//   changeDate: new Date("2023-07-04T07:19:54.617Z"),
//   sequence: 854,
//   factors: {
//     user: {
//       displayName: "John Doe",
//       id: "221256020561756190",
//       loginName: "johndoe@zitadel.com",
//       verifiedAt: "2023-07-04T07:19:49.168Z",
//       sessionId: "221390781972217886",
//     },
//   },
//   metadata: {},
//   domain: "localhost",
// });
