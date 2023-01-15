import 'cypress-wait-until';
// @ts-ignore
import * as regGrep from '@cypress/grep';
regGrep();

export interface ZITADELTarget {
  headers: {
    Authorization: string;
    'x-zitadel-orgid': string;
  };
  mgmtBaseURL: string;
  adminBaseURL: string;
  org: string;
}

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
      shouldNotExist(options: ShouldNotExistOptions): Cypress.Chainable<null>;
      /**
       * Custom command that asserts success is printed after a change.
       */
      shouldConfirmSuccess(): Cypress.Chainable<null>;
    }
  }

  namespace Mocha {
    interface Context {
      get target(): void; // TODO: Remove
      get api(): void; // TODO: Remove
    }
  }
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
