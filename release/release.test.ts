/**
 * Unit tests for the release script.
 *
 * This test suite covers the testable functions extracted from release.ts including:
 * - Branch pattern validation for conventional commits
 * - Environment variable setup based on different release scenarios
 * - Docker build command construction
 * - Main release orchestration logic
 */

import { execFileSync, execSync } from 'node:child_process';
import { releaseVersion, releaseChangelog, releasePublish } from 'nx/release';
import { afterEach, beforeEach, describe, expect, test, vi } from 'vitest';
import {
  shouldUseConventionalCommits,
  setupWorkspaceVersionEnvironmentVariables,
  executeDockerBuild,
  parseReleaseOptions,
  executeRelease,
  type GitInfo,
  type ReleaseOptions,
  configureGithubRepo,
} from './release';

// Mock external dependencies
vi.mock('node:child_process', () => ({
  execSync: vi.fn(),
  execFileSync: vi.fn(),
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

  test('returns false for non-maintenance branches', () => {
    expect(shouldUseConventionalCommits('main')).toBe(false);
    expect(shouldUseConventionalCommits('next')).toBe(false);
    expect(shouldUseConventionalCommits('feature/my-feature')).toBe(false);
    expect(shouldUseConventionalCommits('fix/bug-123')).toBe(false);
    expect(shouldUseConventionalCommits('v1.0.0')).toBe(false);
  });
});

describe('setupWorkspaceVersionEnvironmentVariables', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    process.env = { ...originalEnv, GITHUB_TOKEN: 'classic pat with package:write' };
    // Clear relevant env vars
    delete process.env.ZITADEL_RELEASE_VERSION;
    delete process.env.ZITADEL_RELEASE_REVISION;
    delete process.env.ZITADEL_RELEASE_IS_LATEST;
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  const mockGitInfo: GitInfo = {
    branch: 'v2.x',
    sha: 'abc123def456',
    highestVersionBefore: '1.0.0',
  };

  const mockOptions: ReleaseOptions = {
    dryRun: false,
    verbose: false,
    githubRepo: 'zitadel/zitadel',
  };

  test('sets version with v prefix', () => {
    setupWorkspaceVersionEnvironmentVariables(mockOptions, mockGitInfo, '2.0.0');

    expect(process.env.ZITADEL_RELEASE_VERSION).toBe('v2.0.0');
  });

  test('sets isLatest to true when new version is higher than before', () => {
    setupWorkspaceVersionEnvironmentVariables(
      mockOptions,
      { ...mockGitInfo, highestVersionBefore: '1.0.0' },
      '2.0.0'
    );

    expect(process.env.ZITADEL_RELEASE_IS_LATEST).toBe('true');
  });

  test('sets isLatest to true when new version equals highest before', () => {
    setupWorkspaceVersionEnvironmentVariables(
      mockOptions,
      { ...mockGitInfo, highestVersionBefore: '2.0.0' },
      '2.0.0'
    );

    expect(process.env.ZITADEL_RELEASE_IS_LATEST).toBe('true');
  });

  test('sets isLatest to false when new version is lower than before', () => {
    setupWorkspaceVersionEnvironmentVariables(
      mockOptions,
      { ...mockGitInfo, highestVersionBefore: '3.0.0' },
      '2.0.0'
    );

    expect(process.env.ZITADEL_RELEASE_IS_LATEST).toBe('false');
  });

  test('throws error if workspaceVersion is not provided', () => {
    expect(() =>
      setupWorkspaceVersionEnvironmentVariables(mockOptions, mockGitInfo)
    ).toThrowError();
  });

  test('sets env vars for matching minor maintenance branch and version', () => {
    setupWorkspaceVersionEnvironmentVariables(
      mockOptions,
      { ...mockGitInfo, branch: 'v2.5.x' },
      '2.5.1'
    );

    expect(process.env.ZITADEL_RELEASE_VERSION).toBe('v2.5.1');
  });

  test('sets env vars for matching major maintenance branch and version', () => {
    setupWorkspaceVersionEnvironmentVariables(
      mockOptions,
      { ...mockGitInfo, branch: 'v3.x' },
      '3.0.2'
    );

    expect(process.env.ZITADEL_RELEASE_VERSION).toBe('v3.0.2');
  });

  test('throws error if major workspaceVersion does not match major maintenance branch', () => {
    expect(() =>
      setupWorkspaceVersionEnvironmentVariables(
        mockOptions,
        { ...mockGitInfo, branch: 'v2.x' },
        '3.0.0'
      )
    ).toThrowError();
  });

  test('throws error if minor workspaceVersion does not match minor maintenance branch', () => {
    expect(() =>
      setupWorkspaceVersionEnvironmentVariables(
        mockOptions,
        { ...mockGitInfo, branch: 'v2.5.x' },
        '2.6.0'
      )
    ).toThrowError();
  });

  test('throws error if neither GITHUB_TOKEN nor GITHUB_API_TOKEN is set', () => {
    process.env = { ...originalEnv };
    expect(() =>
      setupWorkspaceVersionEnvironmentVariables(
        mockOptions,
        mockGitInfo,
        '2.0.0'
      )
    ).toThrowError();
  });
});

describe('executeDockerBuild', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test('includes the release bake file', () => {
    const mockExecFileSync = vi.mocked(execFileSync);

    executeDockerBuild(false, false);

    const calls = mockExecFileSync.mock.calls.map(call => call[1] as Array<string>);
    expect(calls.every(call => call.find(arg => arg.includes('release/docker-bake-release.hcl')))).toBeTruthy();
  });

  test('includes --push flag when dryRun is false', () => {
    const mockExecFileSync = vi.mocked(execFileSync);

    executeDockerBuild(false, false);

    const calls = mockExecFileSync.mock.calls.map(call => call[1] as Array<string>);
    expect(calls.every(args => args.find(arg => arg.includes('--push')))).toBeTruthy();
  });

  test('includes api-debug target for conventional commits', () => {
    const mockExecFileSync = vi.mocked(execFileSync);

    executeDockerBuild(true, false);

    const calls = mockExecFileSync.mock.calls.map(call => call[1] as Array<string>);
    expect(calls.find(args => args.find(arg => arg.includes('api-debug')))).toBeTruthy();
  });

  test('passes environment variables to execFileSync', () => {
    const mockExecFileSync = vi.mocked(execFileSync);

    executeDockerBuild(false, false);

    expect(mockExecFileSync).toHaveBeenCalledWith(
      expect.any(String),
      expect.any(Array<String>),
      expect.objectContaining({
        stdio: 'inherit',
        env: process.env,
      })
    );
  });
});

describe('parseReleaseOptions', () => {

  const githubRepo = ['--githubRepo', 'some/repo'];

  test('parses githubRepo option', async () => {
    const options = await parseReleaseOptions(githubRepo);
    expect(options.githubRepo).toBe('some/repo');
  });

  test('throws error if githubRepo is missing', async () => {
    await expect(parseReleaseOptions([])).rejects.toThrowError();
  });

  test('parses dryRun flag', async () => {
    const options = await parseReleaseOptions(['--dryRun', ...githubRepo]);
    expect(options.dryRun).toBe(true);
  });

  test('parses verbose flag', async () => {
    const options = await parseReleaseOptions(['--verbose', ...githubRepo]);
    expect(options.verbose).toBe(true);
  });

  test('parses both flags', async () => {
    const options = await parseReleaseOptions(['--dryRun', '--verbose', ...githubRepo]);
    expect(options.dryRun).toBe(true);
    expect(options.verbose).toBe(true);
  });

  test('defaults dryRun to true', async () => {
    const options = await parseReleaseOptions(githubRepo);
    expect(options.dryRun).toBe(true);
  });

  test('defaults verbose to false', async () => {
    const options = await parseReleaseOptions(githubRepo);
    expect(options.verbose).toBe(false);
  });

  test('accepts --no-dryRun to set dryRun to false', async () => {
    const options = await parseReleaseOptions(['--no-dryRun', ...githubRepo]);
    expect(options.dryRun).toBe(false);
  });

  test('accepts --no-dry-run to set dryRun to false', async () => {
    const options = await parseReleaseOptions(['--no-dry-run', ...githubRepo]);
    expect(options.dryRun).toBe(false);
  });
});

describe('configureGithubRepo', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    process.env = { ...originalEnv };
    delete process.env.ZITADEL_RELEASE_GITHUB_ORG;
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  test('sets ZITADEL_RELEASE_GITHUB_ORG to zitadel for zitadel/zitadel repo', () => {
    const options: ReleaseOptions = {
      dryRun: false,
      verbose: false,
      githubRepo: 'zitadel/zitadel',
    };

    const mockExecSync = vi.mocked(execSync);
    configureGithubRepo(options);
    expect(mockExecSync).not.toHaveBeenCalled();

    expect(process.env.ZITADEL_RELEASE_GITHUB_ORG).toBe('zitadel');
  });

  test('sets ZITADEL_RELEASE_GITHUB_ORG to custom org for other repos', () => {
    const mockExecSync = vi.mocked(execSync);
    mockExecSync.mockReturnValue('true\n');
    const options: ReleaseOptions = {
      dryRun: false,
      verbose: false,
      githubRepo: 'customorg/customrepo',
    };

    configureGithubRepo(options);
    expect(mockExecSync).toHaveBeenCalledOnce();
    expect(process.env.ZITADEL_RELEASE_GITHUB_ORG).toBe('customorg');
  });

  test('throws error if org is zitadel for non-zitadel repo', () => {
    const mockExecSync = vi.mocked(execSync);
    const options: ReleaseOptions = {
      dryRun: false,
      verbose: false,
      githubRepo: 'zitadel/customrepo',
    };
    expect(() => configureGithubRepo(options)).toThrowError();
    expect(mockExecSync).not.toHaveBeenCalled();
  });

  test('throws error if repo is not zitadel/zitadel and not a fork', () => {
    const mockExecSync = vi.mocked(execSync);
    mockExecSync.mockReturnValue('false\n');
    const options: ReleaseOptions = {
      dryRun: false,
      verbose: false,
      githubRepo: 'customorg/customrepo',
    };
    expect(() => configureGithubRepo(options)).toThrowError();
    expect(mockExecSync).toHaveBeenCalledOnce();
  });
});

describe('executeRelease', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    process.env = { ...originalEnv, GITHUB_TOKEN: 'classic pat with package:write' };
    delete process.env.ZITADEL_RELEASE_VERSION;
    delete process.env.ZITADEL_RELEASE_REVISION;
    delete process.env.ZITADEL_RELEASE_IS_LATEST;
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  const mockGitInfo: GitInfo = {
    branch: 'v1.x',
    sha: 'abc123def456',
    highestVersionBefore: '1.0.0',
  };

  const mockOptions: ReleaseOptions = {
    dryRun: false,
    verbose: false,
    githubRepo: 'zitadel/zitadel',
  };

  test('returns 0 for non-conventional commits after docker build', async () => {
    const mockGitInfo: GitInfo = {
      branch: 'main',
      sha: 'abc123def456',
      highestVersionBefore: '1.0.0',
    };

    const exitCode = await executeRelease(mockGitInfo, mockOptions);
    expect(exitCode).toBe(0);
  });

  test('calls releaseVersion for conventional commits', async () => {

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: {} as any,
    });

    vi.mocked(releasePublish).mockResolvedValue({});

    await executeRelease(mockGitInfo, mockOptions);

    expect(releaseVersion).toHaveBeenCalledWith({
      dryRun: false,
      verbose: false,
    });
  });

  test('throws error when workspace version cannot be determined', async () => {

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: undefined,
      projectsVersionData: {} as any,
    });

    await expect(
      executeRelease(mockGitInfo, mockOptions)
    ).rejects.toThrowError();
  });

  test('returns 0 when all publish results succeed', async () => {

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: {} as any,
    });

    vi.mocked(releasePublish).mockResolvedValue({
      project1: { code: 0 },
      project2: { code: 0 },
    } as any);

    const exitCode = await executeRelease(mockGitInfo, mockOptions);

    expect(exitCode).toBe(0);
  });

  test('returns 1 when any publish result fails', async () => {

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: {} as any,
    });

    vi.mocked(releasePublish).mockResolvedValue({
      project1: { code: 0 },
      project2: { code: 1 },
    } as any);

    const exitCode = await executeRelease(mockGitInfo, mockOptions);

    expect(exitCode).toBe(1);
  });

  test('calls releaseChangelog with correct parameters', async () => {

    const mockProjectsVersionData = { someData: true };

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: mockProjectsVersionData as any,
    });

    vi.mocked(releasePublish).mockResolvedValue({});

    await executeRelease(mockGitInfo, mockOptions);

    expect(releaseChangelog).toHaveBeenCalledWith({
      versionData: mockProjectsVersionData,
      version: '1.1.0',
      dryRun: false,
      verbose: false,
    });
  });

  test('does not call gh release edit when current version equals highest version', async () => {
    const mockExecSync = vi.mocked(execSync);

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.0.0',
      projectsVersionData: {} as any,
    });
    // Avoid triggering resetChangedFiles (which uses execSync) by simulating a failed publish
    vi.mocked(releasePublish).mockResolvedValue({
      project1: { code: 1 },
    } as any);

    await executeRelease(
      { ...mockGitInfo, highestVersionBefore: '1.0.0' },
      { ...mockOptions, dryRun: false }
    );

    // Ensure we did not attempt to edit the latest release
    expect(
      mockExecSync.mock.calls.some(([cmd]) =>
        String(cmd).includes('gh release edit')
      )
    ).toBe(false);
  });

  test('does not call gh release edit when dryRun is true even if tags differ', async () => {
    const mockExecSync = vi.mocked(execSync);

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: {} as any,
    });
    vi.mocked(releasePublish).mockResolvedValue({
      project1: { code: 1 },
    } as any);

    await executeRelease(
      { ...mockGitInfo, highestVersionBefore: '1.0.0' },
      { ...mockOptions, dryRun: true }
    );

    expect(
      mockExecSync.mock.calls.some(([cmd]) =>
        String(cmd).includes('gh release edit')
      )
    ).toBe(false);
  });

  test('calls gh release edit to set latest when current version tag differs and not dryRun', async () => {
    const mockExecSync = vi.mocked(execSync);

    vi.mocked(releaseVersion).mockResolvedValue({
      workspaceVersion: '1.1.0',
      projectsVersionData: {} as any,
    });
    vi.mocked(releasePublish).mockResolvedValue({
      project1: { code: 1 },
    } as any);

    await executeRelease(
      { ...mockGitInfo, highestVersionBefore: '2.0.0' },
      { ...mockOptions, dryRun: false }
    );

    expect(mockExecSync).toHaveBeenCalledWith(
      'gh release edit v2.0.0 --latest',
      expect.objectContaining({ stdio: 'inherit' })
    );
  });
});
