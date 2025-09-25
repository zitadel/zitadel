import { defineConfig } from "cypress";
import { unlinkSync } from 'fs'

export default defineConfig({
  reporter: "list",
  video: true,
  retries: {
    runMode: 2
  },
  e2e: {
    baseUrl: `http://localhost:3000${process.env.NEXT_PUBLIC_BASE_PATH || ""}`,
    specPattern: "integration/integration/**/*.cy.{js,jsx,ts,tsx}",
    supportFile: "integration/support/e2e.{js,jsx,ts,tsx}",
    pageLoadTimeout: 120_0000,
    env: {
      API_MOCK_STUBS_URL: process.env.API_MOCK_STUBS_URL || "http://localhost:22220/v1/stubs"
    },
    setupNodeEvents(on, config) {
      on('after:spec', (_, results) => {
        // We don't want to keep and cache videos of successful runs
        // This is an implementation according to the official Cypress docs:
        // https://docs.cypress.io/app/guides/screenshots-and-videos#Delete-videos-for-specs-without-failing-or-retried-tests
        if (results && results.video) {
          // Do we have failures for any retry attempts?
          const failures = results.tests.some((test) =>
            test.attempts.some((attempt) => attempt.state === 'failed')
          )
          if (!failures) {
            // delete the video if the spec passed and no tests retried
            try {
              unlinkSync(results.video)
            } catch (err) {
              // Ignore errors when deleting video file
            }
          }
        }
      })      
    },
  },
});
