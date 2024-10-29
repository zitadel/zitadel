import { stub } from "../support/mock";

describe("verify invite", () => {
  //   beforeEach(() => {
  //     stub("zitadel.org.v2.OrganizationService", "ListOrganizations", {
  //       data: {
  //         details: {
  //           totalResult: 1,
  //         },
  //         result: [{ id: "256088834543534543" }],
  //       },
  //     });

  //     stub("zitadel.user.v2.UserService", "ListAuthenticationMethodTypes", {
  //       data: {
  //         authMethodTypes: [],
  //       },
  //     });

  //     stub("zitadel.user.v2.UserService", "GetUserById", {
  //       data: {
  //         user: {
  //           userId: "221394658884845598",
  //           state: 1,
  //           username: "john@zitadel.com",
  //           loginNames: ["john@zitadel.com"],
  //           preferredLoginName: "john@zitadel.com",
  //           human: {
  //             userId: "221394658884845598",
  //             state: 1,
  //             username: "john@zitadel.com",
  //             loginNames: ["john@zitadel.com"],
  //             preferredLoginName: "john@zitadel.com",
  //             profile: {
  //               givenName: "John",
  //               familyName: "Doe",
  //               avatarUrl: "https://zitadel.com/avatar.jpg",
  //             },
  //             email: {
  //               email: "john@zitadel.com",
  //               isVerified: true,
  //             },
  //           },
  //         },
  //       },
  //     });
  //   });

  it.only("shows authenticators after successful invite verification", () => {
    stub("zitadel.user.v2.UserService", "VerifyInviteCode");
    cy.visit("/verify?userId=123&code=abc&invite=true");
    cy.location("pathname", { timeout: 10_000 }).should(
      "eq",
      "/authenticator/set",
    );
  });

  it("shows an error if invite code validation failed", () => {
    stub("zitadel.user.v2.UserService", "VerifyInviteCode", {
      code: 3,
      error: "error validating code",
    });
    // TODO: Avoid uncaught exception in application
    cy.once("uncaught:exception", () => false);
    cy.visit("/verify?userId=123&code=abc&invite=true");
    cy.contains("Could not verify invite", { timeout: 10_000 });
  });
});

describe("verify email", () => {
  //   beforeEach(() => {
  //     stub("zitadel.org.v2.OrganizationService", "ListOrganizations", {
  //       data: {
  //         details: {
  //           totalResult: 1,
  //         },
  //         result: [{ id: "256088834543534543" }],
  //       },
  //     });

  //     stub("zitadel.user.v2.UserService", "ListAuthenticationMethodTypes", {
  //       data: {
  //         authMethodTypes: [],
  //       },
  //     });
  //   });

  it("shows password and passkey method after successful invite verification", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail");
    cy.visit("/verify?userId=123&code=abc");
    cy.location("pathname", { timeout: 10_000 }).should(
      "eq",
      "/authenticator/set",
    );
  });

  it("shows an error if invite code validation failed", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail", {
      code: 3,
      error: "error validating code",
    });
    // TODO: Avoid uncaught exception in application
    cy.once("uncaught:exception", () => false);
    cy.visit("/verify?userId=123&code=abc&submit=true");
    cy.contains("Could not verify email", { timeout: 10_000 });
  });
});
