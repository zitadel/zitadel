const notificationPath = `/instance?id=notifications`;
const smtpPath = `/instance?id=smtpprovider`;
const smsPath = `/instance?id=smsprovider`;

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
      cy.contains('SMTP Provider');
    });
    it(`should add Mailgun SMTP provider settings`, () => {
      let rowSelector = `a:contains('Mailgun')`;
      cy.visit(smtpPath);
      cy.get(rowSelector).click();
      cy.get('[formcontrolname="hostAndPort"]').should('have.value', 'smtp.mailgun.org:587');
      cy.get('[formcontrolname="user"]').clear().type('user@example.com');
      cy.get('[formcontrolname="password"]').clear().type('password');
      cy.get('[data-e2e="continue-button"]').click();
      cy.get('[formcontrolname="senderAddress"]').clear().type('sender1@example.com');
      cy.get('[formcontrolname="senderName"]').clear().type('Test1');
      cy.get('[formcontrolname="replyToAddress"]').clear().type('replyto1@example.com');
      cy.get('[data-e2e="create-button"]').click();
      cy.shouldConfirmSuccess();
      cy.get('tr').contains('mailgun');
      cy.get('tr').contains('smtp.mailgun.org:587');
      cy.get('tr').contains('sender1@example.com');
    });
    it(`should change Mailgun SMTP provider settings`, () => {
      let rowSelector = `tr:contains('mailgun')`;
      cy.visit(smtpPath);
      cy.get(rowSelector).click();
      cy.get('[formcontrolname="hostAndPort"]').should('have.value', 'smtp.mailgun.org:587');
      cy.get('[formcontrolname="user"]').should('have.value', 'user@example.com');
      cy.get('[formcontrolname="user"]').clear().type('change@example.com');
      cy.get('[data-e2e="continue-button"]').click();
      cy.get('[formcontrolname="senderAddress"]').should('have.value', 'sender1@example.com');
      cy.get('[formcontrolname="senderName"]').should('have.value', 'Test1');
      cy.get('[formcontrolname="replyToAddress"]').should('have.value', 'replyto1@example.com');
      cy.get('[formcontrolname="senderAddress"]').clear().type('senderchange1@example.com');
      cy.get('[formcontrolname="senderName"]').clear().type('Change1');
      cy.get('[data-e2e="create-button"]').click();
      cy.shouldConfirmSuccess();
      rowSelector = `tr:contains('mailgun')`;
      cy.get(rowSelector).contains('mailgun');
      cy.get(rowSelector).contains('smtp.mailgun.org:587');
      cy.get(rowSelector).contains('senderchange1@example.com');
    });
    it(`should activate Mailgun SMTP provider settings`, () => {
      let rowSelector = `tr:contains('smtp.mailgun.org:587')`;
      cy.visit(smtpPath);
      cy.get(rowSelector).find('[data-e2e="activate-provider-button"]').click({ force: true });
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
      rowSelector = `tr:contains('smtp.mailgun.org:587')`;
      cy.get(rowSelector).find('[data-e2e="active-provider"]');
      cy.get(rowSelector).contains('mailgun');
      cy.get(rowSelector).contains('smtp.mailgun.org:587');
      cy.get(rowSelector).contains('senderchange1@example.com');
    });
    it(`should add Mailjet SMTP provider settings`, () => {
      let rowSelector = `a:contains('Mailjet')`;
      cy.visit(smtpPath);
      cy.get(rowSelector).click();
      cy.get('[formcontrolname="hostAndPort"]').should('have.value', 'in-v3.mailjet.com:587');
      cy.get('[formcontrolname="user"]').clear().type('user@example.com');
      cy.get('[formcontrolname="password"]').clear().type('password');
      cy.get('[data-e2e="continue-button"]').click();
      cy.get('[formcontrolname="senderAddress"]').clear().type('sender2@example.com');
      cy.get('[formcontrolname="senderName"]').clear().type('Test2');
      cy.get('[formcontrolname="replyToAddress"]').clear().type('replyto2@example.com');
      cy.get('[data-e2e="create-button"]').click();
      cy.shouldConfirmSuccess();
      rowSelector = `tr:contains('mailjet')`;
      cy.get(rowSelector).contains('mailjet');
      cy.get(rowSelector).contains('in-v3.mailjet.com:587');
      cy.get(rowSelector).contains('sender2@example.com');
    });
    it(`should activate Mailjet SMTP provider settings an disable Mailgun`, () => {
      let rowSelector = `tr:contains('in-v3.mailjet.com:587')`;
      cy.visit(smtpPath);
      cy.get(rowSelector).find('[data-e2e="activate-provider-button"]').click({ force: true });
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
      cy.get(rowSelector).find('[data-e2e="active-provider"]');
      cy.get(rowSelector).contains('mailjet');
      cy.get(rowSelector).contains('in-v3.mailjet.com:587');
      cy.get(rowSelector).contains('sender2@example.com');
      rowSelector = `tr:contains('mailgun')`;
      cy.get(rowSelector).find('[data-e2e="active-provider"]').should('not.exist');
    });
    it(`should deactivate Mailjet SMTP provider`, () => {
      let rowSelector = `tr:contains('mailjet')`;
      cy.visit(smtpPath);
      cy.get(rowSelector).find('[data-e2e="deactivate-provider-button"]').click({ force: true });
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
      rowSelector = `tr:contains('mailjet')`;
      cy.get(rowSelector).find('[data-e2e="active-provider"]').should('not.exist');
      rowSelector = `tr:contains('mailgun')`;
      cy.get(rowSelector).find('[data-e2e="active-provider"]').should('not.exist');
    });
    it(`should delete Mailjet SMTP provider`, () => {
      let rowSelector = `tr:contains('mailjet')`;
      cy.visit(smtpPath);
      cy.get(rowSelector).find('[data-e2e="delete-provider-button"]').click({ force: true });
      cy.get('[data-e2e="confirm-dialog-input"]').focus().type('Test2');
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
      rowSelector = `tr:contains('mailjet')`;
      cy.get(rowSelector).should('not.exist');
    });
    it(`should delete Mailgun SMTP provider`, () => {
      let rowSelector = `tr:contains('mailgun')`;
      cy.visit(smtpPath);
      cy.get(rowSelector).find('[data-e2e="delete-provider-button"]').click({ force: true });
      cy.get('[data-e2e="confirm-dialog-input"]').focus().type('Change1');
      cy.get('[data-e2e="confirm-dialog-button"]').click();
      cy.shouldConfirmSuccess();
      rowSelector = `tr:contains('mailgun')`;
      cy.get(rowSelector).should('not.exist');
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
