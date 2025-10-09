import { mergeTests } from '@playwright/test';
import { test as apiTest } from './api.js';
import { test as userRegistratorTest } from './user-registrator.js';
import { test as userCreatorTest } from './user-creator.js';

export const test = mergeTests(apiTest, userRegistratorTest, userCreatorTest);
