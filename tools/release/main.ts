import { releaseVersion } from 'nx/release/index.js';
import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';
import { execSync } from 'child_process';
import * as fs from 'fs';
import * as path from 'path';
import { Octokit } from 'octokit';

// Helper to sanitize branch names for versions
function getSanitizedBranch() {
    let branch = process.env.GITHUB_HEAD_REF || process.env.GITHUB_REF_NAME;
    if (!branch) {
        try {
            branch = execSync('git rev-parse --abbrev-ref HEAD').toString().trim();
        } catch (e) {
            branch = 'unknown';
        }
    }
    return branch.replace(/[^a-zA-Z0-9-]/g, '-');
}

// Helper to write version artifact
function writeVersionArtifact(version: string) {
    const artifactsDir = '.artifacts';
    if (!fs.existsSync(artifactsDir)) fs.mkdirSync(artifactsDir);
    fs.writeFileSync(path.join(artifactsDir, 'version'), version);
    if (process.env.GITHUB_ENV) {
        fs.appendFileSync(process.env.GITHUB_ENV, `ZITADEL_VERSION=${version}\n`);
    }
}

// Subcommand: VERSION
async function cmdVersion(argv: any) {
    const isMain = process.env.GITHUB_REF === 'refs/heads/main';
    const isPR = process.env.GITHUB_EVENT_NAME === 'pull_request';

    if (isPR || !isMain) {
        // Preview/Dev Logic
        const { workspaceVersion } = await releaseVersion({
            dryRun: true,
            verbose: argv.verbose,
            firstRelease: true,
        });

        const previewSuffix = getSanitizedBranch();
        const previewVersion = `${workspaceVersion}+${previewSuffix}`;
        console.log(`Preview Version: ${previewVersion}`);
        writeVersionArtifact(previewVersion);
    } else {
        // Main Logic
        const { workspaceVersion } = await releaseVersion({
            dryRun: argv.dryRun !== false, // Default to true unless explicitly false (though usually controlled by release cmd)
            verbose: argv.verbose,
        });

        // If we are just calculating version for a real release, dryRun might be false. 
        // But for 'version' command standalone, we usually just want to know what it IS.
        // However, standard use is: 'version' target runs first.

        console.log(`Release Version: ${workspaceVersion}`);
        if (workspaceVersion) {
            writeVersionArtifact(workspaceVersion);
        } else {
            console.error('Failed to determine workspace version.');
            process.exit(1);
        }
    }
}

// Subcommand: RELEASE
async function cmdRelease(argv: any) {
    // Safety: Default to Dry-Run unless CI_RELEASE=true
    const isLive = process.env.CI_RELEASE === 'true';
    const dryRun = !isLive;

    console.log(`Release Mode: ${isLive ? 'LIVE 🚀' : 'PLAN (Dry-Run) 🧪'}`);

    // 1. Calculate Version & Changelog
    // we use nx release to handle tagging and changelogs
    let version: string | undefined;

    try {
        const result = await releaseVersion({
            dryRun: dryRun,
            verbose: argv.verbose,
            gitCommit: true,
            gitTag: true,
        });
        version = result.workspaceVersion;

        // If dry run, we might get a preview version if not on main, or a calculated next version
        console.log(`Target Version: ${version}`);

    } catch (e) {
        console.error('Error during workspace version calculation:', e);
        process.exit(1);
    }

    if (!version) {
        // Fallback or read from artifact if nx release didn't return (e.g. no changes)
        // But for a release command we expect action.
        console.warn('No version changes detected by NX Release.');
        // Try reading artifact if exists
        try {
            version = fs.readFileSync(path.join(process.cwd(), '.artifacts/version'), 'utf-8').trim();
        } catch (e) { }
    }

    if (!version) {
        console.error('Could not determine version for release.');
        process.exit(1);
    }

    // 2. PR Comment (if PR and DryRun)
    const isPR = process.env.GITHUB_EVENT_NAME === 'pull_request';
    const token = process.env.GITHUB_TOKEN;
    const octokit = token ? new Octokit({ auth: token }) : null;
    const owner = 'zitadel';
    const repo = 'zitadel';

    if (isPR && dryRun && octokit) {
        const prNumber = process.env.GITHUB_REF?.split('/')[2];
        if (prNumber && !isNaN(parseInt(prNumber))) {
            const body = `### 🚀 Release Preview
**Version**: \`${version}\`
**Mode**: Plan (Dry-Run)
**Artifacts**: Ready for build.
`;
            console.log(`[Plan] Would comment on PR #${prNumber}: ${body}`);
            // We could uncomment this to actually post the preview comment even in dry-run
            try {
                await octokit.request('POST /repos/{owner}/{repo}/issues/{issue_number}/comments', {
                    owner, repo, issue_number: parseInt(prNumber), body,
                });
                console.log(`Commented on PR #${prNumber}`);
            } catch (e) {
                console.error('Failed to comment on PR:', e);
            }
        }
    }

    // 3. GitHub Release & Artifact Uploads (Live Mode Only)
    // In Live Mode, nx releaseVersion above should have created the Git Tag.
    // Now we need to create the GitHub Release and upload assets.
    // NOTE: NX Release can create GH Releases, but we want custom asset logic.

    if (isLive && octokit) {
        console.log('Creating/Updating GitHub Release...');
        const tag = `v${version}`;

        // Ensure assets exist
        const artifactsDir = path.join(process.cwd(), '.artifacts/pack');
        const files = fs.readdirSync(artifactsDir).filter(f => f.endsWith('.tar.gz') || f.endsWith('.zip') || f === 'checksums.txt');

        // ... (Logic to create release if not exists, or get it) ...
        // Simplification: We assume nx release might have created it if configured, OR we create it now.
        // Let's explicitly create/get it.

        try {
            // Create or Get Release
            let release;
            try {
                release = (await octokit.request('GET /repos/{owner}/{repo}/releases/tags/{tag}', { owner, repo, tag })).data;
            } catch (e) {
                console.log('Release not found, creating...');
                release = (await octokit.request('POST /repos/{owner}/{repo}/releases', {
                    owner, repo, tag_name: tag, name: tag, draft: false, prerelease: false, generate_release_notes: true
                })).data;
            }

            // Upload Assets
            for (const file of files) {
                console.log(`Uploading ${file}...`);
                const filePath = path.join(artifactsDir, file);
                const data = fs.readFileSync(filePath);
                await octokit.request('POST /repos/{owner}/{repo}/releases/{release_id}/assets{?name,label}', {
                    owner, repo, release_id: release.id, name: file, label: file,
                    headers: { 'content-type': 'application/octet-stream', 'content-length': data.length },
                    data: data as any
                });
            }
        } catch (e) {
            console.error('Failed GitHub Release operations:', e);
            process.exit(1);
        }
    } else if (dryRun) {
        console.log('[Plan] Would create GitHub Release and upload assets.');
    }

    // 4. Docker Push (Live Mode Only)
    // Targets renamed to publish-container
    const dockerTargets = [
        '@zitadel/api:publish-container',
        '@zitadel/login:publish-container'
    ];

    console.log('Processing Container Images...');
    for (const target of dockerTargets) {
        if (isLive) {
            console.log(`Running target: ${target}`);
            try {
                execSync(`npx nx run ${target}`, { stdio: 'inherit', env: process.env });
            } catch (e) {
                console.error(`Failed to run target ${target}:`, e);
                process.exit(1);
            }
        } else {
            console.log(`[Plan] Would run target: ${target}`);
        }
    }
}

// MAIN ENTRY POINT
(async () => {
    await yargs(hideBin(process.argv))
        .scriptName('release-tool')
        .command('version', 'Calculate and output version', {}, cmdVersion)
        .command('release', 'Execute release process (Plan or Live)', {}, cmdRelease)
        .option('verbose', { type: 'boolean', default: false })
        .demandCommand()
        .help()
        .parseAsync();
})();
