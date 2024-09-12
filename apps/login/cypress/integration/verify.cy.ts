import { timestampFromDate } from "@zitadel/client";
import { stub } from "../support/mock";

describe("/verify", () => {
  it("if no MFA required, redirects to loginname after successful email verification", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail");
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
              loginName: "john@zitadel.com",
            },
            otpEmail: {
              // set a factor
              verifiedAt: timestampFromDate(
                new Date("2024-04-04T09:40:55.577Z"),
              ),
            },
            password: undefined,
            webAuthN: undefined,
            intent: undefined,
          },
          metadata: {},
        },
      },
    });
    cy.visit("/verify?userId=123&code=abc&submit=true");
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/loginname");
  });
  it("if MFA is required and no mfa factor is found, redirects to mfa/set after successful email verification", () => {
    stub("zitadel.settings.v2.SettingsService", "GetLoginSettings", {
      data: {
        settings: {
          forceMfa: true,
        },
      },
    });
    stub("zitadel.user.v2.UserService", "VerifyEmail");
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
              loginName: "john@zitadel.com",
            },
            otpEmail: undefined,
            password: undefined,
            webAuthN: undefined,
            intent: undefined,
          },
          metadata: {},
        },
      },
    });
    cy.visit("/verify?userId=123&code=abc&submit=true");
    cy.location("pathname", { timeout: 10_000 }).should("eq", "/mfa/set");
  });
  it("shows an error if validation failed", () => {
    stub("zitadel.user.v2.UserService", "VerifyEmail", {
      code: 3,
      error: "error validating code",
    });
    // TODO: Avoid uncaught exception in application
    cy.once("uncaught:exception", () => false);
    cy.visit("/verify?userId=123&code=abc&submit=true");
    cy.contains("Could not verify email");
  });
});
