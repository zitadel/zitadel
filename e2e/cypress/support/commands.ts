import 'cypress-wait-until';
import { apiAuth, systemAuth } from './api/apiauth';
import { API, SystemAPI } from './api/types';
import { ensureQuotaIsRemoved, Unit } from './api/quota';
import { instanceUnderTest } from './api/instances';

interface ShouldNotExistOptions {
  selector: string;
  timeout?: {
    errMessage: string;
    ms: number;
  };
}

// Goal is to reduce the speed of operations executed to better mimic user interaction
const COMMAND_DELAY = Cypress.env('COMMAND_DELAY') || 0;
if (COMMAND_DELAY > 0) {
    for (const command of ['visit', 'click', 'trigger', 'reload']) {
        Cypress.Commands.overwrite(command as unknown as keyof Cypress.Chainable<any>, (originalFn, ...args) => {
            const origVal = originalFn(...args);

            return new Promise((resolve) => {
              setTimeout(() => {}, COMMAND_DELAY/2);
              resolve(origVal);
              setTimeout(() => {}, COMMAND_DELAY/2);
            });
        });
    }
}

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Custom command that asserts on clipboard text.
       *
       * @example cy.clipboardMatches('hodor', 'hodor1234')
       */
      clipboardMatches(pattern: RegExp | string): Cypress.Chainable<null>;

      /**
       * Custom command that waits until the selector finds zero elements.
       */
      shouldNotExist(options?: ShouldNotExistOptions): Cypress.Chainable<null>;

      /**
       * Custom command that ensures a reliable testing context and returns it
       */
      context(): Cypress.Chainable<Context>;

      /**
       * Custom command that asserts success is printed after a change.
       */
      shouldConfirmSuccess(): Cypress.Chainable<null>;
    }
  }
}

export interface Context {
  api: API;
  system: SystemAPI;
  instanceId: number;
}

Cypress.Commands.add('clipboardMatches', { prevSubject: false }, (pattern: RegExp | string) => {
  /* doesn't work reliably
    return cy.window()
        .then(win => {
            win.focus()
            return cy.waitUntil(() => win.navigator.clipboard.readText()
                .then(clipboadText => {
                    win.focus()
                    const matches = typeof pattern === "string"
                    ? clipboadText.includes(pattern)
                    : pattern.test(clipboadText)
                    if (!matches) {
                        cy.log(`text in clipboard ${clipboadText} doesn't match the pattern ${pattern}, yet`)
                    }
                    return matches
                })
            )
        })
        .then(() => null)
    */
});

Cypress.Commands.add('shouldNotExist', { prevSubject: false }, (options: ShouldNotExistOptions) => {
  if (!options.timeout) {
    const elements = Cypress.$(options.selector);
    expect(elements.text()).to.be.empty;
    expect(elements.length).to.equal(0);
    return null;
  }
  return cy
    .waitUntil(
      () => {
        const elements = Cypress.$(options.selector);
        if (!elements.length) {
          return cy.wrap(true);
        }
        return cy.log(`elements with selector ${options.selector} and text ${elements.text()} exist`).wrap(false);
      },
      {
        timeout: options.timeout.ms,
        errorMsg: options.timeout.errMessage,
      },
    )
    .then(() => null);
});

Cypress.Commands.add('shouldConfirmSuccess', { prevSubject: false }, () => {
  cy.get('.data-e2e-message');
  cy.shouldNotExist({ selector: '.data-e2e-failure' });
  cy.get('.data-e2e-success');
});

Cypress.Commands.add('context', { prevSubject: false }, () => {
  return systemAuth().then((system) => {
    return instanceUnderTest(system).then((instanceId) => {
      return ensureQuotaIsRemoved(
        {
          system: system,
          api: null,
          instanceId: instanceId,
        },
        Unit.AuthenticatedRequests,
      ).then(() => {
        return ensureQuotaIsRemoved(
          {
            system: system,
            api: null,
            instanceId: instanceId,
          },
          Unit.ExecutionSeconds,
        ).then(() => {
          return apiAuth().then((api) => {
            return {
              system: system,
              api: api,
              instanceId: instanceId,
            };
          });
        });
      });
    });
  });
});
