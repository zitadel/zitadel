import { addStub, removeStub } from "../support/mock";

describe("/passkey/login", () => {
  it("should redirect a user with password authentication to /password", () => {
    removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
    addStub("zitadel.user.v2alpha.SessionService", "SetSession", {
      //   authMethodTypes: [2, 1],
      id: "221390781972217886",
      creationDate: new Date("2023-07-04T07:19:49.178Z"),
      changeDate: new Date("2023-07-04T07:19:54.617Z"),
      sequence: 854,
      factors: {
        user: {
          displayName: "John Doe",
          id: "221256020561756190",
          loginName: "johndoe@zitadel.com",
          verifiedAt: "2023-07-04T07:19:49.168Z",
          sessionId: "221390781972217886",
        },
      },
      metadata: {},
      domain: "localhost",
    });

    removeStub(
      "zitadel.user.v2alpha.SessionService",
      "ListAuthenticationMethodTypes"
    );
    addStub(
      "zitadel.user.v2alpha.SessionService",
      "ListAuthenticationMethodTypes",
      {
        authMethodTypes: [1], // 1 for password authentication
      }
    );

    cy.visit("/loginname?loginName=johndoe%40zitadel.com&submit=true");
    cy.location("pathname", { timeout: 10_000 }).should(
      "eq",
      "/password?loginName=johndoe%40zitadel.com"
    );
  });
  it("should redirect a user with passwordless authentication to /passkey/login", () => {
    removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
    addStub("zitadel.user.v2alpha.SessionService", "SetSession", {
      //   authMethodTypes: [2, 1],
      id: "221390781972217886",
      creationDate: new Date("2023-07-04T07:19:49.178Z"),
      changeDate: new Date("2023-07-04T07:19:54.617Z"),
      sequence: 854,
      factors: {
        user: {
          displayName: "John Doe",
          id: "221256020561756190",
          loginName: "johndoe@zitadel.com",
          verifiedAt: "2023-07-04T07:19:49.168Z",
          sessionId: "221390781972217886",
        },
      },
      metadata: {},
      domain: "localhost",
    });

    removeStub(
      "zitadel.user.v2alpha.SessionService",
      "ListAuthenticationMethodTypes"
    );
    addStub(
      "zitadel.user.v2alpha.SessionService",
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
