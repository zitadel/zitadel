import 'cypress-wait-until';
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
  selector?: string;
  timeout?: number;
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

Cypress.Commands.add('shouldNotExist', { prevSubject: false }, (options?: ShouldNotExistOptions) => {
  return cy.waitUntil(() => Cypress.$(options?.selector).length === 0, {
    errorMsg: () => `Timed out while waiting for element to not exist: ${Cypress.$(options?.selector).text()}`,
    timeout: typeof options?.timeout === 'number' ? options.timeout : 500,
  });
});
