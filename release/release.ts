import { execSync, execFileSync } from 'node:child_process';
import { releaseVersion, releaseChangelog, releasePublish } from 'nx/release';
import yargs from 'yargs';

// ZITADEL_RELEASE_VERSION defaults to git SHA if not a conventional commit release.
// It is needed in the following places:
// - to compile the version into the API binary (nx-release-publish target)
// - to upload GitHub release assets by referencing a release by its Git tag
// - to build and push the docker images with the version tag
const versionEnvVar = "ZITADEL_RELEASE_VERSION";
// ZITADEL_RELEASE_REVISION is used in docker-bake-release.hcl to add the git revision as a label to the docker images
const revisionEnvVar = "ZITADEL_RELEASE_REVISION";
// ZITADEL_RELEASE_IS_LATEST is used in docker-bake-release.hcl to determine whether to tag the docker images as latest
const isLatestEnvVar = "ZITADEL_RELEASE_IS_LATEST";
// ZITADEL_RELEASE_GITHUB_ORG is used to specify the GitHub organization for which Docker images should be created.
const githubOrgEnvVar = "ZITADEL_RELEASE_GITHUB_ORG";
// ZITADEL_RELEASE_GITHUB_REPO is used to link npm packages to the repository they are published to.
const githubRepoEnvVar = "ZITADEL_RELEASE_GITHUB_REPO";
// NX_DRY_RUN is used to determine whether to perform a dry-run of the release process
// If NX_DRY_RUN is true, the nx-release-publish targets don't try to upload assets to a GitHub release
const dryRunEnvVar = "NX_DRY_RUN";

export interface GitInfo {
  branch: string;
  sha: string;
  highestVersionBefore: string;
}

export interface ReleaseOptions {
  dryRun: boolean;
  verbose: boolean;
  githubRepo: string;
}

/**
 * Determines git information needed for the release process.
 * Extracted for testability - can be mocked in tests.
 */
export function determineGitInfo(): GitInfo {
  const branch = execSync('git rev-parse --abbrev-ref HEAD').toString().trim();
  const sha = execSync('git rev-parse HEAD').toString().trim();
  // highestVersionBefore is the highest semantic version tag in the repository that follows the format v[0-9]*.[0-9]*.[0-9]*
  // By comparing it to the determined workspace version, we can decide whether to tag the docker images as latest
  // The filter "v[0-9]*.[0-9]*.[0-9]*" excludes pre-release and build metadata tags like v1.0.0-beta or v1.0.0+build.1
  // The --sort=-v:refname flag sorts the tags by version number in descending order
  // -v is needed to sort by version number instead of lexicographically
  // :refname is needed to sort by tag name instead of commit date
  // head -n 1 gets the first line of the output, which is the highest version tag
  const highestVersionBefore = execSync('git tag --list "v[0-9]*.[0-9]*.[0-9]*" --sort=-v:refname | head -n 1').toString().trim().replace(/^v/, '');

  return { branch, sha, highestVersionBefore };
}

/**
 * Determines whether to use conventional commits based on the git branch.
 * Returns true for maintenance branches (v[0-9].x or v[0-9].[0-9].x).
 */
export function shouldUseConventionalCommits(branch: string): boolean {
  return /^v[0-9]+\.(x|[0-9]+\.x)$/.test(branch);
}

/**
 * Parses command line arguments for the release process.
 */
export async function parseReleaseOptions(argv: string[]): Promise<ReleaseOptions> {
  const result = await yargs(argv)
    .option('dryRun', {
      alias: 'd',
      description:
        'Whether or not to perform a dry-run of the release process, defaults to true',
      type: 'boolean',
      default: true,
    })
    .option('verbose', {
      alias: 'v',
      description:
        'Whether or not to enable verbose logging, defaults to false',
      type: 'boolean',
      default: false,
    })
    .option('githubRepo', {
      alias: 'r',
      description:
        'The GitHub repository for which the release should be created, defaults to zitadel/zitadel',
      type: 'string',
      requiresArg: true,
    })
    .demandOption('githubRepo', 'GitHub repository is required')
    .parseAsync();

  return {
    dryRun: result.dryRun,
    verbose: result.verbose,
    githubRepo: result.githubRepo,
  };
}

// configureGithubRepo makes sure that we can release to a different GitHub repository than zitadel/zitadel for testing purposes.
export function configureGithubRepo(options: ReleaseOptions) {
  const repo = options.githubRepo;
  const org = repo.trim().split('/')[0];
  if (repo.trim() !== 'zitadel/zitadel') {
    if (org === 'zitadel') {
      throw new Error('GitHub organization must not be zitadel when releasing to a different repository than zitadel/zitadel.');
    }
    if (execSync('gh repo view --json isFork --jq .isFork', { stdio: 'pipe' }).toString().trim() !== 'true') {
      throw new Error(`GitHub repository ${repo} of the current directory must be a fork of zitadel/zitadel.`);
    }
  }
  process.env[githubOrgEnvVar] = org;
  console.log(`Setting ${githubOrgEnvVar}=${process.env[githubOrgEnvVar]} for Docker image creation`);
  process.env[githubRepoEnvVar] = repo;
  console.log(`Setting ${githubRepoEnvVar}=${process.env[githubRepoEnvVar]} for npm package publication`);
}

export function setupDefaultEnvironmentVariables(
  gitSha: string,
  dryRun: boolean
): void {
  process.env[dryRunEnvVar] = dryRun ? 'true' : 'false';
  console.log(`Setting default ${dryRunEnvVar}=${process.env[dryRunEnvVar]}`);
  process.env[revisionEnvVar] = gitSha;
  console.log(`Setting default ${revisionEnvVar}=${process.env[revisionEnvVar]}`);
  process.env[versionEnvVar] = gitSha;
  console.log(`Setting default ${versionEnvVar}=${process.env[versionEnvVar]}`);
  process.env[isLatestEnvVar] = 'false';
  console.log(`Setting default ${isLatestEnvVar}=${process.env[isLatestEnvVar]}`);
}

/**
 * Sets up environment variables for the release process.
 */
export function setupWorkspaceVersionEnvironmentVariables(
  options: ReleaseOptions,
  gitInfo: GitInfo,
  workspaceVersion?: string | null
): void {

  if (!workspaceVersion) {
    throw new Error('Could not determine workspace version. No relevant changes found in conventional commits.');
  }

  if (!process.env["GITHUB_TOKEN"]) {
    throw new Error('GITHUB_TOKEN env must be set with a classic PAT and scope write:packages to create a release.');
  }

  const versionMatch = workspaceVersion.match(/^(\d+)\.(\d+)\.(\d+)$/); // Ensure it's in semver format
  if (!versionMatch) {
    throw new Error(`Workspace version ${workspaceVersion} is not a valid semver (e.g., 1.2.3).`);
  }
  const [, major, minor] = versionMatch;

  const branchMatch = gitInfo.branch.match(/^v(\d+)\.(x|\d+)(?:\.x)?$/);
  if (!branchMatch) {
    throw new Error(`Branch ${gitInfo.branch} is not a valid maintenance branch (e.g., v1.x or v1.2.x).`);
  }
  const [, branchMajor, branchMinor] = branchMatch;

  if (parseInt(major, 10) !== parseInt(branchMajor, 10) || (branchMinor !== 'x' && parseInt(minor, 10) !== parseInt(branchMinor, 10))) {
    throw new Error(`Workspace version ${workspaceVersion} does not match the maintenance branch ${gitInfo.branch}.`);
  }

  process.env[versionEnvVar] = `v${workspaceVersion}`;
  console.log(`Overwriting ${versionEnvVar}=${process.env[versionEnvVar]} based on workspace version ${workspaceVersion}  according to conventional commits`);
  const workspaceVersionIsHigherThanBeforeOrEqual = workspaceVersion.localeCompare(gitInfo.highestVersionBefore, undefined, { numeric: true, sensitivity: 'base' }) >= 0;
  process.env[isLatestEnvVar] = workspaceVersionIsHigherThanBeforeOrEqual ? 'true' : 'false';
  console.log(`Overwriting ${isLatestEnvVar}=${process.env[isLatestEnvVar]} because ${versionEnvVar}=${process.env[versionEnvVar]} is ${workspaceVersionIsHigherThanBeforeOrEqual ? 'higher than or equal to' : 'lower than'} the previously highest regular semantic tag v${gitInfo.highestVersionBefore}`);

}

export function executeDockerBuild(conventionalCommits: boolean, dryRun: boolean): void {
  let apiBakeArgs = '--file apps/api/docker-bake-release.hcl';
  if (conventionalCommits) {
    console.log('Publishing Debug API image as conventional commits were detected.');
    apiBakeArgs += ' api-debug';
  }
  executeDockerBuildForProject(dryRun, '@zitadel/api', apiBakeArgs);
  executeDockerBuildForProject(dryRun, '@zitadel/login', '--file apps/login/docker-bake-release.hcl');
}

/**
 * Executes docker build commands with the appropriate configuration.
 */
export function executeDockerBuildForProject(dryRun: boolean, nxProject: string, projectBakeArgs: string): void {
  // Construct argument array for safer invocation.
  let bakeArgsArray = ['--file', 'release/docker-bake-release.hcl'];
  if (projectBakeArgs) {
    // Split projectBakeArgs string by spaces for now. To be robust, if projectBakeArgs can contain quoted args,
    // consider using a parser like shell-quote, but in this codebase, args are hardcoded and simple.
    bakeArgsArray = bakeArgsArray.concat(projectBakeArgs.split(' '));
  }
  if (!dryRun) {
    console.log('Docker images will be pushed to the registry after build, because dryRun is false.');
    bakeArgsArray.push('--push');
  }
  // Build full nx arguments array
  const nxArgs = [
    'nx',
    'run',
    `${nxProject}:build-docker`,
    '--no-tui',
    `--args=${bakeArgsArray.join(' ')}`
  ];
  console.log(`Executing docker build with command: pnpm ${nxArgs.map(a => `"${a}"`).join(' ')}`);
  execFileSync('pnpm', nxArgs, {
    stdio: 'inherit', env: process.env
  });
}

function resetChangedFiles(): void {
  try {
    execSync('git checkout .npmrc package.json packages/zitadel-client/package.json packages/zitadel-proto/package.json apps/login/package.json', { stdio: 'inherit' });
    console.log('Reset changed files to clean state.');
  } catch (error) {
    console.error('Failed to reset changed files:', error);
  }
}

/**
 * Fix GitHub's "latest" designation after Nx creates the release with the hardcoded make_latest=legacy value.
 * https://github.com/nrwl/nx/blob/master/packages/nx/src/command-line/release/utils/remote-release-clients/github.ts#L405-L406
 * GitHub's legacy mode after some timeframe uses creation date instead of the semantic version, so backport releases steal the "latest" badge
 * We need to explicitly set the highest semantic version release as latest
 */
function fixGitHubReleaseLatestBadge(options: ReleaseOptions, gitInfo: GitInfo): void {
  const highestVersionTag = `v${gitInfo.highestVersionBefore}`;
  const currentVersionTag = process.env[versionEnvVar];
  const shouldFixLatest = currentVersionTag !== highestVersionTag;
  if (shouldFixLatest) {
    console.log(`Correcting GitHub "latest" release badge: moving from ${currentVersionTag} to ${highestVersionTag}`);
    const ghArgs = ['release', 'edit', highestVersionTag, '--latest'];
    if (options.dryRun) {
      console.log(`[Dry Run] Would execute: gh ${ghArgs.map(a => JSON.stringify(a)).join(' ')}`);
      return;
    }
    try {
      execFileSync('gh', ghArgs, { stdio: 'inherit' });
    } catch (error) {
      console.error(`Failed to update latest release badge to ${highestVersionTag}:`, error);
    }
  }
}

/**
 * Main execution logic for the release process.
 * Returns exit code instead of calling process.exit for testability.
 */
export async function executeRelease(
  gitInfo: GitInfo,
  options: ReleaseOptions
): Promise<number> {

  configureGithubRepo(options);

  setupDefaultEnvironmentVariables(gitInfo.sha, options.dryRun);

  const conventionalCommits = shouldUseConventionalCommits(gitInfo.branch);
  console.log(`Determined conventional commits = ${conventionalCommits} based on git branch = ${gitInfo.branch}`);

  if (!conventionalCommits) {
    console.log(`Skipping GitHub release creation based on conventionalCommits=${conventionalCommits}. Instead setting ${versionEnvVar}=${process.env[versionEnvVar]} ${isLatestEnvVar}=${process.env[isLatestEnvVar]} and running the build-docker targets with additional docker-bake-release.hcl files to push SHA tagged Docker images for production.`);
    executeDockerBuild(false, options.dryRun);
    return 0;
  }

  const { workspaceVersion, projectsVersionData } = await releaseVersion({
    dryRun: options.dryRun,
    verbose: options.verbose,
  });

  setupWorkspaceVersionEnvironmentVariables(options, gitInfo, workspaceVersion);
  console.log(`Setting ${versionEnvVar}=${process.env[versionEnvVar]}`);
  console.log(`Setting ${isLatestEnvVar}=${process.env[isLatestEnvVar]} because ${versionEnvVar}=${process.env[versionEnvVar]} is higher or equal to the previously highest regular semantic tag v${gitInfo.highestVersionBefore}. Running the build-docker and build-docker-debug targets with additional docker-bake-release.hcl files to push Docker images.`);
  executeDockerBuild(true, options.dryRun);

  await releaseChangelog({
    versionData: projectsVersionData,
    version: workspaceVersion,
    dryRun: options.dryRun,
    verbose: options.verbose,
  });

  fixGitHubReleaseLatestBadge(options, gitInfo);

  const publishResults = await releasePublish({
    dryRun: options.dryRun,
    verbose: options.verbose,
    tag: process.env[isLatestEnvVar] == 'true' ? 'latest' : gitInfo.branch.replaceAll('.', '-'),
  });

  const code = Object.values(publishResults).every((result) => result.code === 0) ? 0 : 1;
  if (code === 0) {
    console.log(`Release process completed successfully for version ${workspaceVersion}. Resetting changed files.`);
    resetChangedFiles();
  }
  return code;
}

/**
 * Main entry point for the release script.
 */
export async function main(argv: string[] = process.argv.slice(2)): Promise<number> {
  const gitInfo = determineGitInfo();
  const options = await parseReleaseOptions(argv);
  return executeRelease(gitInfo, options);
}

// Execute main when run directly
if (import.meta.url === `file://${process.argv[1]}`) {
  main().then((exitCode) => {
    process.exit(exitCode);
  }).catch((error) => {
    console.error('Release failed:', error);
    process.exit(1);
  });
}