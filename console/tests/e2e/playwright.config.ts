import { PlaywrightTestConfig } from '@playwright/test';

// Environment variables should be typed and initially loaded

export const ORG_OWNER_PW=process.env.ORG_OWNER_PW
export const ORG_OWNER_VIEWER_PW=process.env.ORG_OWNER_VIEWER_PW
export const ORG_PROJECT_CREATOR_PW=process.env.ORG_PROJECT_CREATOR_PW
export const SERVICEACCOUNT_KEY=process.env.SERVICEACCOUNT_KEY
export const CONSOLE_URL=process.env.CONSOLE_URL
export const API_CALLS_DOMAIN=process.env.API_CALLS_DOMAIN
export const ZITADEL_PROJECT_RESOURCE_ID=process.env.ZITADEL_PROJECT_RESOURCE_ID

export const RESULTSPATH='./tests/e2e/results'

const config: PlaywrightTestConfig = {    
  use: {
    ignoreHTTPSErrors: true,
    video: 'retain-on-failure',
    baseURL: CONSOLE_URL,    
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
