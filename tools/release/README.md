# Release Process

This directory contains the tooling and configuration for the ZITADEL release process. The release pipeline handles versioning, changelog generation, artifact creation, and publishing (Docker images and GitHub Releases).

## Quick Start

To run a release (or preview) locally:

```bash
# Preview / Dry-Run
pnpm release

# Detailed output
pnpm release --verbose --no-tui
```

> **Note**: Always use `pnpm release` (defined in `package.json`). Do **not** run `nx run @zitadel/release:release` directly, as it requires the Nx Daemon to be disabled to prevent hangs.

## Release Modes

The release script detects the environment and behaves accordingly:

### 1. Preview Mode (Dry-Run)
*   **Trigger**: Default behavior when running locally or on non-main branches (e.g., Pull Requests).
*   **Actions**:
    *   Calculates a preview version (e.g., `2.1.0-feature-branch`).
    *   Generates a changelog preview.
    *   Builds all artifacts (binaries, archives).
    *   **Docker**: Builds and pushes Docker images (e.g., `ghcr.io/zitadel/test-api:<version>`).
    *   **GitHub**: Comments on the PR with a summary of the release plan.

### 2. Live Mode (Production)
*   **Trigger**: Only when `CI_RELEASE=true` (set in CI on the `main` branch).
*   **Actions**:
    *   Commits version bumps and changelogs to Git.
    *   Creates and pushes a Git Tag.
    *   Creates a generic GitHub Release.
    *   Uploads all artifacts to the GitHub Release.
    *   **Docker**: Builds and pushes `linux/amd64` and `linux/arm64` images to GHCR with the `latest` tag.

## Configuration

*   **`tools/release/main.ts`**: The main orchestration script. Handles the logic for versioning, changelog generation, and invoking the build targets.
*   **`project.json`**: Defines the Nx targets for the release tool.
*   **`package.json`**: Contains the `"release"` script which explicitly sets `NX_DAEMON=false`.

## Environment Variables

| Variable | Description | Required? |
| :--- | :--- | :--- |
| `GITHUB_TOKEN` | GitHub PAT for commenting on PRs and publishing releases. | Yes (for CI/Publish) |
| `CI_RELEASE` | Set to `true` to enable Live Mode. | No (Default: false) |
| `NX_VERBOSE_LOGGING` | Set to `true` for detailed Nx logs. | No |

## Troubleshooting

### Script Hangs
If the release script hangs, especially on repeated runs, it is likely due to the **Nx Daemon**.
**Solution**: Ensure you are using `pnpm release`. This script sets `NX_DAEMON=false` for the runner process, ensuring a clean state for every run.

### Missing Artifacts
If artifacts are not uploading, check that the `pack` targets in `apps/api` and `apps/login` are executing correctly and outputting to `.artifacts/pack`.
