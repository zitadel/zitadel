# Release Script Tests

This directory contains the refactored release script with comprehensive test coverage.

## Running Tests

Install dependencies first:
```bash
pnpm install
```

Run the tests:
```bash
pnpm nx test-unit @zitadel/release
# or
pnpm nx run @zitadel/release:test-unit
```

Run tests in watch mode during development:
```bash
pnpm nx test-unit-watch @zitadel/release
```

## Running the Release

Execute the release process:
```bash
pnpm release
# or
pnpm nx release @zitadel/release
```

With options:
```bash
pnpm nx release @zitadel/release -- --no-dryRun --verbose
```

## Linting

Run TypeScript type checking:
```bash
pnpm nx lint @zitadel/release
```

## Test Coverage

The test suite covers:

- **`shouldUseConventionalCommits`**: Branch pattern validation for maintenance branches
- **`setupEnvironmentVariables`**: Environment variable configuration for different release scenarios
- **`executeDockerBuild`**: Docker build command construction with correct targets and bake files
- **`parseReleaseOptions`**: Command-line argument parsing
- **`executeRelease`**: Full release orchestration with mocked external calls

## Mocking Strategy

The test suite mocks:
- `node:child_process` - to prevent actual git/docker commands
- `nx/release` - to simulate Nx release operations
- `process.env` - to test environment variable setup in isolation

This allows testing the business logic without side effects.
