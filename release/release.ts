import { execSync } from 'node:child_process';
import { releaseVersion, releaseChangelog, releasePublish } from 'nx/release';
import yargs from 'yargs';
import { config } from 'yargs';

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
// ZITADEL_RELEASE_PUSH is used in docker-bake-release.hcl to determine whether to push the docker images.
// It is false for dry runs and true otherwise.
const doPushEnvVar = "ZITADEL_RELEASE_PUSH";
// ZITADEL_RELEASE_GITHUB_ORG is used to specify the GitHub organization for which Docker images should be created.
const githubOrgEnvVar = "ZITADEL_RELEASE_GITHUB_ORG";

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

export interface EnvironmentConfig {
  versionEnvVar: string;
  revisionEnvVar: string;
  isLatestEnvVar: string;
  doPushEnvVar: string;
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
        'Whether or not to perform a dry-run of the release process, defaults to false',
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
export function configureGithubRepo(options: ReleaseOptions): void {
  const repo = options.githubRepo;
  if (repo.trim() !== 'zitadel/zitadel') {
    // Verify neither GH_TOKEN nor GITHUB_TOKEN are set to avoid accidental releases to the main repository
    if (process.env['GH_TOKEN'] || process.env['GITHUB_TOKEN']) {
      throw new Error(`When specifying a custom GitHub repository (${repo}), neither GH_TOKEN nor GITHUB_TOKEN must be set to avoid accidental releases to the main repository.`);
    }
    // Verify the gh CLI is authenticated and defaults to the specified repository
    if (!execSync('gh auth status', { stdio: 'pipe' }).toString().includes('Logged in to GitHub.com as')) {
      throw new Error(`When specifying a custom GitHub repository (${repo}), the gh CLI must be authenticated to avoid accidental releases to the main repository.`);
    }
    if (execSync('gh repo view | head -1', { stdio: 'pipe' }).toString() !== `${repo.trim()}\n`) {
      throw new Error(`When specifying a custom GitHub repository (${repo}), the gh CLI must default to the specified repository to avoid accidental releases to the main repository. Run "gh repo set-default ${repo}" to set the default repository.`);
    }
  }
  process.env[githubOrgEnvVar] = repo;
  console.log(`Setting ${githubOrgEnvVar}=${process.env[githubOrgEnvVar]} for Docker image creation`);
}

/**
 * Sets up environment variables for the release process.
 */
export function setupEnvironmentVariables(
  config: EnvironmentConfig,
  gitInfo: GitInfo,
  options: ReleaseOptions,
  conventionalCommits: boolean,
  workspaceVersion?: string
): void {
  process.env[config.revisionEnvVar] = gitInfo.sha;
  console.log(`Setting ${config.revisionEnvVar}=${process.env[config.revisionEnvVar]} for docker image labels`);

  process.env[config.doPushEnvVar] = options.dryRun ? 'false' : 'true';
  console.log(`Setting ${config.doPushEnvVar}=${process.env[config.doPushEnvVar]} based on dryRun = ${options.dryRun}`);

  if (!conventionalCommits) {
    process.env[config.versionEnvVar] = gitInfo.sha;
    process.env[config.isLatestEnvVar] = 'false';
  } else if (workspaceVersion) {
    process.env[config.versionEnvVar] = `v${workspaceVersion}`;
    const workspaceVersionIsHigherThanBeforeOrEqual = workspaceVersion.localeCompare(gitInfo.highestVersionBefore, undefined, { numeric: true, sensitivity: 'base' }) >= 0;
    process.env[config.isLatestEnvVar] = workspaceVersionIsHigherThanBeforeOrEqual ? 'true' : 'false';
  }
}

/**
 * Executes docker build commands with the appropriate configuration.
 */
export function executeDockerBuild(conventionalCommits: boolean): void {
  const baseCommand = 'pnpm nx run-many --target build-docker';
  const debugTarget = conventionalCommits ? ' build-docker-debug' : '';
  const bakeFiles = ' --file release/docker-bake-release.hcl --file apps/api/docker-bake-release.hcl --file apps/login/docker-bake-release.hcl';
  
  execSync(`${baseCommand}${debugTarget}${bakeFiles}`, { stdio: 'inherit', env: process.env });
}

/**
 * Main execution logic for the release process.
 * Returns exit code instead of calling process.exit for testability.
 */
export async function executeRelease(
  gitInfo: GitInfo,
  options: ReleaseOptions,
  envConfig: EnvironmentConfig
): Promise<number> {

  configureGithubRepo(options);

  const conventionalCommits = shouldUseConventionalCommits(gitInfo.branch);
  console.log(`Determined conventional commits = ${conventionalCommits} based on git branch = ${gitInfo.branch}`);

  setupEnvironmentVariables(envConfig, gitInfo, options, conventionalCommits);

  if (!conventionalCommits) {
    console.log(`Skipping GitHub release creation based on conventionalCommits=${conventionalCommits}. Instead setting ${envConfig.versionEnvVar}=${process.env[envConfig.versionEnvVar]} ${envConfig.isLatestEnvVar}=${process.env[envConfig.isLatestEnvVar]} and running the build-docker targets with additional docker-bake-release.hcl files to push SHA tagged Docker images for production.\n`);
    executeDockerBuild(conventionalCommits);
    return 0;
  }

  const { workspaceVersion, projectsVersionData } = await releaseVersion({
    dryRun: options.dryRun,
    verbose: options.verbose,    
  });

  if (!workspaceVersion) {
    throw new Error('Could not determine workspace version. No relevant changes found in conventional commits.');
  }
  
  setupEnvironmentVariables(envConfig, gitInfo, options, conventionalCommits, workspaceVersion);
  console.log(`Setting ${envConfig.versionEnvVar}=${process.env[envConfig.versionEnvVar]}`);
  console.log(`Setting ${envConfig.isLatestEnvVar}=${process.env[envConfig.isLatestEnvVar]} because ${envConfig.versionEnvVar}=${process.env[envConfig.versionEnvVar]} is higher or equal to the previously highest regular semantic tag v${gitInfo.highestVersionBefore}. Running the build-docker and build-docker-debug targets with additional docker-bake-release.hcl files to push Docker images.\n`);
  executeDockerBuild(conventionalCommits);

  await releaseChangelog({
    versionData: projectsVersionData,
    version: workspaceVersion,
    dryRun: options.dryRun,
    verbose: options.verbose,
  });

  const publishResults = await releasePublish({
    dryRun: options.dryRun,
    verbose: options.verbose,    
  });

  return Object.values(publishResults).every((result) => result.code === 0) ? 0 : 1;
}

/**
 * Main entry point for the release script.
 */
export async function main(argv: string[] = process.argv.slice(2)): Promise<number> {
  const gitInfo = determineGitInfo();
  const options = await parseReleaseOptions(argv);

  const envConfig: EnvironmentConfig = {
    versionEnvVar,
    revisionEnvVar,
    isLatestEnvVar,
    doPushEnvVar,
  };

  return executeRelease(gitInfo, options, envConfig);
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