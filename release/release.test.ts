/**
 * Unit tests for the release script.
 *
 * This test suite covers the testable functions extracted from release.ts including:
 * - Branch pattern validation for conventional commits
 * - Environment variable setup based on different release scenarios
 * - Docker build command construction
 * - Main release orchestration logic
 */

import { execSync } from 'node:child_process';
import { releaseVersion, releaseChangelog, releasePublish } from 'nx/release';
import { afterEach, beforeEach, describe, expect, test, vi } from 'vitest';
import {
  shouldUseConventionalCommits,
  setupEnvironmentVariables,
  executeDockerBuild,
  parseReleaseOptions,
  executeRelease,
  type GitInfo,
  type ReleaseOptions,
  type EnvironmentConfig,
} from './release';

// Mock external dependencies
vi.mock('node:child_process', () => ({
  execSync: vi.fn(),
}));

vi.mock('nx/release', () => ({
  releaseVersion: vi.fn(),
  releaseChangelog: vi.fn(),
  releasePublish: vi.fn(),
}));

// Store original environment
const originalEnv = process.env;

describe('shouldUseConventionalCommits', () => {
  test('returns true for major maintenance branch', () => {
    expect(shouldUseConventionalCommits('v1.x')).toBe(true);
    expect(shouldUseConventionalCommits('v2.x')).toBe(true);
  });

  test('returns true for minor maintenance branch', () => {
    expect(shouldUseConventionalCommits('v1.0.x')).toBe(true);
    expect(shouldUseConventionalCommits('v2.15.x')).toBe(true);
  });

  test('returns false for main branch', () => {
    expect(shouldUseConventionalCommits('main')).toBe(false);
  });

  test('returns false for feature branches', () => {
    expect(shouldUseConventionalCommits('feature/my-feature')).toBe(false);
    expect(shouldUseConventionalCommits('fix/bug-123')).toBe(false);
  });
});

describe('setupEnvironmentVariables', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    process.env = { ...originalEnv };
    // Clear relevant env vars
    delete process.env.ZITADEL_RELEASE_VERSION;
    delete process.env.ZITADEL_RELEASE_REVISION;
    delete process.env.ZITADEL_RELEASE_IS_LATEST;
    delete process.env.ZITADEL_RELEASE_PUSH;
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  const mockGitInfo: GitInfo = {
    branch: 'main',
    sha: 'abc123def456',
    highestVersionBefore: '1.0.0',
  };

  const mockConfig: EnvironmentConfig = {
    versionEnvVar: 'ZITADEL_RELEASE_VERSION',
    revisionEnvVar: 'ZITADEL_RELEASE_REVISION',
    isLatestEnvVar: 'ZITADEL_RELEASE_IS_LATEST',
    doPushEnvVar: 'ZITADEL_RELEASE_PUSH',
  };

  const mockOptions: ReleaseOptions = {
    dryRun: false,
    verbose: false,
  };

  test('sets revision and push env vars for all scenarios', () => {
    setupEnvironmentVariables(mockConfig, mockGitInfo, mockOptions, false);

    expect(process.env.ZITADEL_RELEASE_REVISION).toBe('abc123def456');
    expect(process.env.ZITADEL_RELEASE_PUSH).toBe('true');
  });

  test('sets push to false when dryRun is true', () => {
    setupEnvironmentVariables(mockConfig, mockGitInfo, { ...mockOptions, dryRun: true }, false);

    expect(process.env.ZITADEL_RELEASE_PUSH).toBe('false');
  });

  test('sets version to SHA and isLatest to false for non-conventional commits', () => {
    setupEnvironmentVariables(mockConfig, mockGitInfo, mockOptions, false);

    expect(process.env.ZITADEL_RELEASE_VERSION).toBe('abc123def456');
    expect(process.env.ZITADEL_RELEASE_IS_LATEST).toBe('false');
  });

  test('sets version with v prefix for conventional commits with workspace version', () => {
    setupEnvironmentVariables(mockConfig, mockGitInfo, mockOptions, true, '2.0.0');

    expect(process.env.ZITADEL_RELEASE_VERSION).toBe('v2.0.0');
  });

  test('sets isLatest to true when new version is higher than before', () => {
    setupEnvironmentVariables(
      mockConfig,
      { ...mockGitInfo, highestVersionBefore: '1.0.0' },
      mockOptions,
      true,
      '2.0.0'
    );

    expect(process.env.ZITADEL_RELEASE_IS_LATEST).toBe('true');
  });

  test('sets isLatest to true when new version equals highest before', () => {
    setupEnvironmentVariables(
      mockConfig,
      { ...mockGitInfo, highestVersionBefore: '2.0.0' },
      mockOptions,
      true,
      '2.0.0'
    );

    expect(process.env.ZITADEL_RELEASE_IS_LATEST).toBe('true');
  });

  test('sets isLatest to false when new version is lower than before', () => {
    setupEnvironmentVariables(
      mockConfig,
      { ...mockGitInfo, highestVersionBefore: '3.0.0' },
      mockOptions,
      true,
      '2.0.0'
    );

    expect(process.env.ZITADEL_RELEASE_IS_LATEST).toBe('false');
  });
});

describe('executeDockerBuild', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test('includes build-docker-debug target for conventional commits', () => {
    const mockExecSync = vi.mocked(execSync);

    executeDockerBuild(true);

    expect(mockExecSync).toHaveBeenCalledWith(
      expect.stringContaining('build-docker build-docker-debug'),
      expect.any(Object)
    );
  });

  test('excludes build-docker-debug target for non-conventional commits', () => {
    const mockExecSync = vi.mocked(execSync);

    executeDockerBuild(false);

    expect(mockExecSync).toHaveBeenCalledWith(
      expect.not.stringContaining('build-docker-debug'),
      expect.any(Object)
    );
    expect(mockExecSync).toHaveBeenCalledWith(
      expect.stringContaining('--target build-docker --file'),
      expect.any(Object)
    );
  });

  test('includes all required bake files', () => {
    const mockExecSync = vi.mocked(execSync);

    executeDockerBuild(false);

    const call = mockExecSync.mock.calls[0][0] as string;
    expect(call).toContain('--file release/docker-bake-release.hcl');
    expect(call).toContain('--file apps/api/docker-bake-release.hcl');
    expect(call).toContain('--file apps/login/docker-bake-release.hcl');
  });

  test('passes environment variables to execSync', () => {
    const mockExecSync = vi.mocked(execSync);

    executeDockerBuild(false);

    expect(mockExecSync).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        stdio: 'inherit',
        env: process.env,
      })
    );
  });
});

describe('parseReleaseOptions', () => {
  test('parses dryRun flag', async () => {
    const options = await parseReleaseOptions(['--dryRun']);
    expect(options.dryRun).toBe(true);
  });

  test('parses verbose flag', async () => {
    const options = await parseReleaseOptions(['--verbose']);
    expect(options.verbose).toBe(true);
  });

  test('parses both flags', async () => {
    const options = await parseReleaseOptions(['--dryRun', '--verbose']);
    expect(options.dryRun).toBe(true);
    expect(options.verbose).toBe(true);
  });

  test('defaults dryRun to true', async () => {
    const options = await parseReleaseOptions([]);
    expect(options.dryRun).toBe(true);
  });

  test('defaults verbose to false', async () => {
    const options = await parseReleaseOptions([]);
    expect(options.verbose).toBe(false);
  });

  test('accepts --no-dryRun to set dryRun to false', async () => {
    const options = await parseReleaseOptions(['--no-dryRun']);
    expect(options.dryRun).toBe(false);
  });
});

describe('executeRelease', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    process.env = { ...originalEnv };
    delete process.env.ZITADEL_RELEASE_VERSION;
    delete process.env.ZITADEL_RELEASE_REVISION;
    delete process.env.ZITADEL_RELEASE_IS_LATEST;
    delete process.env.ZITADEL_RELEASE_PUSH;
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  const mockGitInfo: GitInfo = {
    branch: 'main',
    sha: 'abc123def456',
    highestVersionBefore: '1.0.0',
  };

  const mockOptions: ReleaseOptions = {
    dryRun: false,
    verbose: false,
  };

  const mockConfig: EnvironmentConfig = {
    versionEnvVar: 'ZITADEL_RELEASE_VERSION',
    revisionEnvVar: 'ZITADEL_RELEASE_REVISION',
    isLatestEnvVar: 'ZITADEL_RELEASE_IS_LATEST',
    doPushEnvVar: 'ZITADEL_RELEASE_PUSH',
  };

  test('returns 0 for non-conventional commits after docker build', async () => {
    const mockExecSync = vi.mocked(execSync);

    const exitCode = await executeRelease(mockGitInfo, mockOptions, mockConfig);

    expect(exitCode).toBe(0);
    expect(mockExecSync).toHaveBeenCalledOnce();
    expect(vi.mocked(releaseVersion)).not.toHaveBeenCalled();
  });

  test('calls releaseVersion for conventional commits', async () => {
    const mockGitInfo: GitInfo = {
      branch: 'v1.x',
      sha: 'abc123def456',
      highestVersionBefore: '1.0.0',
    };

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: {} as any,
    });

    vi.mocked(releasePublish).mockResolvedValue({});

    await executeRelease(mockGitInfo, mockOptions, mockConfig);

    expect(releaseVersion).toHaveBeenCalledWith({
      dryRun: false,
      verbose: false,
    });
  });

  test('throws error when workspace version cannot be determined', async () => {
    const mockGitInfo: GitInfo = {
      branch: 'v1.x',
      sha: 'abc123def456',
      highestVersionBefore: '1.0.0',
    };

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: undefined,
      projectsVersionData: {} as any,
    });

    await expect(
      executeRelease(mockGitInfo, mockOptions, mockConfig)
    ).rejects.toThrow('Could not determine workspace version');
  });

  test('returns 0 when all publish results succeed', async () => {
    const mockGitInfo: GitInfo = {
      branch: 'v1.x',
      sha: 'abc123def456',
      highestVersionBefore: '1.0.0',
    };

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: {} as any,
    });

    vi.mocked(releasePublish).mockResolvedValue({
      project1: { code: 0 },
      project2: { code: 0 },
    } as any);

    const exitCode = await executeRelease(mockGitInfo, mockOptions, mockConfig);

    expect(exitCode).toBe(0);
  });

  test('returns 1 when any publish result fails', async () => {
    const mockGitInfo: GitInfo = {
      branch: 'v1.x',
      sha: 'abc123def456',
      highestVersionBefore: '1.0.0',
    };

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: {} as any,
    });

    vi.mocked(releasePublish).mockResolvedValue({
      project1: { code: 0 },
      project2: { code: 1 },
    } as any);

    const exitCode = await executeRelease(mockGitInfo, mockOptions, mockConfig);

    expect(exitCode).toBe(1);
  });

  test('calls releaseChangelog with correct parameters', async () => {
    const mockGitInfo: GitInfo = {
      branch: 'v1.x',
      sha: 'abc123def456',
      highestVersionBefore: '1.0.0',
    };

    const mockProjectsVersionData = { someData: true };

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: mockProjectsVersionData as any,
    });

    vi.mocked(releasePublish).mockResolvedValue({});

    await executeRelease(mockGitInfo, mockOptions, mockConfig);

    expect(releaseChangelog).toHaveBeenCalledWith({
      versionData: mockProjectsVersionData,
      version: '1.1.0',
      dryRun: false,
      verbose: false,
    });
  });
});
