import { execSync, execFileSync } from 'node:child_process';
import { releaseVersion, releaseChangelog, releasePublish } from 'nx/release';
import yargs from 'yargs';
import { Octokit } from 'octokit';

// ZITADEL_RELEASE_VERSION defaults to git SHA if not a conventional commit release.
// It is needed in the following places:
// - to compile the version into the API binary (nx-release-publish target)
// - to upload GitHub release assets by referencing a release by its Git tag
// - to build and push the docker images with the version tag
// ZITADEL_RELEASE_IS_LATEST is used in docker-bake-release.hcl to determine whether to tag the docker images as latest
// ZITADEL_RELEASE_GITHUB_ORG is used to specify the GitHub organization for which Docker images should be created.
// ZITADEL_RELEASE_GITHUB_REPO is used to link npm packages to the repository they are published to.
// NX_DRY_RUN is used to determine whether to perform a dry-run of the release process
// If NX_DRY_RUN is true, the nx-release-publish targets don't try to upload assets to a GitHub release

export interface GitInfo {
  branch: string;
  sha: string;
}

export interface ReleaseOptions {
  dryRun: boolean;
  verbose: boolean;
  isLatest: boolean;
}

/**
 * Determines git information needed for the release process.
 * Extracted for testability - can be mocked in tests.
 */
export function determineGitInfo(): GitInfo {
  const branch = execSync('git rev-parse --abbrev-ref HEAD').toString().trim();
  const sha = execSync('git rev-parse HEAD').toString().trim();
  return { branch, sha };
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
    .option('isLatest', {
      alias: 'l',
      description:
        'Whether or not the release is the latest version, defaults to true',
      type: 'boolean',
      default: true,
    })
    .parseAsync();

  return {
    dryRun: result.dryRun,
    verbose: result.verbose,
    isLatest: result.isLatest,
  };
}

/**
 * Main execution logic for the release process.
 * Returns exit code instead of calling process.exit for testability.
 */
export async function main(argv: string[] = process.argv.slice(2)): Promise<number> {
  const options = await parseReleaseOptions(argv);
  const gitInfo = determineGitInfo();
  const { workspaceVersion } = await releaseVersion({
    dryRun: true,
    verbose: true,
    dockerVersionScheme: 'semantic',
  });
  if (!workspaceVersion) {
    console.error('Failed to determine workspace version.');
    return 1;
  }
  console.log(`Preparing release for version: v${workspaceVersion}`);

  // If a release for the determined version already exists, we fail the release process unless the flag --create-or-update is set.
  failIfReleaseAlreadyExists(options, workspaceVersion);

  // We make the version available for the 
  process.env['ZITADEL_RELEASE_VERSION'] = `v${workspaceVersion}`;
  await releaseVersion({
    dryRun: options.dryRun,
    verbose: true,
    dockerVersionScheme: 'semantic',
  });
  await releaseChangelog({
    version: workspaceVersion,
    dryRun: options.dryRun,
    verbose: options.verbose,
  });
  const publishResults = await releasePublish({
    dryRun: options.dryRun,
    verbose: options.verbose,
    tag: options.isLatest ? 'latest' : gitInfo.branch.replaceAll('.', '-'),
  });

  // Nx Release uses the hardcoded value 'legacy' for the make_latest mode to create GitHub releases.
  // The legacy mode compares past releases by semantic version, but only for a limited time frame.
  // This causes backport releases to steal the "latest" badge from the highest semantic version release.
  // With fixGitHubReleaseLatestBadge we ensure that the highest semantic version is always marked as latest.
  fixGitHubReleaseLatestBadge(options, workspaceVersion);

  // Nx Release does not currently support uploading assets to releases, so we do this with uploadGitHubReleaseAssets and Octokit.
  uploadGitHubReleaseAssets(options, workspaceVersion);

  // After the release is published, we push the already built Docker images with pushDockerImages.
  // Runs pnpm nx run-many --target push-docker
  pushDockerImages(options, workspaceVersion);

  // When we created a new latest release, we trigger version bumping workflows in other GitHub repositories with bumpBrewtabVersion and bumpHelmChartAppVersion.
  bumpBrewtabVersion(options, workspaceVersion);
  bumpHelmChartAppVersion(options, workspaceVersion);

  // To keep the git working directory clean, we reset any changed files after a successful release with resetChangedFiles.
  resetChangedFiles();
  const code = Object.values(publishResults).every((result) => result.code === 0) ? 0 : 1;
  if (code === 0) {
    console.log(`Release process completed successfully for version ${workspaceVersion}. Resetting changed files.`);
    resetChangedFiles();
  }
  return code;
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
function fixGitHubReleaseLatestBadge(options: ReleaseOptions, workspaceVersion: string, octokit: Octokit): void {
  const ghArgs = ['release', 'edit', `v${workspaceVersion}`, '--latest'];
  if (options.dryRun) {
    console.log(`[Dry Run] Would execute: gh ${ghArgs.map(a => JSON.stringify(a)).join(' ')}`);
    return;
  }
  try {
    execFileSync('gh', ghArgs, { stdio: 'inherit' });
  } catch (error) {
    console.error(`Failed to update latest release badge to v${workspaceVersion}:`, error);
  }
}