import { Octokit } from 'octokit';
import { execSync } from 'child_process';
import * as fs from 'fs';
import * as path from 'path';

async function main() {
    const version = process.env.NX_RELEASE_VERSION || process.argv[2];
    if (!version) {
        console.error('No version provided. Usage: ts-node release-artifacts.ts <version>');
        process.exit(1);
    }

    const dryRun = process.env.NX_DRY_RUN === 'true' || process.argv.includes('--dry-run');

    console.log(`Preparing release artifacts for version: ${version}`);

    // 1. Build Artifacts
    console.log('Running pack target...');
    if (!dryRun) {
        execSync('pnpm nx run pack', { stdio: 'inherit' });
    } else {
        console.log('[Dry-Run] Would run: pnpm nx run pack');
    }

    // 2. Upload to GitHub Release
    const token = process.env.GITHUB_TOKEN;
    if (!token) {
        console.warn('GITHUB_TOKEN not set. Skipping GitHub Release upload.');
    } else {
        const octokit = new Octokit({ auth: token });
        const owner = 'zitadel';
        const repo = 'zitadel'; // Adjust if needed
        const tag = `v${version}`; // Assuming v-prefix

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
                console.log('Release not found, waiting or skipping...'); // nx release should have created it
            }

            if (release) {
                console.log(`Found release ${release.id}. Uploading artifacts...`);
                const artifactsDir = path.join(process.cwd(), '.artifacts/pack');
                if (fs.existsSync(artifactsDir)) {
                    const files = fs.readdirSync(artifactsDir).filter(f => f.endsWith('.tar.gz') || f.endsWith('.zip') || f === 'checksums.txt');
                    for (const file of files) {
                        const filePath = path.join(artifactsDir, file);
                        const data = fs.readFileSync(filePath);
                        if (!dryRun) {
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
                                data: data as any, // octokit types valid here
                            });
                            console.log(`Uploaded ${file}`);
                        } else {
                            console.log(`[Dry-Run] Would upload ${file}`);
                        }
                    }
                } else {
                    console.warn('Artifacts directory not found.');
                }
            }
        } catch (error) {
            console.error('Error uploading artifacts:', error);
        }
    }

    // 3. Docker Tagging
    console.log('Tagging Docker images...');
    const images = [
        'ghcr.io/zitadel/zitadel', // Example, need to match actual image names from workflows
        'europe-docker.pkg.dev/zitadel/zitadel/zitadel',
    ];

    // Need to implement the logic to retag existing images (from sha?) to version
    // This depends on how the build was done.
    // The pack target might not build docker images? 
    // existing pack.yml builds docker images using sha tags.

    // For now, let's assume we run this after standard CI build or we implement docker build here.
    // The User wanted "locally testable".
    // Local tests might not push to registry.

}

main().catch(e => {
    console.error(e);
    process.exit(1);
});
