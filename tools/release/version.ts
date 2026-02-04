import { releaseVersion } from 'nx/release/index.js';
import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';
import { execSync } from 'child_process';
import * as fs from 'fs';
import * as path from 'path';
import { DefaultArtifactClient } from '@actions/artifact';

(async () => {
    const argv = await yargs(hideBin(process.argv))
        .version(false) // disable default --version
        .option('dryRun', {
            alias: 'd',
            type: 'boolean',
            description: 'Whether or not to perform a dry-run of the release process, defaults to true',
            default: true,
        })
        .option('verbose', {
            description: 'Whether or not to enable verbose logging, defaults to false',
            type: 'boolean',
            default: false,
        })
        .parseAsync();

    const isMain = process.env.GITHUB_REF === 'refs/heads/main';
    const isPR = process.env.GITHUB_EVENT_NAME === 'pull_request';
    const dryRun = argv.dryRun;

    try {
        if (isPR || !isMain) {
            // Preview Version Logic
            const { workspaceVersion } = await releaseVersion({
                dryRun: true, // Always true for previews
                verbose: argv.verbose,
                firstRelease: true,
            });



            // Get branch name (try env first for CI, then git)
            let branch = process.env.GITHUB_HEAD_REF || process.env.GITHUB_REF_NAME;
            if (!branch) {
                try {
                    branch = execSync('git rev-parse --abbrev-ref HEAD').toString().trim();
                } catch (e) {
                    branch = 'unknown';
                }
            }
            // Sanitize branch name (replace non-alphanumeric-dash with dash)
            const sanitizedBranch = branch.replace(/[^a-zA-Z0-9-]/g, '-');

            // Use stable branch name as suffix (no SHA) to optimize caching
            const previewSuffix = sanitizedBranch;

            const previewVersion = `${workspaceVersion}+${previewSuffix}`;
            console.log(`Preview Version: ${previewVersion}`);

            // Output to env var for other steps
            console.log(`ZITADEL_VERSION=${previewVersion}`);
            // Write to file
            const artifactsDir = '.artifacts';
            if (!fs.existsSync(artifactsDir)) fs.mkdirSync(artifactsDir);
            fs.writeFileSync(path.join(artifactsDir, 'version'), previewVersion);

            if (process.env.GITHUB_ENV) {
                fs.appendFileSync(process.env.GITHUB_ENV, `ZITADEL_VERSION=${previewVersion}\n`);
            }

        } else {
            // Main Release Logic
            const { workspaceVersion } = await releaseVersion({
                dryRun: dryRun,
                verbose: argv.verbose,
            });
            console.log(`Release Version: ${workspaceVersion}`);

            if (!workspaceVersion) {
                console.error('Failed to determine workspace version.');
                process.exit(1);
            }

            console.log(`ZITADEL_VERSION=${workspaceVersion}`);
            // Write to file
            const artifactsDir = '.artifacts';
            if (!fs.existsSync(artifactsDir)) fs.mkdirSync(artifactsDir);
            fs.writeFileSync(path.join(artifactsDir, 'version'), workspaceVersion);

            if (process.env.GITHUB_ENV) {
                fs.appendFileSync(process.env.GITHUB_ENV, `ZITADEL_VERSION=${workspaceVersion}\n`);
            }
        }

        // Upload version as artifact if in CI
        if (process.env.GITHUB_ACTIONS) {
            const versionToUpload = path.join('.artifacts', 'version');
            if (fs.existsSync(versionToUpload)) {
                console.log('Uploading version artifact to GitHub Actions...');
                try {
                    const artifactClient = new DefaultArtifactClient();
                    const artifactName = `release-version`;
                    await artifactClient.uploadArtifact(
                        artifactName,
                        [versionToUpload],
                        process.cwd()
                    );
                    console.log(`Uploaded version artifact.`);
                } catch (err: any) {
                    if (err.message?.includes('ACTIONS_RUNTIME_TOKEN')) {
                        console.warn('Skipping version artifact upload: ACTIONS_RUNTIME_TOKEN not available.');
                    } else {
                        console.error('Failed to upload version artifact:', err);
                    }
                }
            }
        }

        process.exit(0);

    } catch (err) {
        console.error(err);
        process.exit(1);
    }
})();
