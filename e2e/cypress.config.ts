import { defineConfig } from 'cypress';

import * as webpack from '@cypress/webpack-batteries-included-preprocessor';
import { addCucumberPreprocessorPlugin } from '@badeball/cypress-cucumber-preprocessor';

let tokensCache = new Map<string, string>();

export default defineConfig({
  reporter: 'mochawesome',

  reporterOptions: {
    reportDir: 'cypress/results',
    overwrite: false,
    html: true,
    json: true,
  },

  chromeWebSecurity: false,
  trashAssetsBeforeRuns: false,
  defaultCommandTimeout: 10000,

  env: {
    ORGANIZATION: process.env.CYPRESS_ORGANIZATION || 'zitadel',
    BACKEND_URL: process.env.CYPRESS_BACKEND_URL || baseUrl().replace('/ui/console', ''),
  },

  e2e: {
    baseUrl: baseUrl(),
    experimentalSessionAndOrigin: true,
    specPattern: '**/*.{feature,cy.ts}',
    setupNodeEvents,
  },
});

function baseUrl() {
  return process.env.CYPRESS_BASE_URL || 'http://localhost:8080/ui/console';
}

async function setupNodeEvents(on, config) {
  await addCucumberPreprocessorPlugin(on, config);

  on('task', {
    safetoken({ key, token }) {
      tokensCache.set(key, token);
      return null;
    },
  });
  on('task', {
    loadtoken({ key }): string | null {
      return tokensCache.get(key) || null;
    },
  });

  const preprocessorConfig = {
    ...webpack.defaultOptions,
    typescript: require.resolve('typescript'),
  };

  preprocessorConfig.webpackOptions.module.rules.push({
    test: /\.feature$/,
    use: [
      {
        loader: '@badeball/cypress-cucumber-preprocessor/webpack',
        options: config,
      },
    ],
  });

  on('file:preprocessor', webpack(preprocessorConfig));

  return config;
}
