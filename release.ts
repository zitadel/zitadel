import { releaseChangelog, releasePublish, releaseVersion } from 'nx/release';
import { writeFileSync, mkdirSync } from 'fs';
import { dirname } from 'path';
import { execSync } from 'child_process';
import yargs from 'yargs';

const versionFilePath = "./.artifacts/next-version.txt";

(async () => {
  const options = await yargs(process.argv.slice(2))
    .version(false) // don't use the default meaning of version in yargs
    .option('version', {
      description:
        'Explicit version specifier to use, if overriding conventional commits',
      type: 'string',
    })
    .option('dryRun', {
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
    specifier: options.version,
    dryRun: options.dryRun,
    verbose: options.verbose,
  });

  // If no version change is needed, exit early
  if (!workspaceVersion) {
    console.log('\nNo version changes detected. Skipping changelog and publish steps.\n');
    process.exit(0);
  }

  // Restore git tracked changes
  // This prevents accidental commits of version changes
  // Unfortunately, Nx can't be configured to avoid modifying files during versioning
  if (!options.dryRun) {
    console.log('\nRestoring all git tracked changes...\n');
    execSync('git restore .', { stdio: 'inherit' });
  }
  
  // Write the version to a file so the nx-release-publish targets can read it
  // This is needed to compile the version into the release artifacts, like the API
  mkdirSync(dirname(versionFilePath), { recursive: true });
  writeFileSync(versionFilePath, workspaceVersion, 'utf-8');
  console.log(`\nWrote version ${workspaceVersion} to ${versionFilePath}\n`);

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
