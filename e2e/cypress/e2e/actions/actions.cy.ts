import { ensureActionExists, resetAllTriggers, triggerActions } from 'support/api/actions';
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
  const specPath = 'cypress/e2e/actions';
  const emailVerifiedScript = 'e2eSetEmailVerified';
  const addGrantScript = 'e2eAddGrant';
  const setMetadataScript = 'e2eSetMetadata';
  const projectName = 'e2eaction';
  const roleKey = 'e2eactionrole';
  const preCreationEmail = 'precre@action.com';
  const postCreationEmail = 'postcre@action.com';
  const postAuthPWEmail = 'postauthpw@action.com';
  const postAuthOTPEmail = 'postauthotp@action.com';
  const postAuthU2FEmail = 'postauthu2f@action.com';
  const postAuthPWLessEmail = 'postauthpwless@action.com';

  beforeEach(() => {
    newTarget('e2eactions')
      .as('target')
      .then((target) => {
        ensureDomainPolicy(target, false, false, false);
        cy.readFile(`${specPath}/${emailVerifiedScript}.js`).then((script) => {
          ensureActionExists(target, emailVerifiedScript, script)
            .as('emailVerifiedId')
            .then((emailVerifiedId) => {
              resetAllTriggers(target);
              triggerActions(target, 3, 2, [emailVerifiedId]);
            });
        });
      });
  });
  describe('pre creation', () => {
    const validKey = `a valid key`;
    const validValue = `a valid value`;
    const tooLongKey = `${'toolongkey'.repeat(20)}-overflow`;
    const mustNotExist = 'must not exist';
    const validActionName = 'e2eSetValidMetadata';
    const invalidActionName = 'e2eSetInvalidMetadata';

    beforeEach(() => {
      cy.get<ZITADELTarget>('@target').then((target) => {
        ensureHumanDoesntExist(target, preCreationEmail);
        cy.get<number>('@emailVerifiedId').then((emailVerifiedId) => {
          cy.readFile(`${specPath}/${setMetadataScript}.js`).then((metadataScript) => {
            ensureActionExists(
              target,
              validActionName,
              metadataScript
                .replace('<IDENTIFIER>', validActionName)
                .replace('<KEY>', validKey)
                .replace('<VALUE>', validValue),
              false,
            ).then((setValidMetaId) => {
              ensureActionExists(
                target,
                invalidActionName,
                metadataScript
                  .replace('<IDENTIFIER>', invalidActionName)
                  .replace('<KEY>', tooLongKey)
                  .replace('<VALUE>', mustNotExist),
                true,
              ).then((setInvalidMetaId) => {
                triggerActions(target, 3, 2, [emailVerifiedId, setValidMetaId, setInvalidMetaId]);
              });
            });
          });
        });
      });
    });
    it(`shouldn't prompt for email code and add metadata`, () => {
      cy.get<ZITADELTarget>('@target').then((target) => {
        register(preCreationEmail, target.orgId).then((userId) => {
          sessionAsPredefinedUser(User.IAMAdminUser);
          cy.visit(`/users/${userId}?id=metadata&org=${target.orgId}`);
          cy.contains('tr', 'akey').contains('avalue');
          cy.contains('tr', mustNotExist).should('not.exist')
          cy.contains('tr', tooLongKey).should('not.exist');
        });
      });
    });
  });
  describe('post creation', () => {
    beforeEach(() => {
      cy.get<ZITADELTarget>('@target').then((target) => {
        cy.get<number>('@emailVerifiedId').then((emailVerifiedId) => {
          triggerActions(target, 3, 2, [emailVerifiedId]);
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
              triggerActions(target, 3, 3, [addGrantId]);
            });
          });
        });
      });
    });
    it(`should add a grant when registering via UI`, () => {
      cy.get<ZITADELTarget>('@target').then((target) => {
        register(postCreationEmail, target.orgId).then((userId) => {
          sessionAsPredefinedUser(User.IAMAdminUser);
          cy.visit(`/users/${userId}?id=grants&org=${target.orgId}`);
          cy.contains('tr', roleKey);
        });
      });
    });
  });
  describe('post authentication', () => {
    beforeEach(() => {
      cy.get<ZITADELTarget>('@target').then((target) => {
        cy.readFile<string>(`${specPath}/${setMetadataScript}.js`).then((script) => {
          const setAuthError = 'setAuthError';
          ensureActionExists(
            target,
            setAuthError,
            script
              .replace('<IDENTIFIER>', setAuthError)
              .replace('<KEY>', '${ctx.v1.authMethod} error')
              .replace('<VALUE>', '${ctx.v1.authError}'),
          ).then((addMetadataId) => {
            triggerActions(target, 3, 1, [addMetadataId]);
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

      it('should store password error none in metadata after successful password authentication', () => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          login(postAuthPWEmail, target.orgId);
          cy.get('@userId').then((userId) => {
            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.visit(`/users/${userId}?id=metadata&org=${target.orgId}`);
            cy.get('tr').should('have.length', 2);
            cy.contains('tr', 'password error').contains('none');
          });
        });
      });

      it('should store password error authentication failed in metadata after failed password authentication', () => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          login(postAuthPWEmail, target.orgId, false, 'this password is wrong');
          cy.get('@userId').then((userId) => {
            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.visit(`/users/${userId}?id=metadata&org=${target.orgId}`);
            cy.get('tr').should('have.length', 2);
            cy.contains('tr', 'password error').contains('Errors.User.Password.Invalid');
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

      it(
        'should store password error none and OTP error none in metadata after successful otp authentication',
        {
          // authentication can fail sometimes
          retries: {
            openMode: null,
            runMode: 2,
          },
        },
        () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            login(postAuthOTPEmail, target.orgId);
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
            login(postAuthOTPEmail, target.orgId, true, undefined, () => {
              cy.task<string>('generateOTP').then((token) => {
                cy.get('#code').should('be.visible').type(token);
                cy.get('#submit-button').should('be.visible').click();
              });
            });

            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.get('@userId').then((userId) => {
              cy.visit(`/users/${userId}?org=${target.orgId}&id=metadata`);
              cy.get('tr').should('have.length', 3);
              cy.contains('tr', 'password error').contains('none');
              cy.contains('tr', 'OTP error').contains('none');
            });
          });
        },
      );
    });
    describe('u2f', () => {
      beforeEach(() => {
        cy.get<ZITADELTarget>('@target').then((target) => {
          ensureHumanDoesntExist(target, postAuthU2FEmail);
          ensureHumanExists(target, postAuthU2FEmail).as('userId');
        });
      });
      it(
        'should store password error none and U2F error none in metadata after successful u2f authentication',
        {
          // Verifying a key that was registered in an origin other than the RP's origin fails.
          // It is tagged here so that it can be grepped and skipped when run against a dev server.
          tags: ['@same-origin'],
          browser: 'chrome',
          // authentication can fail sometimes
          retries: {
            openMode: null,
            runMode: 2,
          },
        },
        () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            login(postAuthU2FEmail, target.orgId);
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
            login(postAuthU2FEmail, target.orgId, true, undefined, () => {
              cy.get('#btn-login').should('be.visible').click();
            });
            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.get('@userId').then((userId) => {
              cy.visit(`/users/${userId}?org=${target.orgId}&id=metadata`);
              cy.get('tr').should('have.length', 3);
              cy.contains('tr', 'password error').contains('none');
              cy.contains('tr', 'U2F error').contains('none');
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
        'should store password error none and passwordless error none in metadata after successful passwordless authentication',
        {
          // Verifying a key that was registered in an origin other than the RP's origin fails.
          // It is tagged here so that it can be grepped and skipped when run against a dev server.
          tags: ['@same-origin'],
          browser: 'chrome',
          // authentication can fail sometimes
          retries: {
            openMode: null,
            runMode: 2,
          },
        },
        () => {
          cy.get<ZITADELTarget>('@target').then((target) => {
            login(postAuthPWLessEmail, target.orgId);
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
            cy.get('[data-e2e="passwordless-new"]').should('be.visible').click();
            cy.contains('[data-e2e="passwordless-names"]', pwlessName);
            login(postAuthPWLessEmail, target.orgId, true, undefined, undefined, true);
            sessionAsPredefinedUser(User.IAMAdminUser);
            cy.get('@userId').then((userId) => {
              cy.visit(`/users/${userId}?org=${target.orgId}&id=metadata`);
              cy.get('tr').should('have.length', 3);
              // TODO: Is it wrong that the action runs twice here?
              cy.contains('tr', 'password error').contains('none');
              cy.contains('tr', 'passwordless error').contains('none');
            });
          });
        },
      );
    });
  });
});
