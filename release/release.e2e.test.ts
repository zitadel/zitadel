/**
 * End-to-end test for the release process.
 * 
 * This test runs the actual release process against a real forked repository.
 * 
 * Prerequisites:
 * - E2E_TEST_REPO: The forked repo (e.g., "test-org/zitadel")
 * - E2E_TEST_BRANCH: The branch to test (e.g., "e2e-test-release")
 * - GH_TOKEN: GitHub token with repo and package write permissions
 * 
 * What this test does:
 * 1. Fork/sync and clone/fetch the repository
 * 2. Checkout or create the test branch
 * 3. Pull latest changes
 * 4. Run the release process
 * 5. Verify GitHub release was created with correct assets
 * 6. Download and test the binary for current platform
 * 7. Verify Docker images exist with correct tags
 * 8. Clean up test release and artifacts
 */

import { execSync, spawnSync } from 'node:child_process';
import { existsSync, mkdirSync } from 'node:fs';
import { platform, arch } from 'node:os';
import { beforeAll, describe, expect, test } from 'vitest';

const E2E_TEST_DIR = '/tmp/zitadel-e2e-test';
const TEST_REPO_ENV_VAR = "ZITADEL_RELEASE_E2E_TEST_REPO";
const TEST_BRANCH_ENV_VAR = "ZITADEL_RELEASE_E2E_TEST_BRANCH";
const EXPECT_RELEASE_VERSION_ENV_VAR = "ZITADEL_RELEASE_E2E_EXPECT_VERSION";
const EXPECT_GITHUB_RELEASE_ENV_VAR = "ZITADEL_RELEASE_E2E_EXPECT_GITHUB_RELEASE";
const EXPECT_LATEST_RELEASE_ENV_VAR = "ZITADEL_RELEASE_E2E_EXPECT_LATEST";
const REQUIRED_ENV_VARS = [TEST_REPO_ENV_VAR, TEST_BRANCH_ENV_VAR, EXPECT_RELEASE_VERSION_ENV_VAR, EXPECT_GITHUB_RELEASE_ENV_VAR, EXPECT_LATEST_RELEASE_ENV_VAR];

// Check if E2E tests should run
// Ensure all required environment variables are set
// Ensure E2E_TEST_REPO is not on the zitadel organization
const SHOULD_RUN_E2E = REQUIRED_ENV_VARS.every(envVar => process.env[envVar]) && !process.env[TEST_REPO_ENV_VAR]?.startsWith('zitadel/');

// Get platform-specific info
const GOOS = platform() === 'darwin' ? 'darwin' : platform() === 'win32' ? 'windows' : 'linux';
const GOARCH = arch() === 'x64' ? 'amd64' : arch() === 'arm64' ? 'arm64' : arch();
const GITHUB_RELEASE_GITHUB_ORG = process.env[TEST_REPO_ENV_VAR]?.split('/')[0] || '';

// Helper to run commands in the test repo
function runInTestRepo(command: string, options: { silent?: boolean; ignoreError?: boolean } = {}): string {
    try {
        return execSync(command, {
            cwd: E2E_TEST_DIR,
            encoding: 'utf-8',
            stdio: options.silent ? 'pipe' : 'inherit',
            env: {
                ...process.env,
                ZITADEL_RELEASE_GITHUB_ORG: GITHUB_RELEASE_GITHUB_ORG,
            },
        }).trim();
    } catch (error: any) {
        if (options.ignoreError) {
            return '';
        }
        console.error(`Command failed: ${command}`);
        console.error(`Error: ${error.message}`);
        if (error.stdout) console.error(`Stdout: ${error.stdout}`);
        if (error.stderr) console.error(`Stderr: ${error.stderr}`);
        throw error;
    }
}

// Helper to run commands and capture output
function runCommand(command: string, cwd: string = E2E_TEST_DIR): { stdout: string; stderr: string; exitCode: number } {
    const result = spawnSync(command, {
        cwd,
        encoding: 'utf-8',
        shell: true,
        env: {
            ...process.env,
            ZITADEL_RELEASE_GITHUB_ORG: GITHUB_RELEASE_GITHUB_ORG,
        },
    });
    return {
        stdout: result.stdout || '',
        stderr: result.stderr || '',
        exitCode: result.status || 0,
    };
}

describe.skipIf(SHOULD_RUN_E2E)('Release E2E (Real Fork)', () => {
    test('should skip E2E tests when environment variables are not set correctly', () => {
        console.log('Skipping E2E tests. To run them, ensure the following environment variables are set:');
        console.log(`- ${TEST_REPO_ENV_VAR}`);
        console.log(`- ${TEST_BRANCH_ENV_VAR}`);
        console.log(`- ${EXPECT_RELEASE_VERSION_ENV_VAR}`);
        console.log(`- ${EXPECT_GITHUB_RELEASE_ENV_VAR}`);
        console.log(`- ${EXPECT_LATEST_RELEASE_ENV_VAR}`);
        console.log('Also ensure that the test repository is not under the "zitadel" organization.');
    });
});

describe.skipIf(!SHOULD_RUN_E2E)('Release E2E (Real Fork)', () => {
    const testRepo = process.env[TEST_REPO_ENV_VAR]!;
    const testBranch = process.env[TEST_BRANCH_ENV_VAR]!;
    const testcodeBranch = execSync('git rev-parse --abbrev-ref HEAD').toString().trim();

    beforeAll(async () => {
        console.log('\n🔧 Setting up E2E test with real fork...');
        console.log(`  Repository: ${testRepo}`);
        console.log(`  Release Branch: ${testBranch}`);
        console.log(`  Testcode Branch: ${testcodeBranch}`);
        console.log(`  Platform: ${GOOS}/${GOARCH}`);

        // Step 1: Clone or fetch the repository
        if (existsSync(E2E_TEST_DIR)) {
            console.log('  Fetching latest changes...');
            runInTestRepo('git fetch --all');
        } else {
            console.log('  Cloning forked repository...');
            mkdirSync(E2E_TEST_DIR, { recursive: true });
            execSync(`git clone https://github.com/${testRepo}.git ${E2E_TEST_DIR}`, {
                encoding: 'utf-8',
                stdio: 'inherit',
                env: {
                    ...process.env,
                    GH_TOKEN: process.env.GH_TOKEN,
                },
            });
        }

        // Step 2: Checkout or create the test branch
        console.log(`  Checking out branch ${testBranch}...`);
        try {
            runInTestRepo(`git checkout ${testBranch}`, { ignoreError: true });
        } catch {
            // Branch doesn't exist locally, try to track remote
            try {
                runInTestRepo(`git checkout -b ${testBranch} origin/${testBranch}`, { ignoreError: true });
            } catch {
                // Remote branch doesn't exist either, create new branch
                console.log(`  Creating new branch ${testBranch}...`);
                runInTestRepo(`git checkout -b ${testBranch}`);
            }
        }

        // Step 3: Pull latest changes
        console.log('  Pulling latest changes...');
        runInTestRepo(`git pull origin ${testBranch}`, { ignoreError: true });

        console.log('✓ E2E test environment ready\n');

        // Step 4: Merge the branch containing the release code from where we initiated the test in the main repository
        runInTestRepo(`git merge ${testcodeBranch}`, { ignoreError: true });
    }, 180000); // 3 minute timeout for setup

    test('should install dependencies', () => {
        console.log('  Installing dependencies...');
        runInTestRepo('pnpm install --frozen-lockfile');
        console.log('✓ Dependencies installed');
    }, 300000); // 5 minute timeout

    test('should run unit tests before release', () => {
        console.log('  Running unit tests...');
        runInTestRepo('pnpm nx run @zitadel/release:test-unit');
        console.log('✓ Unit tests passed');
    }, 120000); // 2 minute timeout

    test('should run release process in dry run mode', async () => {
        // Run release in dry-run mode first to validate
        console.log('  Testing with dry-run...');
        runInTestRepo('pnpm nx run @zitadel/release:release -- --dry-run');
        console.log('  ✓ Dry-run completed successfully');
    }, 120000); // 2 minute timeout

    test('should run full release process', async () => {
        console.log('\n🚀 Running full release process...')
        runInTestRepo('pnpm nx run @zitadel/release:release -- --no-dry-run')
        console.log('✓ Full release process completed');
    }, 600000); // 10 minute timeout

    test('should verify Docker images exist with correct tags', async () => {
        console.log('\n🐳 Verifying Docker images...');

        // Extract registry and repository from test repo
        const [owner] = testRepo.split('/');
        const registry = 'ghcr.io';
        const testTag = process.env[EXPECT_RELEASE_VERSION_ENV_VAR];

        // Expected images
        const expectedImages = [
            { name: 'zitadel', tag: testTag },
            { name: 'login', tag: testTag },
        ];

        // Add 'latest' tag if this is a main branch release
        if (process.env[EXPECT_LATEST_RELEASE_ENV_VAR] == 'true') {
            expectedImages.push(
                { name: 'zitadel', tag: 'latest' },
                { name: 'login', tag: 'latest' }
            );
        }

        for (const { name, tag } of expectedImages) {
            const imageName = `${registry}/${owner}/${name}:${tag}`;
            console.log(`  Checking: ${imageName}`);

            // Pull the image to verify it exists
            const pullResult = runCommand(`docker pull ${imageName}`);
            expect(pullResult.exitCode).toBe(0);
            console.log(`  ✓ Image exists: ${imageName}`);

            // Run container with --version
            console.log(`  Testing ${name} container with --version...`);
            const runResult = runCommand(`docker run --rm ${imageName} --version`);

            expect(runResult.exitCode).toBe(0);
            expect(runResult.stdout).toContain(testTag);
            console.log(`  ✓ Container version: ${runResult.stdout.trim()}`);

            // Clean up the pulled image
            runCommand(`docker rmi ${imageName}`, E2E_TEST_DIR);
        }

        console.log('✓ All Docker images verified');
    }, 600000); // 10 minute timeout for Docker operations

    test.skipIf(process.env[EXPECT_GITHUB_RELEASE_ENV_VAR] === 'true')('should verify GitHub release exists with correct assets', async () => {
        console.log('Verifying GitHub release...');
        const testTag = process.env[EXPECT_RELEASE_VERSION_ENV_VAR]!;

        // Get release information
        const releaseInfo = runInTestRepo(
            `gh release view ${testTag} --json tagName,assets,body`,
            { silent: true }
        );
        const release = JSON.parse(releaseInfo);

        // Verify release exists
        expect(release.tagName).toBe(testTag);
        console.log(`  ✓ Release ${testTag} exists`);

        // Verify expected assets are present
        const assetNames = release.assets.map((a: any) => a.name);
        console.log(`  Found ${assetNames.length} assets`);

        // Check for platform-specific API tarballs
        const expectedBinaryPattern = new RegExp(`zitadel-${GOOS}-${GOARCH}.tar.gz`);
        const hasBinary = assetNames.some((name: string) => expectedBinaryPattern.test(name));
        expect(hasBinary).toBe(true);
        console.log(`  ✓ Binary for ${GOOS}/${GOARCH} found`);

        // Check for checksums file
        expect(assetNames).toContain('checksums.txt');
        console.log('  ✓ Checksums file found');

        // Check for login tarball
        expect(assetNames.some((name: string) => name.includes('zitadel-login.tar.gz'))).toBe(true);
        console.log('  ✓ Login tarball found');
    }, 60000);

    test.skipIf(process.env[EXPECT_GITHUB_RELEASE_ENV_VAR] === 'true')('should download and test binary', async () => {
        console.log('\n📥 Downloading and testing binary...');
        const testTag = process.env[EXPECT_RELEASE_VERSION_ENV_VAR]!;

        const artifactsDir = `${E2E_TEST_DIR}/.e2e-artifacts`;
        mkdirSync(artifactsDir, { recursive: true });

        // Download the binary for current platform
        const binaryPattern = `*${GOOS}-${GOARCH}*`;
        console.log(`  Downloading binary matching: ${binaryPattern}`);
        runInTestRepo(
            `gh release download ${testTag} --pattern "${binaryPattern}" --dir ${artifactsDir}`
        );

        // Find the downloaded tarball
        const files = execSync(`ls ${artifactsDir}`, { encoding: 'utf-8' }).split('\n').filter(Boolean);
        const tarballFile = files.find(f => f.includes(GOOS) && f.includes(GOARCH));
        expect(tarballFile).toBeDefined();

        const tarballPath = `${artifactsDir}/${tarballFile}`;
        console.log(`  ✓ Downloaded: ${tarballFile}`);

        // Unpack the tarball
        console.log('  Unpacking tarball...');
        runInTestRepo(`tar -xzf ${tarballPath} -C ${artifactsDir}`);
        console.log('  ✓ Tarball unpacked');

        // Run binary with --version flag
        console.log('  Testing binary with --version...');
        const result = runCommand(`${tarballPath} --version`, artifactsDir);

        expect(result.exitCode).toBe(0);
        expect(result.stdout).toContain(testTag);
        console.log(`  ✓ Binary version: ${result.stdout.trim()}`);
    }, 120000);
})
