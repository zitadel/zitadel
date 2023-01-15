import { ensureActionExists, setTriggerTypes } from 'support/api/actions';
import { ensureDomainPolicy } from 'support/api/policies';
import { ensureProjectExists } from 'support/api/projects';
import { ensureRoleExists } from 'support/api/roles';
import { newTarget } from 'support/api/target';
import { ensureHumanDoesntExist, ensureHumanExists } from 'support/api/users';
import { ZITADELTarget } from 'support/commands';
import { login } from 'support/login/login';
import { register } from 'support/login/register';
import { sessionAsPredefinedUser, User } from 'support/login/session';

describe('actions', () => {
  const emailVerifiedScript = 'e2eSetEmailVerified',
    addGrantScript = 'e2eAddGrant',
    addMetadataScript = 'e2eAddMetadata',
    storeUsernameScript = 'e2eSetLastUsernameMD',
    projectName = 'e2eaction',
    roleKey = 'e2eactionrole',
    preCreationEmail = 'precre@action.com',
    postCreationEmail = 'postcre@action.com',
    postAuthPWEmail = 'postauthpw@action.com',
    postAuthOTPEmail = 'postauthotp@action.com',
    postAuthU2FEmail = 'postauthu2f@action.com',
    postAuthPWLessEmail = 'postauthpwless@action.com',
    userPw = 'Password1!',
    specPath = 'cypress/e2e/actions';

  beforeEach(() => {
    newTarget('e2eactions')
      .as('target')
      .then((target) => {
        ensureDomainPolicy(target, false, false, false);
        cy.readFile(`${specPath}/${emailVerifiedScript}.js`).then((script) => {
          ensureActionExists(target, emailVerifiedScript, script)
            .as('emailVerifiedId')
            .then((emailVerifiedId) => {
              setTriggerTypes(target, 3, 2, [emailVerifiedId]);
            });
        });
      });
  });
  describe('triggers', () => {
    describe('creation', () => {
      describe('pre', () => {
        beforeEach(() => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            ensureHumanDoesntExist(target, preCreationEmail);
            cy.get<number>('@emailVerifiedId').then((emailVerifiedId) => {
              cy.readFile(`${specPath}/${addMetadataScript}.js`).then((script) => {
                ensureActionExists(target, addMetadataScript, script).then((metadataId) => {
                  setTriggerTypes(target, 3, 2, [emailVerifiedId, metadataId]);
                });
              });
            });
          });
        });
        it(`shouldn't prompt for email code and add metadata`, () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            register(preCreationEmail, target.headers['x-zitadel-orgid']).then((userId) => {
              sessionAsPredefinedUser(User.IAMAdminUser);
              cy.visit(`/users/${userId}?id=metadata&org=${target.headers['x-zitadel-orgid']}`);
              cy.contains('tr', 'akey').contains('avalue');
            });
          });
        });
      });
    });
    describe('post', () => {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          cy.get<number>('@emailVerifiedId').then((emailVerifiedId) => {
            setTriggerTypes(target, 3, 2, [emailVerifiedId]);
          });
          ensureHumanDoesntExist(target, postCreationEmail);
          ensureProjectExists(target, projectName).then((projectId) => {
            ensureRoleExists(target, projectId, roleKey);
            cy.readFile<string>(`${specPath}/${addGrantScript}.js`).then((script) => {
              ensureActionExists(
                target,
                addGrantScript,
                script.replace('<PROJECT_ID>', `${projectId}`).replace('<ROLE_KEY>', roleKey),
              ).then((addGrantId) => {
                setTriggerTypes(target, 3, 3, [addGrantId]);
              });
            });
          });
        });
      });
      it(`should add a grant when registering via UI`, () => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          register(postCreationEmail, target.headers['x-zitadel-orgid']).then((userId) => {
            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.visit(`/users/${userId}?id=grants&org=${target.headers['x-zitadel-orgid']}`);
            cy.contains('tr', roleKey);
          });
        });
      });
    });
  });
  describe('post authentication', () => {
    beforeEach(() => {
      cy.get<ZITADELTarget>('@target').then((target) => {
        cy.readFile<string>(`${specPath}/${storeUsernameScript}.js`).then((script) => {
          ensureActionExists(target, storeUsernameScript, script).then((storeUsernameId) => {
            setTriggerTypes(target, 3, 1, [storeUsernameId]);
          });
        });
      });
    });
    describe('pw', () => {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureHumanDoesntExist(target, postAuthPWEmail);
          ensureHumanExists(target, postAuthPWEmail).as('userId');
        });
      });

      it('should store the username in metadata after password authentication', () => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          login(postAuthPWEmail, userPw, target.headers['x-zitadel-orgid']);
          cy.get('@userId').then((userId) => {
            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.visit(`/users/${userId}?id=metadata&org=${target.headers['x-zitadel-orgid']}`);
            cy.contains('tr', 'last username used').contains(postAuthPWEmail);
          });
        });
      });
    });
    describe('otp', () => {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureHumanDoesntExist(target, postAuthOTPEmail);
          ensureHumanExists(target, postAuthOTPEmail).as('userId');
        });
      });

      it('should store the username in metadata after otp authentication', () => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          login(postAuthOTPEmail, userPw, target.headers['x-zitadel-orgid']);
          cy.visit('/users/me?id=security');
          cy.get('[data-e2e="add-factor"]').should('be.visible').click();
          cy.get('[data-e2e="add-factor-otp"]').should('be.visible').click();
          cy.get('[data-e2e="otp-secret"]')
            .as('otpSecret')
            .then((secret) => {
              cy.task<string>('generateOTP', secret.text().trim()).then((token) => {
                cy.get('[data-e2e="otp-code-input"]').should('be.visible').type(token, { force: true });
                cy.get('[data-e2e="save-factor"]').should('be.visible').click();
              });
            });
          login(postAuthOTPEmail, userPw, target.headers['x-zitadel-orgid'], () => {
            cy.task<string>('generateOTP').then((token) => {
              cy.get('#code').should('be.visible').type(token);
              cy.get('#submit-button').should('be.visible').click();
            });
          });

          sessionAsPredefinedUser(User.IAMAdminUser);
          cy.get('@userId').then((userId) => {
            cy.visit(`/users/${userId}?org=${target.headers['x-zitadel-orgid']}&id=metadata`);
            cy.contains('tr', 'last username used').contains(postAuthOTPEmail);
          });
        });
      });
    });
    describe('u2f', () => {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureHumanDoesntExist(target, postAuthU2FEmail);
          ensureHumanExists(target, postAuthU2FEmail).as('userId');
        });
      });
      it(
        'should store the username in metadata after u2f authentication',
        {
          // Verifying a key that was registered in an origin other than the RP's origin fails.
          // It is tagged here so that it can be grepped and skipped when run against a dev server.
          tags: ['@same-origin'],
        },
        () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            login(postAuthU2FEmail, userPw, target.headers['x-zitadel-orgid']);
            cy.visit('/users/me?id=security');
            cy.get('[data-e2e="add-factor"]').should('be.visible').click();
            cy.get('[data-e2e="add-factor-u2f"]').should('be.visible').click();
            const factorName = 'virtualAuthenticator';
            cy.get('[data-e2e="u2f-factor-name"]').should('be.visible').type(factorName);
            cy.task('remoteDebuggerCommand', {
              event: 'WebAuthn.enable',
            });
            cy.task('remoteDebuggerCommand', {
              event: 'WebAuthn.addVirtualAuthenticator',
              params: {
                options: {
                  protocol: 'ctap2',
                  transport: 'usb',
                  hasResidentKey: true,
                  hasUserVerification: true,
                  isUserVerified: true,
                },
              },
            });
            cy.get('[data-e2e="save-factor"]').should('be.visible').click();
            cy.contains('[data-e2e="u2f-factor-names"]', factorName);
            login(postAuthU2FEmail, userPw, target.headers['x-zitadel-orgid'], () => {
              cy.get('#btn-login').should('be.visible').click();
            });
            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.get('@userId').then((userId) => {
              cy.visit(`/users/${userId}?org=${target.headers['x-zitadel-orgid']}&id=metadata`);
              cy.contains('tr', 'last username used').contains(postAuthU2FEmail);
            });
          });
        },
      );
    });
    describe('passwordless', () => {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureHumanDoesntExist(target, postAuthPWLessEmail);
          ensureHumanExists(target, postAuthPWLessEmail).as('userId');
        });
      });
      it(
        'should store the username in metadata after passwordless authentication',
        {
          // Verifying a key that was registered in an origin other than the RP's origin fails.
          // It is tagged here so that it can be grepped and skipped when run against a dev server.
          tags: ['@same-origin'],
        },
        () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            login(postAuthPWLessEmail, userPw, target.headers['x-zitadel-orgid']);
            cy.visit('/users/me?id=security');
            cy.get('[data-e2e="add-passwordless"]').should('be.visible').click();
            const pwlessName = 'virtualPasswordless';
            cy.get('[data-e2e="passwordless-name"]').should('be.visible').type(pwlessName);
            cy.task('remoteDebuggerCommand', {
              event: 'WebAuthn.enable',
            });
            cy.task('remoteDebuggerCommand', {
              event: 'WebAuthn.addVirtualAuthenticator',
              params: {
                options: {
                  protocol: 'ctap2',
                  transport: 'usb',
                  hasResidentKey: true,
                  hasUserVerification: true,
                  isUserVerified: true,
                },
              },
            });
            cy.get('[data-e2e="passwordless-new"]').should('be.visible').click()
            cy.contains('[data-e2e="passwordless-names"]', pwlessName);
            login(postAuthPWLessEmail, userPw, target.headers['x-zitadel-orgid'], undefined, true);
            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.get('@userId').then((userId) => {
              cy.visit(`/users/${userId}?org=${target.headers['x-zitadel-orgid']}&id=metadata`);
              cy.contains('tr', 'last username used').contains(postAuthPWLessEmail);
            });
          });
        },
      );
    });
  });
});
