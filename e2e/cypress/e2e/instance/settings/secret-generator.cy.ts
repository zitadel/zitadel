const secretGeneratorSettingsPath = `/instance?id=secrets`;

beforeEach(() => {
  cy.context().as('ctx');
});

describe('instance secret generators', () => {
  describe('secret generator settings', () => {
    it(`should show secret generator cards`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.contains('Initialization Mail');
      cy.contains('Email verification');
      cy.contains('Phone verification');
      cy.contains('Password Reset');
      cy.contains('Passwordless Initialization');
      cy.contains('App Secret');
      cy.contains('One Time Password (OTP) - SMS');
      cy.contains('One Time Password (OTP) - Email');
    });

    it(`Initialization Mail should contain default settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.get('input[id="length1"]').should('have.value', '6');
      cy.get('input[id="expiry1"]').should('have.value', '4320');
      cy.get('mat-slide-toggle#includeLowerLetters1 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters1 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits1 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols1 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`Email verification should contain default settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.get('input[id="length2"]').should('have.value', '6');
      cy.get('input[id="expiry2"]').should('have.value', '60');
      cy.get('mat-slide-toggle#includeLowerLetters2 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters2 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits2 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols2 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`Phone verification should contain default settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.get('input[id="length3"]').should('have.value', '6');
      cy.get('input[id="expiry3"]').should('have.value', '60');
      cy.get('mat-slide-toggle#includeLowerLetters3 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters3 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits3 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols3 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`Password Reset should contain default settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.get('input[id="length4"]').should('have.value', '6');
      cy.get('input[id="expiry4"]').should('have.value', '60');
      cy.get('mat-slide-toggle#includeLowerLetters4 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters4 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits4 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols4 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`Passwordless Initialization should contain default settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.get('input[id="length5"]').should('have.value', '12');
      cy.get('input[id="expiry5"]').should('have.value', '60');
      cy.get('mat-slide-toggle#includeLowerLetters5 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeUpperLetters5 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits5 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols5 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`App Secret should contain default settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.get('input[id="length6"]').should('have.value', '64');
      cy.get('input[id="expiry6"]').should('have.value', '0');
      cy.get('mat-slide-toggle#includeLowerLetters6 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeUpperLetters6 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits6 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols6 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`One Time Password (OTP) - SMS should contain default settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.get('input[id="length7"]').should('have.value', '8');
      cy.get('input[id="expiry7"]').should('have.value', '5');
      cy.get('mat-slide-toggle#includeLowerLetters7 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters7 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeDigits7 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols7 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`One Time Password (OTP) - Email should contain default settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.get('input[id="length8"]').should('have.value', '8');
      cy.get('input[id="expiry8"]').should('have.value', '5');
      cy.get('mat-slide-toggle#includeLowerLetters8 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters8 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeDigits8 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols8 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`Initialization Mail should update settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.wait(1000);
      cy.get('input[id="length1"]').clear().type('64');
      cy.get('mat-slide-toggle#includeLowerLetters1 button').click();
      cy.get('button[id="saveSecretGenerator1"]').click();
      cy.wait(1000);
      cy.get('input[id="length1"]').should('have.value', '64');
      cy.get('input[id="expiry1"]').should('have.value', '4320');
      cy.get('mat-slide-toggle#includeLowerLetters1 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeUpperLetters1 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits1 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols1 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`Email verification should update settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.wait(1000);
      cy.get('input[id="length2"]').clear().type('64');
      cy.get('mat-slide-toggle#includeUpperLetters2 button').click();
      cy.get('button[id="saveSecretGenerator2"]').click();
      cy.wait(1000);
      cy.get('input[id="length2"]').should('have.value', '64');
      cy.get('input[id="expiry2"]').should('have.value', '60');
      cy.get('mat-slide-toggle#includeLowerLetters2 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters2 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeDigits2 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols2 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`Phone verification should update settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.wait(1000);
      cy.get('input[id="expiry3"]').clear().type('10');
      cy.get('mat-slide-toggle#includeSymbols3 button').click();
      cy.get('button[id="saveSecretGenerator3"]').click();
      cy.wait(1000);
      cy.get('input[id="length3"]').should('have.value', '6');
      cy.get('input[id="expiry3"]').should('have.value', '10');
      cy.get('mat-slide-toggle#includeLowerLetters3 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters3 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits3 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols3 button').should('have.attr', 'aria-checked', 'true');
    });

    it(`Password Reset should update settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.wait(1000);
      cy.get('input[id="expiry4"]').clear().type('5');
      cy.get('mat-slide-toggle#includeDigits4 button').click();
      cy.get('button[id="saveSecretGenerator4"]').click();
      cy.wait(1000);
      cy.get('input[id="length4"]').should('have.value', '6');
      cy.get('input[id="expiry4"]').should('have.value', '5');
      cy.get('mat-slide-toggle#includeLowerLetters4 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters4 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits4 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeSymbols4 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`Passwordless Initialization should update settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.wait(1000);
      cy.get('input[id="length5"]').clear().type('64');
      cy.get('mat-slide-toggle#includeDigits5 button').click();
      cy.get('button[id="saveSecretGenerator5"]').click();
      cy.wait(1000);
      cy.get('input[id="length5"]').should('have.value', '64');
      cy.get('input[id="expiry5"]').should('have.value', '60');
      cy.get('mat-slide-toggle#includeLowerLetters5 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeUpperLetters5 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeDigits5 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeSymbols5 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`App Secret should update settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.wait(1000);
      cy.get('input[id="length6"]').clear().type('32');
      cy.get('input[id="expiry6"]').clear().type('120');
      cy.get('mat-slide-toggle#includeUpperLetters6 button').click();
      cy.get('button[id="saveSecretGenerator6"]').click();
      cy.wait(1000);
      cy.get('input[id="length6"]').should('have.value', '32');
      cy.get('input[id="expiry6"]').should('have.value', '120');
      cy.get('mat-slide-toggle#includeLowerLetters6 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeUpperLetters6 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeDigits6 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols6 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`One Time Password (OTP) - SMS should update settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.wait(1000);
      cy.get('input[id="expiry7"]').clear().type('120');
      cy.get('mat-slide-toggle#includeLowerLetters7 button').click();
      cy.get('button[id="saveSecretGenerator7"]').click();
      cy.wait(1000);
      cy.get('input[id="length7"]').should('have.value', '8');
      cy.get('input[id="expiry7"]').should('have.value', '120');
      cy.get('mat-slide-toggle#includeLowerLetters7 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeUpperLetters7 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeDigits7 button').should('have.attr', 'aria-checked', 'true');
      cy.get('mat-slide-toggle#includeSymbols7 button').should('have.attr', 'aria-checked', 'false');
    });

    it(`One Time Password (OTP) should update settings`, () => {
      cy.visit(secretGeneratorSettingsPath);
      cy.wait(1000);
      cy.get('input[id="length8"]').clear().type('12');
      cy.get('input[id="expiry8"]').clear().type('90');
      cy.get('mat-slide-toggle#includeDigits8 button').click();
      cy.get('mat-slide-toggle#includeSymbols8 button').click();
      cy.get('button[id="saveSecretGenerator8"]').click();
      cy.wait(1000);
      cy.get('input[id="length8"]').should('have.value', '12');
      cy.get('input[id="expiry8"]').should('have.value', '90');
      cy.get('mat-slide-toggle#includeLowerLetters8 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeUpperLetters8 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeDigits8 button').should('have.attr', 'aria-checked', 'false');
      cy.get('mat-slide-toggle#includeSymbols8 button').should('have.attr', 'aria-checked', 'true');
    });
  });
});
