import { releaseVersion } from 'nx/release/index.js';
import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';
import { execSync } from 'child_process';
import * as fs from 'fs';
import * as path from 'path';

(async () => {
    const argv = await yargs(hideBin(process.argv))
        .version(false) // disable default --version
        .option('verbose', {
            description: 'Whether or not to enable verbose logging, defaults to false',
            type: 'boolean',
            default: false,
        })
        .parseAsync();

    // logic: considered main only if explicitly on main branch.
    // Local (no GITHUB_REF) will fall through to 'false' (Preview), satisfying "locally ... [preview] version".
    const isMain = process.env.GITHUB_REF === 'refs/heads/main';

    const isPR = process.env.GITHUB_EVENT_NAME === 'pull_request' || !isMain;

    // We strictly dry-run in PR/Local to avoid git tagging side effects.
    // On Main, we allow actual modifications (unless env var overrides, but for now we assume Main = Release).
    const dryRun = isPR;

    console.log(`Running version task. isMain: ${isMain}, isPR: ${isPR}, dryRun: ${dryRun}`);

    try {
        if (isPR) {
            // Preview Version Logic
            const { workspaceVersion } = await releaseVersion({
                dryRun: true, // Always valid for preview/calculation
                verbose: argv.verbose,
                firstRelease: true,
            });

            const sha = execSync('git rev-parse --short HEAD').toString().trim();
            const previewVersion = `${workspaceVersion}+${sha}`;
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
                dryRun: false, // On Main, we effectively release
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

        process.exit(0);

    } catch (err) {
        console.error(err);
        process.exit(1);
    }
})();
