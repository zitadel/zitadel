import { v4 as uuidv4 } from 'uuid';

import { Context } from '../../../support/commands';
import { ensureOrgExists } from '../../../support/api/orgs';
import { activateSMTPProvider, ensureSMTPProviderExists } from '../../../support/api/smtp';
import { ensureSMSProviderDoesntExist, ensureSMSProviderExists } from '../../../support/api/sms';

const notificationPath = `/instance?id=notifications`;
const smtpPath = `/instance?id=smtpprovider`;
const smsPath = `/instance?id=smsprovider`;

beforeEach(() => {
  cy.context().as('ctx');
});

type SMTPProvider = {
  description: string;
  rowSelector: string;
};

describe('instance notifications', () => {
  describe('notification settings', () => {
    it(`should show notification settings`, () => {
      cy.visit(notificationPath);
      cy.contains('Notification');
    });
  });

  describe('smtp settings', () => {
    beforeEach(() => {
      const description = `mailgun-${uuidv4()}`;
      cy.wrap<SMTPProvider>({ description, rowSelector: `tr:contains('${description}')` }).as('provider');
    });

    it(`should add Mailgun SMTP provider settings`, () => {
      cy.get<SMTPProvider>('@provider').then((provider) => {
        cy.visit(smtpPath);
        cy.get(`a:contains('Mailgun')`).click();
        cy.get('[formcontrolname="description"]').should('be.enabled').clear().type(provider.description);
        cy.get('[formcontrolname="hostAndPort"]').should('have.value', 'smtp.mailgun.org:587');
        cy.get('[formcontrolname="user"]').should('be.enabled').clear().type('user@example.com');
        cy.get('[formcontrolname="password"]').should('be.enabled').clear().type('password');
        cy.get('[data-e2e="continue-to-2nd-form"]').should('be.enabled').click();
        cy.get('[formcontrolname="senderAddress"]').should('be.enabled').clear().type('sender1@example.com');
        cy.get('[formcontrolname="senderName"]').should('be.enabled').clear().type('Test1');
        cy.get('[formcontrolname="replyToAddress"]').should('be.enabled').clear().type('replyto1@example.com');
        cy.get('[data-e2e="continue-button"]').should('be.enabled').click();
        cy.get('[data-e2e="create-button"]').should('be.enabled').click();
        cy.shouldConfirmSuccess();
        cy.get('[data-e2e="close-button"]').should('be.enabled').click();
        cy.get(provider.rowSelector).contains('smtp.mailgun.org:587');
        cy.get(provider.rowSelector).contains('sender1@example.com');
      });
    });

    it(`should add Mailgun SMTP provider settings and activate it using wizard`, () => {
      cy.get<SMTPProvider>('@provider').then((provider) => {
        cy.visit(smtpPath);
        cy.get(`a:contains('Mailgun')`).click();
        cy.get('[formcontrolname="description"]').should('be.enabled').clear().type(provider.description);
        cy.get('[formcontrolname="hostAndPort"]').should('have.value', 'smtp.mailgun.org:587');
        cy.get('[formcontrolname="user"]').should('be.enabled').clear().type('user@example.com');
        cy.get('[formcontrolname="password"]').should('be.enabled').clear().type('password');
        cy.get('[data-e2e="continue-to-2nd-form"]').should('be.enabled').click();
        cy.get('[formcontrolname="senderAddress"]').should('be.enabled').clear().type('sender1@example.com');
        cy.get('[formcontrolname="senderName"]').should('be.enabled').clear().type('Test1');
        cy.get('[formcontrolname="replyToAddress"]').should('be.enabled').clear().type('replyto1@example.com');
        cy.get('[data-e2e="continue-button"]').should('be.enabled').click();
        cy.get('[data-e2e="create-button"]').click();
        cy.shouldConfirmSuccess();
        cy.get('[data-e2e="activate-button"]').click();
        cy.shouldConfirmSuccess();
        cy.get('[data-e2e="close-button"]').click();
        cy.get(provider.rowSelector).find('[data-e2e="active-provider"]');
        cy.get(provider.rowSelector).contains('smtp.mailgun.org:587');
        cy.get(provider.rowSelector).contains('sender1@example.com');
      });
    });

    describe('with inactive existing', () => {
      beforeEach(() => {
        cy.get<Context>('@ctx').then((ctx) => {
          cy.get<SMTPProvider>('@provider').then(({ description }) => {
            ensureSMTPProviderExists(ctx.api, description);
          });
        });
        cy.visit(smtpPath);
      });

      it(`should change Mailgun SMTP provider settings`, () => {
        cy.get<SMTPProvider>('@provider').then(({ rowSelector }) => {
          cy.get(rowSelector).click();
          cy.get('[data-e2e="continue-to-2nd-form"]').click();
          cy.get('[formcontrolname="senderAddress"]').should('be.enabled').clear().type('senderchange1@example.com');
          cy.get('[formcontrolname="senderName"]').clear().type('Change1');
          cy.get('[data-e2e="continue-button"]').click();
          cy.get('[data-e2e="create-button"]').click();
          cy.shouldConfirmSuccess();
          cy.get('[data-e2e="close-button"]').click();
          cy.get(rowSelector).contains('senderchange1@example.com');
        });
      });
      it(`should activate Mailgun SMTP provider settings`, () => {
        cy.get<SMTPProvider>('@provider').then(({ rowSelector }) => {
          cy.get(rowSelector).find('[data-e2e="activate-provider-button"]').click({ force: true });
          cy.get('[data-e2e="confirm-dialog-button"]').click();
          cy.shouldConfirmSuccess();
          cy.get(rowSelector).find('[data-e2e="active-provider"]');
        });
      });

      it(`should delete Mailgun SMTP provider`, () => {
        cy.get<SMTPProvider>('@provider').then(({ rowSelector }) => {
          cy.get(rowSelector).find('[data-e2e="delete-provider-button"]').click({ force: true });
          cy.get('[data-e2e="confirm-dialog-input"]').focus().should('be.enabled').type('A Sender');
          cy.get('[data-e2e="confirm-dialog-button"]').click();
          cy.shouldConfirmSuccess();
          cy.get(rowSelector).should('not.exist');
        });
      });
    });
    describe('with active existing', () => {
      beforeEach(() => {
        cy.get<Context>('@ctx').then((ctx) => {
          cy.get<SMTPProvider>('@provider').then(({ description }) => {
            ensureSMTPProviderExists(ctx.api, description).then((providerId) => {
              activateSMTPProvider(ctx.api, providerId);
            });
          });
        });
        cy.pause();
        cy.visit(smtpPath);
      });

      it(`should deactivate an existing Mailgun SMTP provider using wizard`, () => {
        cy.get<SMTPProvider>('@provider').then(({ rowSelector }) => {
          cy.get(rowSelector).click();
          cy.get('[data-e2e="continue-to-2nd-form"]').click();
          cy.get('[data-e2e="continue-button"]').click();
          cy.get('[data-e2e="create-button"]').click();
          cy.shouldConfirmSuccess();
          cy.get('[data-e2e="deactivate-button"]').click();
          cy.shouldConfirmSuccess();
          cy.get('[data-e2e="close-button"]').click();
          cy.get(rowSelector).find('[data-e2e="active-provider"]').should('not.exist');
        });
      });
    });
  });

  describe('sms settings', () => {
    beforeEach(() => {
      cy.wrap<string>(`twilio-${uuidv4()}`).as('uniqueSid');
    });

    describe('without existing', () => {
      beforeEach(() => {
        cy.get<Context>('@ctx').then((ctx) => {
          ensureSMSProviderDoesntExist(ctx.api);
        });
      });

      it(`should add SMS provider`, () => {
        cy.visit(smsPath);
        cy.get('[data-e2e="new-twilio-button"]').click();
        cy.get('[formcontrolname="sid"]').should('be.enabled').clear().type('test');
        cy.get('[formcontrolname="token"]').should('be.enabled').clear().type('token');
        cy.get('[formcontrolname="senderNumber"]').should('be.enabled').clear().type('2312123132');
        cy.get('[data-e2e="save-sms-settings-button"]').click();
        cy.shouldConfirmSuccess();
        cy.get('h4').contains('Twilio');
        cy.get('.state').contains('Inactive');
      });
    });

    describe('with inactive existing', () => {
      beforeEach(() => {
        cy.get<Context>('@ctx').then((ctx) => {
          ensureSMSProviderExists(ctx.api);
          cy.visit(smsPath);
        });
      });

      it(`should activate SMS provider`, () => {
        cy.get('h4').contains('Twilio');
        cy.get('.state').contains('Inactive');
        cy.get('[data-e2e="activate-sms-provider-button"]').click();
        cy.shouldConfirmSuccess();
        cy.get('.state').contains('Active');
      });

      it(`should edit SMS provider`, () => {
        cy.get('h4').contains('Twilio');
        cy.get('[data-e2e="new-twilio-button"]').click();
        cy.get('[formcontrolname="sid"]').should('be.enabled').clear().type('test2');
        cy.get('[formcontrolname="senderNumber"]').should('be.enabled').clear().type('6666666666');
        cy.get('[data-e2e="save-sms-settings-button"]').click();
        cy.shouldConfirmSuccess();
        cy.get('[data-e2e="new-twilio-button"]').click();
        cy.get('[formcontrolname="sid"]').should('have.value', 'test2');
        cy.get('[formcontrolname="senderNumber"]').should('have.value', '6666666666');
      });

      it(`should edit SMS provider token`, () => {
        cy.get('h4').contains('Twilio');
        cy.get('[data-e2e="new-twilio-button"]').click();
        cy.get('[data-e2e="edit-sms-token-button"]').click();
        cy.get('[data-e2e="notification-setting-password"]').should('be.enabled').clear().type('newsupertoken');
        cy.get('[data-e2e="save-notification-setting-password-button"]').click();
        cy.shouldConfirmSuccess();
      });

      it(`should remove SMS provider`, () => {
        cy.get('h4').contains('Twilio');
        cy.get('[data-e2e="remove-sms-provider-button"]').click();
        cy.get('[data-e2e="confirm-dialog-button"]').click();
        cy.shouldConfirmSuccess();
      });
    });
  });
});
