# End-to-End Testing for Release Process

This document describes how to run E2E tests for the ZITADEL release process using a real forked repository.

## Overview

The E2E tests run the actual release process against a forked repository with real GitHub releases, Docker registries, and npm packages. This provides the most realistic testing environment.

## Prerequisites

### 1. Fork the Repository

Fork the ZITADEL repository to a test organization or personal account:
- Go to https://github.com/zitadel/zitadel
- Click "Fork"
- Choose your test organization

### 2. Create a Test Branch

In your fork, create a dedicated branch for E2E testing:
```bash
git checkout -b e2e-test-release
git push -u origin e2e-test-release
```

### 3. Generate GitHub Token

Create a Personal Access Token with the following scopes:
- `repo` (Full control of private repositories)
- `write:packages` (Upload packages to GitHub Package Registry)
- `delete:packages` (Delete packages from GitHub Package Registry - for cleanup)

Generate token at: https://github.com/settings/tokens/new

### 4. Set Up Environment Variables

```bash
export E2E_TEST_REPO="your-org/zitadel-fork"
export E2E_TEST_BRANCH="e2e-test-release"
export GH_TOKEN="ghp_your_token_here"
```

## Running E2E Tests

### Configuration Validation (Always Runs)

These tests check if your environment is properly configured:

```bash
pnpm nx test-e2e @zitadel/release
```

Output when not configured:
```
⚠ E2E tests against real fork are skipped.
  To run these tests, set the following environment variables:
    - E2E_TEST_REPO
    - E2E_TEST_BRANCH
    - GH_TOKEN
```

### Full E2E Test (With Configuration)

Once environment variables are set, the full suite runs:

```bash
pnpm nx run @zitadel/release:test-e2e
```

This will:
1. Clone your forked repository
2. Checkout the test branch
3. Install dependencies
4. Run release in dry-run mode
5. Verify the process completed successfully
6. Clean up test artifacts

## What Gets Tested

### ✅ Included in Tests

1. **Repository Cloning** - Verifies access to the fork
2. **Branch Setup** - Ensures test branch exists and is accessible
3. **Dependency Installation** - Validates package.json and pnpm-lock.yaml
4. **Release Dry Run** - Runs the full release process without publishing
5. **GitHub CLI Authentication** - Verifies gh CLI can access the repository
6. **Workspace Validation** - Checks project structure is correct
7. **Unit Tests** - Ensures all unit tests pass before release

### 🚫 Not Included (Requires Manual Testing)

These require actually publishing and should be tested manually on the fork:

1. **Docker Image Publishing** - Run with `--no-dryRun` to actually push images
2. **npm Package Publishing** - Run with `--no-dryRun` to actually publish packages
3. **GitHub Release Creation** - Run with conventional commits to create releases

## Manual Full Release Test

To manually test the complete release flow on your fork:

### Option 1: Using a Maintenance Branch

```bash
# In your fork, create a maintenance branch
git checkout -b v999.x
git push -u origin v999.x

# Run the release (this will create a real release)
export GH_TOKEN="your_token"
pnpm nx release @zitadel/release -- --no-dryRun

# Clean up
gh release delete v999.0.0 --yes
git tag -d v999.0.0
git push origin :refs/tags/v999.0.0
```

### Option 2: Using Main Branch (SHA-tagged images)

```bash
# On main branch, just run dry-run
pnpm nx release @zitadel/release -- --dryRun

# Or push SHA-tagged images (no GitHub release)
pnpm nx release @zitadel/release -- --no-dryRun
# Clean up Docker images manually from GitHub Container Registry
```

## Test Architecture

```
┌─────────────────────────────────────────────────────────┐
│ E2E Test (Local Machine)                                │
│                                                         │
│  1. Clone fork      ──────────────────────────┐         │
│  2. Setup branch                              │         │
│  3. Install deps                              │         │
│  4. Run release     ────────┐                 │         │
│  5. Verify results          │                 │         │
│  6. Cleanup                 │                 │         │
└─────────────────────────────┼─────────────────┼─────────┘
                              │                 │
                              ▼                 ▼
                    ┌──────────────────┐ ┌──────────────┐
                    │ GitHub           │ │ GitHub       │
                    │ - Releases       │ │ - Clone      │
                    │ - API calls      │ │ - Auth       │
                    └──────────────────┘ └──────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │ GitHub Container │
                    │ Registry         │
                    │ - Docker images  │
                    └──────────────────┘
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: E2E Release Tests

on:
  push:
    branches: [e2e-test-release]
  workflow_dispatch:

jobs:
  e2e-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: pnpm/action-setup@v2
      
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
          cache: 'pnpm'
      
      - name: Install dependencies
        run: pnpm install
      
      - name: Run E2E tests
        env:
          E2E_TEST_REPO: ${{ github.repository }}
          E2E_TEST_BRANCH: e2e-test-release
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: pnpm nx test-e2e @zitadel/release
```

## Troubleshooting

### Authentication Failed

```
Error: gh auth status failed
```

**Solution**: Ensure GH_TOKEN is set and has correct permissions:
```bash
echo $GH_TOKEN  # Should output your token
gh auth status  # Should show "Logged in"
```

### Repository Not Found

```
Error: repository not found
```

**Solution**: Verify the fork exists and token has access:
```bash
gh repo view $E2E_TEST_REPO
```

### Branch Doesn't Exist

```
Error: pathspec 'e2e-test-release' did not match
```

**Solution**: The test will create the branch automatically, or create it manually:
```bash
git checkout -b e2e-test-release
git push -u origin e2e-test-release
```

### Tests Skipped

```
⚠ E2E tests against real fork are skipped.
```

**Solution**: Set all required environment variables:
```bash
export E2E_TEST_REPO="your-org/zitadel-fork"
export E2E_TEST_BRANCH="e2e-test-release"  
export GH_TOKEN="ghp_..."
```

## Cleanup

### After Each Test Run

The tests automatically clean up:
- ✅ Test release (deleted via `gh release delete`)
- ✅ Test tags (deleted from git)
- ✅ Test directory (`/tmp/zitadel-e2e-test`)

### Manual Cleanup

If tests are interrupted, clean up manually:

```bash
# Delete test releases
gh release delete v999.* --yes --repo your-org/zitadel-fork

# Delete test tags
git push origin :refs/tags/v999.*

# Delete test directory
rm -rf /tmp/zitadel-e2e-test

# Delete Docker images (via GitHub web UI or API)
# Go to: https://github.com/orgs/your-org/packages
```

## Best Practices

1. **Use a Dedicated Test Organization** - Don't use production repositories
2. **Automate Cleanup** - Always clean up test artifacts
3. **Use Unique Tags** - Include timestamps to avoid conflicts
4. **Test in CI** - Set up automated E2E tests in GitHub Actions
5. **Monitor Costs** - Docker images and packages can accumulate

## Future Improvements

- [ ] Add actual publishing tests (currently only dry-run)
- [ ] Automate Docker image cleanup via API
- [ ] Test with different branch patterns (main, v1.x, v1.0.x)
- [ ] Validate Docker image labels and metadata
- [ ] Verify npm package.json metadata
- [ ] Test rollback scenarios
- [ ] Add performance benchmarks


## Running E2E Tests

### Quick Start

Run only the basic infrastructure tests (mock gh CLI):
```bash
pnpm nx test-e2e @zitadel/release
```

### Requirements

- Docker installed and running
- Sufficient disk space for Docker images
- Ports 5555 (Docker registry) and 4874 (npm registry) available

## Test Architecture

### 1. Local Docker Registry

The tests start a Docker registry on `localhost:5555`. To publish images there:

```typescript
// In docker-bake files, temporarily override the registry:
process.env.DOCKER_REGISTRY = 'localhost:5555';
```

After the test, verify images:
```bash
curl http://localhost:5555/v2/_catalog
curl http://localhost:5555/v2/zitadel-api/tags/list
```

### 2. Local npm Registry (Verdaccio)

The tests start Verdaccio on `localhost:4874`. To publish packages there:

```bash
# Configure npm to use local registry
npm config set registry http://localhost:4874/

# Run your publish command
pnpm publish

# Reset to default
npm config delete registry
```

Verify packages:
```bash
curl http://localhost:4874/-/all
curl http://localhost:4874/@zitadel/client
```

### 3. Mock GitHub CLI

A mock `gh` script is created that:
- Intercepts `gh release create` commands
- Intercepts `gh release upload` commands
- Stores release metadata in `/tmp/zitadel-release-e2e-test/github-releases/`
- Returns success for all operations

The mock is prepended to `PATH`, so all `gh` calls use it automatically.

## Extending the Tests

### Enable Docker Image Publishing Tests

To enable the skipped Docker tests, you need to:

1. **Modify docker-bake files** to use the local registry:

```typescript
// In release.e2e.test.ts
const originalDockerBake = readFileSync('release/docker-bake-release.hcl', 'utf-8');
const modifiedDockerBake = originalDockerBake.replace(
  /ghcr\.io\/[^\/]+\//g,
  `localhost:${DOCKER_REGISTRY_PORT}/`
);
writeFileSync(`${E2E_TEST_DIR}/docker-bake-release.hcl`, modifiedDockerBake);
```

2. **Run the release** with environment overrides:

```typescript
process.env.DOCKER_BAKE_FILE = `${E2E_TEST_DIR}/docker-bake-release.hcl`;
await main(['--no-dryRun']);
```

3. **Verify images** were pushed:

```typescript
const catalog = execSync(
  `curl -s http://localhost:${DOCKER_REGISTRY_PORT}/v2/_catalog`,
  { encoding: 'utf-8' }
);
const catalogData = JSON.parse(catalog);
expect(catalogData.repositories).toContain('zitadel-api');
expect(catalogData.repositories).toContain('login');
```

### Enable npm Package Publishing Tests

To enable the skipped npm tests:

1. **Configure npm** to use the local registry:

```typescript
execSync(`npm config set registry http://localhost:${NPM_REGISTRY_PORT}/`, { stdio: 'inherit' });
```

2. **Run nx release** with publish:

```typescript
// Mock the releasePublish from nx/release to publish to local registry
await main(['--no-dryRun']);
```

3. **Verify packages** were published:

```typescript
const packages = execSync(
  `curl -s http://localhost:${NPM_REGISTRY_PORT}/-/all`,
  { encoding: 'utf-8' }
);
const packagesData = JSON.parse(packages);
expect(packagesData).toHaveProperty('@zitadel/client');
expect(packagesData).toHaveProperty('@zitadel/proto');
```

## Test Isolation

Each test run:
- Uses a fresh `/tmp/zitadel-release-e2e-test` directory
- Starts new Docker containers with unique names
- Cleans up all containers and files after completion
- Restores the original environment variables

## Debugging

### View Docker Registry Contents

```bash
# List all repositories
curl http://localhost:5555/v2/_catalog

# List tags for a specific image
curl http://localhost:5555/v2/zitadel-api/tags/list

# Get image manifest
curl http://localhost:5555/v2/zitadel-api/manifests/v1.0.0
```

### View npm Registry Contents

```bash
# List all packages
curl http://localhost:4874/-/all | jq 'keys'

# Get package info
curl http://localhost:4874/@zitadel/client | jq '.'

# Get specific version
curl http://localhost:4874/@zitadel/client/1.0.0
```

### View Mock GitHub Releases

```bash
ls /tmp/zitadel-release-e2e-test/github-releases/
cat /tmp/zitadel-release-e2e-test/github-releases/v1.0.0.json
```

### Keep Environment Running

To keep the test environment running for manual testing:

```typescript
// Add at the end of your test
await new Promise(() => {}); // Never resolves, keeps containers running
```

Then in another terminal:
```bash
# Test Docker registry
docker push localhost:5555/my-test-image:latest

# Test npm registry
npm config set registry http://localhost:4874/
npm publish
```

## CI Integration

For CI environments, you may want to:

1. **Pre-pull images** to avoid timeout:
```bash
docker pull registry:2
docker pull verdaccio/verdaccio
```

2. **Increase timeouts** if building is slow:
```typescript
beforeAll(async () => {
  // ...
}, 120000); // 2 minutes
```

3. **Run in isolation** to avoid port conflicts:
```bash
pnpm nx test-e2e @zitadel/release --maxWorkers=1
```

## Troubleshooting

### Port Already in Use

If ports 5555 or 4874 are already in use:
```bash
# Find and kill the process
lsof -ti:5555 | xargs kill -9
lsof -ti:4874 | xargs kill -9
```

### Docker Containers Not Stopping

```bash
docker stop zitadel-e2e-registry zitadel-e2e-verdaccio
docker rm zitadel-e2e-registry zitadel-e2e-verdaccio
```

### Permission Denied on Mock gh Script

```bash
chmod +x /tmp/zitadel-release-e2e-test/bin/gh
```

## Future Improvements

- [ ] Add tests for GitHub release notes generation
- [ ] Test multi-platform Docker builds (arm64, amd64)
- [ ] Verify checksums.txt generation
- [ ] Test rollback scenarios
- [ ] Add performance benchmarks
- [ ] Test with different git branch patterns
- [ ] Validate Docker image labels
- [ ] Verify npm package metadata
