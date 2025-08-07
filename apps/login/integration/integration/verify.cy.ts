import { stub } from "../support/e2e";

describe("verify email", () => {
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
        authMethodTypes: [1], // set one method such that we know that the user was not invited
      },
    });

    stub("zitadel.user.v2.UserService", "SendEmailCode");

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
              isVerified: false, // email is not verified yet
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
  });

  it("shows an error if email code validation failed", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail", {
      code: 3,
      error: "error validating code",
    });
    // TODO: Avoid uncaught exception in application
    cy.once("uncaught:exception", () => false);
    cy.visit("/verify?userId=221394658884845598&code=abc");
    cy.contains("Could not verify email");
  });
});
