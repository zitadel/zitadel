import { execSync } from 'node:child_process';
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
// ZITADEL_RELEASE_PUSH is used in docker-bake-release.hcl to determine whether to push the docker images.
// It is false for dry runs and true otherwise.
const doPushEnvVar = "ZITADEL_RELEASE_PUSH";

const gitBranch = execSync('git rev-parse --abbrev-ref HEAD').toString().trim();
const gitSha = execSync('git rev-parse HEAD').toString().trim();
// highestVersionBefore is the highest semantic version tag in the repository that follows the format v[0-9]*.[0-9]*.[0-9]*
// By comparing it to the determined workspace version, we can decide whether to tag the docker images as latest
// The filter "v[0-9]*.[0-9]*.[0-9]*" excludes pre-release and build metadata tags like v1.0.0-beta or v1.0.0+build.1
// The --sort=-v:refname flag sorts the tags by version number in descending order
// -v is needed to sort by version number instead of lexicographically
// :refname is needed to sort by tag name instead of commit date
// head -n 1 gets the first line of the output, which is the highest version tag
const highestVersionBefore = execSync('git tag --list "v[0-9]*.[0-9]*.[0-9]*" --sort=-v:refname | head -n 1').toString().trim().replace(/^v/, '');


(async () => {

  const options = await yargs(process.argv.slice(2))
    .option('dryRun', {
      alias: 'd',
      description:
        'Whether or not to perform a dry-run of the release process, defaults to false',
      type: 'boolean',
      default: true,
    })
    .option('verbose', {
      description:
        'Whether or not to enable verbose logging, defaults to false',
      type: 'boolean',
      default: false,
    })
    .parseAsync();

  process.env[revisionEnvVar] = gitSha;
  console.log(`Setting ${revisionEnvVar}=${process.env[revisionEnvVar]} for docker image labels`);

  const conventionalCommits = /^v[0-9]+\.(x|[0-9]+\.x)$/.test(gitBranch);
  console.log(`Determined conventional commits = ${conventionalCommits} based on git branch = ${gitBranch}`);

  process.env[doPushEnvVar] = options.dryRun ? 'false' : 'true';
  console.log(`Setting ${doPushEnvVar}=${process.env[doPushEnvVar]} based on dryRun = ${options.dryRun}`);

  if (!conventionalCommits) {
    process.env[versionEnvVar] = gitSha;
    process.env[isLatestEnvVar] = 'false';
    console.log(`Skipping GitHub release creation based on conventionalCommits=${conventionalCommits}. Instead setting ${versionEnvVar}=${process.env[versionEnvVar]} ${isLatestEnvVar}=${process.env[isLatestEnvVar]} and running the build-docker targets with additional docker-bake-release.hcl files to push SHA tagged Docker images for production.\n`);
    execSync('pnpm nx run-many --target build-docker --file release/docker-bake-release.hcl --file apps/api/docker-bake-release.hcl --file apps/login/docker-bake-release.hcl', { stdio: 'inherit', env: process.env });
    process.exit(0);
  }

  const { workspaceVersion, projectsVersionData } = await releaseVersion({
    dryRun: options.dryRun,
    verbose: options.verbose,
  });

  if (!workspaceVersion) {
    throw new Error('Could not determine workspace version. No relevant changes found in conventional commits.');
  }

  console.log(`Creating GitHub release for version v${workspaceVersion} based on conventional commits on branch ${gitBranch}`);
  
  process.env[versionEnvVar] = `v${workspaceVersion}`;
  console.log(`Setting ${versionEnvVar}=${process.env[versionEnvVar]}`);
  const workspaceVersionIsHigherThanBeforeOrEqual = highestVersionBefore.localeCompare(workspaceVersion, undefined, { numeric: true, sensitivity: 'base' }) >= 0;
  process.env[isLatestEnvVar] = workspaceVersionIsHigherThanBeforeOrEqual ? 'true' : 'false';
  console.log(`Setting ${isLatestEnvVar}=${process.env[isLatestEnvVar]} because ${versionEnvVar}=${process.env[versionEnvVar]} is higher or equal to the previously highest regular semantic tag v${highestVersionBefore}`);

  await releaseChangelog({
    versionData: projectsVersionData,
    version: workspaceVersion,
    dryRun: options.dryRun,
    verbose: options.verbose
  });

  const publishResults = await releasePublish({
    dryRun: options.dryRun,
    verbose: options.verbose,    
  });

  process.exit(
    Object.values(publishResults).every((result) => result.code === 0) ? 0 : 1
  )
})();
