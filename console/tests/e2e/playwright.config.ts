import { PlaywrightTestConfig } from '@playwright/test';

// Environment variables should be typed and initially loaded

export const E2E_ORG_OWNER_PW=process.env.E2E_ORG_OWNER_PW
export const E2E_ORG_OWNER_VIEWER_PW=process.env.E2E_ORG_OWNER_VIEWER_PW
export const E2E_ORG_PROJECT_CREATOR_PW=process.env.E2E_ORG_PROJECT_CREATOR_PW
export const E2E_SERVICEACCOUNT_KEY=process.env.E2E_SERVICEACCOUNT_KEY
export const E2E_CONSOLE_URL=process.env.E2E_CONSOLE_URL
export const E2E_API_CALLS_DOMAIN=process.env.E2E_API_CALLS_DOMAIN
export const E2E_ZITADEL_PROJECT_RESOURCE_ID=process.env.E2E_ZITADEL_PROJECT_RESOURCE_ID

export const RESULTSPATH='./tests/e2e/results'

const config: PlaywrightTestConfig = {    
  use: {
    ignoreHTTPSErrors: true,
    video: 'retain-on-failure',
    baseURL: E2E_CONSOLE_URL,    
    contextOptions: {
        locale: "en-US",
        recordVideo: {
            // should be overwritten when creating browser context:
            dir: RESULTSPATH,
        },
        recordHar: {
            // should be overwritten when creating browser context:
            path: RESULTSPATH
        }
    }
  },
};
export default config;
