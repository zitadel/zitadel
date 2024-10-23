import { stub } from "../support/mock";

describe("/verify", () => {
  it("shows password and passkey method after successful invite verification", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail");
    cy.visit("/verify?userId=123&code=abc&submit=true&invite=true");
    cy.contains("Password");
    cy.contains("Passkey");
  });
  it("shows an error if validation failed", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail", {
      code: 3,
      error: "error validating code",
    });
    // TODO: Avoid uncaught exception in application
    cy.once("uncaught:exception", () => false);
    cy.visit("/verify?userId=123&code=abc&submit=true");
    cy.contains("Could not verify user");
  });
});
