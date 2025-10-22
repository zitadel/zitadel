import { releaseChangelog, releasePublish, releaseVersion } from 'nx/release';
import yargs from 'yargs';

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
        'Whether or not to perform a dry-run of the release process, defaults to true',
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

  const { workspaceVersion, projectsVersionData } = await releaseVersion({
    specifier: options.version,
    dryRun: options.dryRun,
    verbose: options.verbose,
  });

  await releaseChangelog({
    versionData: projectsVersionData,
    version: workspaceVersion,
    dryRun: options.dryRun,
    verbose: options.verbose,
  });

  // publishResults contains a map of project names and their exit codes
  const publishResults = await releasePublish({
    dryRun: options.dryRun,
    verbose: options.verbose,
    // You can optionally pass through the version data (e.g. if you are using a custom publish executor that needs to be aware of versions)
    // It will then be provided to the publish executor options as `nxReleaseVersionData`
    // This is not required for the default @nx/js publish executor
    versionData: projectsVersionData,
  });

  process.exit(
    Object.values(publishResults).every((result) => result.code === 0) ? 0 : 1
  );
})();