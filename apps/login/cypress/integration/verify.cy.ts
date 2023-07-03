import { addStub, removeStub } from "../support/mock";

describe("/verify", () => {
  it("redirects after successful email verification", () => {
    removeStub("zitadel.user.v2alpha.UserService", "VerifyEmail");
    addStub("zitadel.user.v2alpha.UserService", "VerifyEmail");
    cy.visit("/verify?userID=123&code=abc&submit=true");
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/loginname");
  });
  it("shows an error if validation failed", () => {
    removeStub("zitadel.user.v2alpha.UserService", "VerifyEmail");
    addStub("zitadel.user.v2alpha.UserService", "VerifyEmail", {
      code: 3,
      error: "error validating code",
    });
    // TODO: Avoid uncaught exception in application
    cy.once("uncaught:exception", () => false);
    cy.visit("/verify?userID=123&code=abc&submit=true");
    cy.contains("error validating code");
  });
});
