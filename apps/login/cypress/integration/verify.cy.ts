import { stub } from "../support/mock";

describe("/verify", () => {
  it("redirects after successful email verification", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail");
    cy.visit("/verify?userId=123&code=abc&submit=true");
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/loginname");
  });
  it("shows an error if validation failed", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail", {
      code: 3,
      error: "error validating code",
    });
    // TODO: Avoid uncaught exception in application
    cy.once("uncaught:exception", () => false);
    cy.visit("/verify?userId=123&code=abc&submit=true");
    cy.contains("error validating code");
  });
});
