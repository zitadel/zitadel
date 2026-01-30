import { test, describe, beforeEach, afterEach, mock } from 'node:test';
import assert from 'node:assert';
import fs from 'fs';
import path from 'path';
import { getCurrentRef, resetCache, downloadFileContent, isValidRef, safeLog } from './fetch-remote-content.mjs';

const TEST_TMP_DIR = path.join(process.cwd(), '.test-tmp');
const MOCK_REPO_ROOT = path.join(TEST_TMP_DIR, 'mock-repo');

describe('fetch-remote-content', () => {
    
  // --- Setup & Teardown ---
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
    resetCache();
    
    // Set up a mock repo structure
    if (fs.existsSync(TEST_TMP_DIR)) {
      fs.rmSync(TEST_TMP_DIR, { recursive: true, force: true });
    }
    fs.mkdirSync(MOCK_REPO_ROOT, { recursive: true });
    
    // Create some dummy content
    fs.writeFileSync(path.join(MOCK_REPO_ROOT, 'README.md'), '# Mock Repo');
    fs.mkdirSync(path.join(MOCK_REPO_ROOT, 'secrets'), { recursive: true });
    fs.writeFileSync(path.join(MOCK_REPO_ROOT, 'secrets/passwd'), 'super_secret');
  });

  afterEach(() => {
    process.env = originalEnv;
    mock.restoreAll();
    resetCache();
    if (fs.existsSync(TEST_TMP_DIR)) {
      fs.rmSync(TEST_TMP_DIR, { recursive: true, force: true });
    }
  });

  // --- Helper Tests ---
  describe('isValidRef', () => {
      test('accepts valid alphanumeric refs', () => {
          assert.strictEqual(isValidRef('main'), true);
          assert.strictEqual(isValidRef('v1.0.0'), true);
          assert.strictEqual(isValidRef('feature/new-docs'), true);
          assert.strictEqual(isValidRef('fix_bug-123'), true);
      });

      test('rejects malicious refs', () => {
          assert.strictEqual(isValidRef('; rm -rf'), false);
          assert.strictEqual(isValidRef('main && echo pwned'), false);
          assert.strictEqual(isValidRef('../parent'), false); // contains .. which isn't allowed in regex
      });
  });

  describe('safeLog', () => {
      test('strips newlines', () => {
          assert.strictEqual(safeLog('valid'), 'valid');
          assert.strictEqual(safeLog('malicious\nlog'), 'maliciouslog');
          assert.strictEqual(safeLog('malicious\r\nlog'), 'maliciouslog');
      });
  });

  // --- getCurrentRef Tests ---
  describe('getCurrentRef', () => {
    test('returns VERCEL_GIT_COMMIT_REF if set and valid', () => {
      process.env.VERCEL_GIT_COMMIT_REF = 'vercel-branch';
      delete process.env.GITHUB_REF_NAME;
      assert.strictEqual(getCurrentRef(), 'vercel-branch');
    });

    test('ignores invalid VERCEL_GIT_COMMIT_REF', () => {
        // Should fall through to GITHUB or git command
        process.env.VERCEL_GIT_COMMIT_REF = 'invalid;cmd'; 
        process.env.GITHUB_REF_NAME = 'valid-github';
        assert.strictEqual(getCurrentRef(), 'valid-github');
    });

    test('returns GITHUB_REF_NAME if set and valid', () => {
      delete process.env.VERCEL_GIT_COMMIT_REF;
      process.env.GITHUB_REF_NAME = 'github-branch';
      assert.strictEqual(getCurrentRef(), 'github-branch');
    });
  });

  // --- downloadFileContent Tests ---
  describe('downloadFileContent', () => {
    
    test('refuses traversal attempts (..)', async () => {
        // We expect this to return null and log a warning
        const content = await downloadFileContent('main', '../secrets/passwd');
        assert.strictEqual(content, null);
    });
    
    test('refuses absolute paths', async () => {
        const content = await downloadFileContent('main', '/etc/passwd');
        assert.strictEqual(content, null);
    });

    test('refuses encoded traversal attempts', async () => {
        // %2e%2e is ..
        const content = await downloadFileContent('main', '%2e%2e/secrets/passwd');
        assert.strictEqual(content, null);
    });

    test('reads valid local file (README.md)', async () => {
       // Mock getCurrentRef to return 'main' and pass 'main' to simulate local
       process.env.VERCEL_GIT_COMMIT_REF = 'main';
       
       const content = await downloadFileContent('main', 'README.md');
       // It should find the real README.md of the project since the script resolves from __dirname
       assert.ok(content); 
       assert.ok(content.length > 0);
    });
    
    test('refuses invalid ref strings', async () => {
        // Validation check for command injection protection
        const content = await downloadFileContent('invalid; rm -rf', 'README.md');
        assert.strictEqual(content, null);
    });
  });

});
