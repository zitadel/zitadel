import assert from 'node:assert/strict';
import { describe, test } from 'node:test';
import {
  getVersionPreferenceCookie,
  LATEST_VERSION_LABEL,
  resolveVersion,
} from '../lib/version-preference.mjs';

const versions = [
  { param: 'latest' },
  { param: 'newest-release' },
  { param: 'older-release' },
];

describe('API docs version preference', () => {
  test('labels the unreleased docs consistently', () => {
    assert.equal(LATEST_VERSION_LABEL, 'Latest (unreleased)');
  });

  test('creates a preference cookie for an explicit selection', () => {
    assert.equal(
      getVersionPreferenceCookie('latest'),
      'zitadel-docs-version=latest; Path=/docs; Max-Age=31536000; SameSite=Lax',
    );
  });

  test('defaults a new visitor to the latest released version', () => {
    assert.equal(resolveVersion({ versions }), 'newest-release');
  });

  test('preserves an explicitly selected released version', () => {
    assert.equal(resolveVersion({ storedVersion: 'older-release', versions }), 'older-release');
  });

  test('preserves an explicitly selected unreleased version', () => {
    assert.equal(resolveVersion({ storedVersion: 'latest', versions }), 'latest');
  });

  test('falls back from a stale preference to the latest released version', () => {
    assert.equal(resolveVersion({ storedVersion: 'removed-release', versions }), 'newest-release');
  });

  test('gives a URL version precedence over the stored preference', () => {
    assert.equal(
      resolveVersion({
        requestedVersion: 'older-release',
        storedVersion: 'latest',
        versions,
      }),
      'older-release',
    );
  });
});
