import { stub } from "../support/mock";

describe("/verify", () => {
  it("shows authenticators after successful invite verification", () => {
    stub("zitadel.user.v2.UserService", "VerifyInviteCode");
    cy.visit("/verify?userId=123&code=abc&submit=true&invite=true");
    cy.location("pathname", { timeout: 10000 }).should(
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
    cy.visit("/verify?userId=123&code=abc&submit=true&invite=true");
    cy.contains("Could not verify invite", { timeout: 10000 });
  });

  it("shows password and passkey method after successful invite verification", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail");
    cy.visit("/verify?userId=123&code=abc&submit=true");
    cy.location("pathname", { timeout: 10000 }).should(
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
    cy.contains("Could not verify email", { timeout: 10000 });
  });
});
