const notificationPath = `/settings?id=notifications`;
const smtpPath = `/settings?id=smtpprovider`;
const smsPath = `/settings?id=smsprovider`;

beforeEach(() => {
  cy.context().as('ctx');
});

describe('instance notifications', () => {
  describe('notification settings', () => {
    it(`should show notification settings`, () => {
      cy.visit(notificationPath);
      cy.contains('Notification');
    });
  });

  describe('smtp settings', () => {
    it(`should show SMTP provider settings`, () => {
      cy.visit(smtpPath);
      cy.contains('SMTP Settings');
    });
    it(`should add SMTP provider settings`, () => {
      cy.visit(smtpPath);
      cy.get('[formcontrolname="senderAddress"]').clear().type('sender@example.com');
      cy.get('[formcontrolname="senderName"]').clear().type('Zitadel');
      cy.get('[formcontrolname="hostAndPort"]').clear().type('smtp.mailtrap.io:2525');
      cy.get('[formcontrolname="user"]').clear().type('user@example.com');
      cy.get('[data-e2e="save-smtp-settings-button"]').click();
      cy.shouldConfirmSuccess();
      cy.get('[formcontrolname="senderAddress"]').should('have.value', 'sender@example.com');
      cy.get('[formcontrolname="senderName"]').should('have.value', 'Zitadel');
      cy.get('[formcontrolname="hostAndPort"]').should('have.value', 'smtp.mailtrap.io:2525');
      cy.get('[formcontrolname="user"]').should('have.value', 'user@example.com');
    });
    it(`should add SMTP provider password`, () => {
      cy.visit(smtpPath);
      cy.get('[data-e2e="add-smtp-password-button"]').click();
      cy.get('[data-e2e="notification-setting-password"]').clear().type('dummy@example.com');
      cy.get('[data-e2e="save-notification-setting-password-button"]').click();
      cy.shouldConfirmSuccess();
    });
  });

  describe('sms settings', () => {
    it(`should show SMS provider settings`, () => {
      cy.visit(smsPath);
      cy.contains('SMS Settings');
    });

    it(`should add SMS provider`, () => {
      cy.visit(smsPath);
      cy.get('[data-e2e="new-twilio-button"]').click();
      cy.get('[formcontrolname="sid"]').clear().type('test');
      cy.get('[formcontrolname="token"]').clear().type('token');
      cy.get('[formcontrolname="senderNumber"]').clear().type('2312123132');
      cy.get('[data-e2e="save-sms-settings-button"]').click();
      cy.shouldConfirmSuccess();
      cy.get('h4').contains('Twilio');
      cy.get('.state').contains('Inactive');
    });

    it(`should activate SMS provider`, () => {
      cy.visit(smsPath);
      cy.get('h4').contains('Twilio');
      cy.get('.state').contains('Inactive');
      cy.get('[data-e2e="activate-sms-provider-button"]').click();
      cy.shouldConfirmSuccess();
      cy.get('.state').contains('Active');
    });

    it(`should edit SMS provider`, () => {
      cy.visit(smsPath);
      cy.get('h4').contains('Twilio');
      cy.get('.state').contains('Active');
      cy.get('[data-e2e="new-twilio-button"]').click();
      cy.get('[formcontrolname="sid"]').should('have.value', 'test');
      cy.get('[formcontrolname="senderNumber"]').should('have.value', '2312123132');
      cy.get('[formcontrolname="sid"]').clear().type('test2');
      cy.get('[formcontrolname="senderNumber"]').clear().type('6666666666');
      cy.get('[data-e2e="save-sms-settings-button"]').click();
      cy.shouldConfirmSuccess();
    });

    it(`should contain edited values`, () => {
      cy.visit(smsPath);
      cy.get('h4').contains('Twilio');
      cy.get('.state').contains('Active');
      cy.get('[data-e2e="new-twilio-button"]').click();
      cy.get('[formcontrolname="sid"]').should('have.value', 'test2');
      cy.get('[formcontrolname="senderNumber"]').should('have.value', '6666666666');
    });

    it(`should edit SMS provider token`, () => {
      cy.visit(smsPath);
      cy.get('h4').contains('Twilio');
      cy.get('.state').contains('Active');
      cy.get('[data-e2e="new-twilio-button"]').click();
      cy.get('[data-e2e="edit-sms-token-button"]').click();
      cy.get('[data-e2e="notification-setting-password"]').clear().type('newsupertoken');
      cy.get('[data-e2e="save-notification-setting-password-button"]').click();
      cy.shouldConfirmSuccess();
    });

    it(`should remove SMS provider`, () => {
      cy.visit(smsPath);
      cy.get('h4').contains('Twilio');
      cy.get('.state').contains('Active');
      cy.get('[data-e2e="remove-sms-provider-button"]').click();
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
    });
  });
});
