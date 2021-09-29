import { PlaywrightTestConfig } from '@playwright/test';
import { CONSOLE_URL } from './models/env';
const config: PlaywrightTestConfig = {    
  use: {
    ignoreHTTPSErrors: true,
    video: 'retain-on-failure',
    baseURL: CONSOLE_URL,    
    contextOptions: {
        locale: "en-US",
        recordVideo: {
             // should be overwritten when creating browser context:
            dir: './tests/e2e/videos/',
        },        
    }
  },
};
export default config;
