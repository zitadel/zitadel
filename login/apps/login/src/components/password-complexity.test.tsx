import {
  cleanup,
  render,
  screen,
  waitFor,
  within,
} from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test } from "vitest";
import { PasswordComplexity } from "./password-complexity";

const matchesTitle = `Matches`;
const doesntMatchTitle = `Doesn't match`;

describe("<PasswordComplexity/>", () => {
  describe.each`
    settingsMinLength | password        | expectSVGTitle
    ${5}              | ${"Password1!"} | ${matchesTitle}
    ${30}             | ${"Password1!"} | ${doesntMatchTitle}
    ${0}              | ${"Password1!"} | ${matchesTitle}
    ${undefined}      | ${"Password1!"} | ${false}
  `(
    `With settingsMinLength=$settingsMinLength, password=$password, expectSVGTitle=$expectSVGTitle`,
    ({ settingsMinLength, password, expectSVGTitle }) => {
      const feedbackElementLabel = /password length/i;
      beforeEach(() => {
        render(
          <PasswordComplexity
            password={password}
            equals
            passwordComplexitySettings={{
              minLength: settingsMinLength,
              requiresLowercase: false,
              requiresUppercase: false,
              requiresNumber: false,
              requiresSymbol: false,
              resourceOwnerType: 0, // ResourceOwnerType.RESOURCE_OWNER_TYPE_UNSPECIFIED,
            }}
          />,
        );
      });
      afterEach(cleanup);

      if (expectSVGTitle === false) {
        test(`should not render the feedback element`, async () => {
          await waitFor(() => {
            expect(
              screen.queryByText(feedbackElementLabel),
            ).not.toBeInTheDocument();
          });
        });
      } else {
        test(`Should show one SVG with title ${expectSVGTitle}`, async () => {
          await waitFor(async () => {
            const svg = within(
              screen.getByText(feedbackElementLabel)
                .parentElement as HTMLElement,
            ).findByRole("img");
            expect(await svg).toHaveTextContent(expectSVGTitle);
          });
        });
      }
    },
  );
});
