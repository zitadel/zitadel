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
            default: false,
        })
        .option('verbose', {
            type: 'boolean',
            default: false,
        })
        .parseAsync();

    const version = process.env.ZITADEL_VERSION;
    if (!version) {
        console.error('ZITADEL_VERSION not set.');
        process.exit(1);
    }

    console.log(`Publishing version: ${version}`);
    const dryRun = argv.dryRun;

    // Logic for Main vs PR detection
    const isMain = process.env.GITHUB_REF === 'refs/heads/main';
    const isPR = process.env.GITHUB_EVENT_NAME === 'pull_request' || !isMain;

    console.log(`Context - isMain: ${isMain}, isPR: ${isPR}, dryRun: ${dryRun}`);

    const artifactsDir = path.join(process.cwd(), '.artifacts/pack');
    if (!fs.existsSync(artifactsDir)) {
        console.error(`Artifacts directory not found: ${artifactsDir}`);
        process.exit(1);
    }

    const files = fs.readdirSync(artifactsDir).filter(f => f.endsWith('.tar.gz') || f.endsWith('.zip') || f === 'checksums.txt');
    const filePaths = files.map(f => path.join(artifactsDir, f));

    // 1. Upload to GitHub Actions Artifacts (Available in CI)
    if (process.env.GITHUB_ACTIONS) {
        console.log('Uploading to GitHub Actions Artifacts...');
        const artifactClient = new DefaultArtifactClient();
        const artifactName = `release-artifacts-${version}`;
        try {
            const { id, size } = await artifactClient.uploadArtifact(
                artifactName,
                filePaths,
                artifactsDir
            );
            console.log(`Uploaded artifact ${id} (${size} bytes)`);
        } catch (err) {
            console.error('Failed to upload artifact to GitHub Actions:', err);
        }
    } else {
        console.log('Skipping GitHub Actions Artifact upload (not in GH Actions env).');
    }

    // Octokit setup
    const token = process.env.GITHUB_TOKEN;
    if (!token) {
        console.warn('GITHUB_TOKEN not set.');
        if (!dryRun && (isMain || isPR)) { // Critical for actual publish
            console.error('GITHUB_TOKEN required for publishing.');
            process.exit(1);
        }
    }
    const octokit = token ? new Octokit({ auth: token }) : null;
    const owner = 'zitadel';
    const repo = 'zitadel';

    // 2. PR Logic: Comment on PR
    if (isPR && octokit && process.env.GITHUB_EVENT_NAME === 'pull_request') {
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
