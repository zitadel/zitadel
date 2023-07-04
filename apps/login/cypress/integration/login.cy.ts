import { addStub, removeStub } from "../support/mock";

describe("/passkey/login", () => {
  it("should redirect a user with password authentication to /password", () => {
    removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
    addStub("zitadel.user.v2alpha.SessionService", "SetSession");
    cy.visit("/passkey/login?loginName=zitadel-admin%40zitadel.localhost");
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/accounts");
  });
  it("should redirect a user with passwordless authentication to /passkey/login", () => {
    removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
    addStub("zitadel.user.v2alpha.SessionService", "SetSession");
    cy.visit("/passkey/login?loginName=zitadel-admin%40zitadel.localhost");
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/accounts");
  });
  it("should prompt a user to setup passwordless authentication if passkey is allowed in the login settings", () => {
    removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
    addStub("zitadel.user.v2alpha.SessionService", "SetSession");
    cy.visit("/passkey/login?loginName=zitadel-admin%40zitadel.localhost");
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/accounts");
  });
  it("redirects after successful login", () => {
    removeStub("zitadel.user.v2alpha.SessionService", "SetSession");
    addStub("zitadel.user.v2alpha.SessionService", "SetSession");
    cy.visit("/passkey/login?loginName=zitadel-admin%40zitadel.localhost");
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/accounts");
  });
});
