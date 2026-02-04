import { Octokit } from 'octokit';
import { DefaultArtifactClient } from '@actions/artifact';
import * as fs from 'fs';
import * as path from 'path';
import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';

(async () => {
    const argv = await yargs(hideBin(process.argv))
        .option('dryRun', {
            alias: 'd',
            type: 'boolean',
            default: true,
        })
        .option('verbose', {
            type: 'boolean',
            default: false,
        })
        .parseAsync();

    let version = process.env.ZITADEL_VERSION;

    if (!version) {
        try {
            version = fs.readFileSync(path.join(process.cwd(), '.artifacts/version'), 'utf-8').trim();
        } catch (e) {
            // ignore
        }
    }

    if (!version) {
        console.error('ZITADEL_VERSION not set and .artifacts/version not found.');
        process.exit(1);
    }

    console.log(`Publishing version: ${version}`);
    const dryRun = argv.dryRun;

    // Context detection for routing actions (what to do), not for safety (whether to do it)
    const isMain = process.env.GITHUB_REF === 'refs/heads/main';
    const isPR = process.env.GITHUB_EVENT_NAME === 'pull_request';

    console.log(`Context - isMain: ${isMain}, isPR: ${isPR}, dryRun: ${dryRun}`);

    if (!process.env.GITHUB_TOKEN) {
        console.warn('WARNING: GITHUB_TOKEN is not set. GitHub interactions (release/comment) will be skipped.');
    }

    const artifactsDir = path.join(process.cwd(), '.artifacts/pack');
    if (!fs.existsSync(artifactsDir)) {
        console.error(`Artifacts directory not found: ${artifactsDir}`);
        process.exit(1);
    }

    const files = fs.readdirSync(artifactsDir).filter(f => f.endsWith('.tar.gz') || f.endsWith('.zip') || f === 'checksums.txt');
    const filePaths = files.map(f => path.join(artifactsDir, f));

    // 1. Upload to GitHub Actions Artifacts (Now handled by ci.yml)
    if (process.env.GITHUB_ACTIONS) {
        console.log('Artifacts generated in .artifacts/ directory. CI workflow will handle upload.');
    }

    // Octokit setup
    const token = process.env.GITHUB_TOKEN;
    if (!token) {
        if (!dryRun && isMain) {
            console.error('GITHUB_TOKEN required for publishing to Main.');
            process.exit(1);
        } else {
            console.warn('GITHUB_TOKEN not set. GitHub interactions will be skipped.');
        }
    }
    const octokit = token ? new Octokit({ auth: token }) : null;
    const owner = 'zitadel';
    const repo = 'zitadel';

    // 2. PR Logic: Comment on PR
    if (isPR && octokit) {
        console.log('Posting comment to PR...');
        // Need PR number. GITHUB_REF for PR is refs/pull/:prNumber/merge
        const prNumber = process.env.GITHUB_REF?.split('/')[2];
        if (prNumber && !isNaN(parseInt(prNumber))) {
            const body = `### 🚀 Release Preview
**Version**: \`${version}\`
**Artifacts**: Uploaded to GitHub Actions Summary.
`;
            if (!dryRun) {
                try {
                    await octokit.request('POST /repos/{owner}/{repo}/issues/{issue_number}/comments', {
                        owner,
                        repo,
                        issue_number: parseInt(prNumber),
                        body,
                    });
                    console.log(`Commented on PR #${prNumber}`);
                } catch (e) {
                    console.error('Failed to comment on PR:', e);
                }
            } else {
                console.log(`[Dry-Run] Would comment on PR #${prNumber}: ${body}`);
            }
        }
    }

    // 3. Main Logic: Upload to GitHub Release
    if (isMain && octokit) {
        console.log('Uploading assets to GitHub Release...');
        const tag = `v${version}`; // Assuming convention
        // Release should have been created by nx release
        try {
            let release;
            try {
                const releaseResponse = await octokit.request('GET /repos/{owner}/{repo}/releases/tags/{tag}', {
                    owner,
                    repo,
                    tag,
                });
                release = releaseResponse.data;
            } catch (e) {
                console.log(`Release ${tag} not found.`);
            }

            if (release) {
                for (const file of files) {
                    const filePath = path.join(artifactsDir, file);
                    const data = fs.readFileSync(filePath);
                    if (!dryRun) {
                        // Check if asset exists? overwrite?
                        // For now just upload
                        console.log(`Uploading ${file}...`);
                        await octokit.request('POST /repos/{owner}/{repo}/releases/{release_id}/assets{?name,label}', {
                            owner,
                            repo,
                            release_id: release.id,
                            name: file,
                            label: file,
                            headers: {
                                'content-type': 'application/octet-stream',
                                'content-length': data.length,
                            },
                            data: data as any,
                        });
                    } else {
                        console.log(`[Dry-Run] Would upload ${file} to release ${release.id}`);
                    }
                }
            } else {
                console.warn(`Release ${tag} not found (might be created later or dry-run). Cannot upload assets.`);
            }

        } catch (e) {
            console.error('Error handling GitHub Release:', e);
            process.exit(1);
        }
    }

    // 4. Docker (Excluded for now)

})();
