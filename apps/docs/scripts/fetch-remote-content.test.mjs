
import { test, describe, beforeEach, afterEach, mock } from 'node:test';
import assert from 'node:assert';
import { getCurrentRef } from './fetch-remote-content.mjs';

describe('fetch-remote-content', () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
    mock.restoreAll();
  });

  describe('getCurrentRef', () => {
    test('returns VERCEL_GIT_COMMIT_REF if set', () => {
      process.env.VERCEL_GIT_COMMIT_REF = 'vercel-branch';
      delete process.env.GITHUB_REF_NAME;
      assert.strictEqual(getCurrentRef(true), 'vercel-branch');
    });

    test('returns GITHUB_REF_NAME if set', () => {
      delete process.env.VERCEL_GIT_COMMIT_REF;
      process.env.GITHUB_REF_NAME = 'github-branch';
      assert.strictEqual(getCurrentRef(true), 'github-branch');
    });

    test('prefers VERCEL over GITHUB', () => {
      process.env.VERCEL_GIT_COMMIT_REF = 'vercel-branch';
      process.env.GITHUB_REF_NAME = 'github-branch';
      assert.strictEqual(getCurrentRef(true), 'vercel-branch');
    });

    // Note: Testing the git fallback would require mocking child_process which is complex via ES modules mocking 
    // without a loader or specialized tool. For now, we trust the env priority logic which is the core fix.
  });
});
