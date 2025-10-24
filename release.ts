import { releaseChangelog, releasePublish, releaseVersion } from 'nx/release';
import { writeFileSync, mkdirSync } from 'fs';
import { dirname } from 'path';
import { execSync } from 'child_process';
import yargs from 'yargs';

const versionEnvVar = "ZITADEL_VERSION";

(async () => {
  const options = await yargs(process.argv.slice(2))
    .version(false) // don't use the default meaning of version in yargs
    .option('version', {
      description:
        'Explicit version specifier to use, if overriding conventional commits',
      type: 'string',
    })
    .option('dry-run', {
      alias: 'd',
      description:
        'Whether or not to perform a dry-run of the release process, defaults to false',
      type: 'boolean',
      default: false,
    })
    .option('verbose', {
      description:
        'Whether or not to enable verbose logging, defaults to false',
      type: 'boolean',
      default: false,
    })
    .parseAsync();

  const { workspaceVersion, projectsVersionData } = await releaseVersion({
    specifier: options.version || process.env.ZITADEL_VERSION,
    dryRun: options.dryRun,
    verbose: options.verbose,
  });

  // If no version change is needed, exit early
  if (!workspaceVersion) {
    console.log('\nNo version changes detected. Skipping changelog and publish steps.\n');
    process.exit(0);
  }

  // Write the version to an env variable so subsequent steps can read it.
  // This is needed in the following places:
  // - to compile the version into the API binary (nx-release-publish target)
  // - to resolve the docker versionScheme v{env.ZITADEL_VERSION}
  // - to upload GitHub release assets by referencing a release by its tag v{env.ZITADEL_VERSION}
  const zitadelVersion = options.version || process.env.ZITADEL_VERSION || `v${workspaceVersion}`;
  console.log(`\nUsing environment variable ${versionEnvVar}=${zitadelVersion}\n`);
  process.env[versionEnvVar] = zitadelVersion;

  await releaseChangelog({
    versionData: projectsVersionData,
    version: workspaceVersion,
    dryRun: options.dryRun,
    verbose: options.verbose
  });

  // publishResults contains a map of project names and their exit codes
  const publishResults = await releasePublish({
    dryRun: options.dryRun,
    verbose: options.verbose
  });

  process.exit(
    Object.values(publishResults).every((result) => result.code === 0) ? 0 : 1
  );
})();
