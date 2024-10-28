import { stub } from "../support/mock";

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
      authMethodTypes: [],
    },
  });

  describe("verify invite", () => {
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
});
