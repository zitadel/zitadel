import 'cypress-wait-until';
import { apiAuth, systemAuth } from './api/apiauth';
import { API, SystemAPI } from './api/types';
import { ensureQuotaIsRemoved, Unit } from './api/quota';
import { instanceUnderTest } from './api/instances';
//
//namespace Cypress {
//    interface Chainable {
//        /**
//         * Custom command that authenticates a user.
//         *
//         * @example cy.consolelogin('hodor', 'hodor1234')
//         */
//        consolelogin(username: string, password: string): void
//    }
//}
//
//Cypress.Commands.add('consolelogin', { prevSubject: false }, (username: string, password: string) => {
//
//    window.sessionStorage.removeItem("zitadel:access_token")
//    cy.visit(Cypress.config('baseUrl')/ui/console).then(() => {
//        // fill the fields and push button
//        cy.get('#loginName').type(username, { log: false })
//        cy.get('#submit-button').click()
//        cy.get('#password').type(password, { log: false })
//        cy.get('#submit-button').click()
//        cy.location('pathname', {timeout: 5 * 1000}).should('eq', '/');
//    })
//})
//

interface ShouldNotExistOptions {
  selector: string;
  timeout?: {
    errMessage: string;
    ms: number;
  };
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
       * Custom command that has to be called before each test
       */
      resetContext(): Cypress.Chainable<null>;
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
/*
Cypress.Commands.add('authenticate', {prevSubject:false}, ()=>{
  return systemAuth().then((system) => {
    return apiAuth().then((api) => {
      return {
        api: api,
        system: system,
      }
    });
})
})*/

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
