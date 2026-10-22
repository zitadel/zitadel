// Load the additional workspace in memory. The administrative Angular workspace stays unchanged.
import { createRequire } from 'node:module';
import { readFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import { dirname, resolve } from 'node:path';
const require = createRequire(import.meta.url);
const cliRequire = createRequire(require.resolve('@angular/cli/package.json'));
const { workspaces, logging, schema } = cliRequire('@angular-devkit/core');
const { NodeJsSyncHost } = cliRequire('@angular-devkit/core/node');
const { Architect } = cliRequire('@angular-devkit/architect');
const { WorkspaceNodeModulesArchitectHost } = cliRequire('@angular-devkit/architect/node');
const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const host = workspaces.createWorkspaceHost(new NodeJsSyncHost());
const { workspace } = await workspaces.readWorkspace(root, host);
const config = JSON.parse(await readFile(resolve(root, 'account/angular.json'), 'utf8')).projects.account;
workspace.projects.add({
  name: 'account',
  root: '',
  sourceRoot: 'src/account',
  projectType: 'application',
  targets: config.architect,
});
const nativeTheme = await readFile(
  resolve(root, 'src/app/modules/policies/private-labeling-policy/private-labeling-policy.component.ts'),
  'utf8',
);
if (!/export enum Theme\s*\{\s*DARK,\s*LIGHT,?\s*\}/.test(nativeTheme))
  throw new Error('Native asset theme contract changed; review Account build boundary');
const registry = new schema.CoreSchemaRegistry();
registry.addPostTransform(schema.transforms.addUndefinedDefaults);
const architect = new Architect(new WorkspaceNodeModulesArchitectHost(workspace, root), registry);
const logger = new logging.Logger('account');
logger.subscribe((entry) => (entry.level === 'error' ? process.stderr : process.stdout).write(entry.message + '\n'));
const serve = process.argv.includes('--serve');
const configuration = process.argv.includes('development') ? 'development' : 'production';
const run = await architect.scheduleTarget(
  { project: 'account', target: serve ? 'serve' : 'build', configuration: serve ? 'development' : configuration },
  {},
  { logger },
);
process.on('SIGINT', () => {
  void run.stop().then(() => process.exit(130));
});
if (serve) {
  run.output.subscribe((result) => {
    if (result.baseUrl) logger.info(result.baseUrl);
  });
} else {
  const result = await run.result;
  await run.stop();
  process.exitCode = result.success ? 0 : 1;
}
