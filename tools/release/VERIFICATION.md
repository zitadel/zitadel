# Verifying the Release Process

## Prerequisites
- Ensure you differ in a devcontainer or have `pnpm` and `nx` (v22+) installed.
- Ensure `GH_TOKEN` is set if you want to test artifact upload (dry-run will skip actual upload).

## Dry Run Release Artifacts
This tests the artifact generation and upload logic (mocked).
```bash
pnpm nx run release-tools:release-artifacts --args="1.0.0-test --dry-run"
```

## Verify NX Release Configuration
To see what version NX would calculate and what changelog it would generate:
```bash
pnpm nx release --dry-run
```
(This should automatically detect `v*` tags and propose a new version based on commits).

## Verifying `checkAllBranchesWhen`
This feature ensures tags on other branches (like `next`) are considered when determining the current version.
To verify:
1. Ensure you are on a feature branch or `next`.
2. Run `pnpm nx release --dry-run --verbose`.
3. Check the output logs for "Searching for git tags on all branches" or similar indication that it's looking beyond HEAD.

## Docker Container Images
The `release-tools:release-artifacts` target currently logs the images that would be tagged.
Ensure the logic in `tools/release/release-artifacts.ts` matches the expected image names:
- `ghcr.io/zitadel/zitadel`
- `europe-docker.pkg.dev/zitadel/zitadel/zitadel`
