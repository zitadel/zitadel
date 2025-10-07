import { mergeTests } from '@playwright/test';
import { test as apiTest } from './api.js';
import { test as registeredTest } from './user-anonymous.js';
import { test as anonymousTest } from './user-creator.js';

export const test = mergeTests(apiTest, registeredTest, anonymousTest);