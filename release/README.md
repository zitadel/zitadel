# Releasing Zitadel

Releasing Zitadel is portable to forks on other organizations than @zitadel.
This allows to easily change, publish and use the Zitadel code in any GitHub organization.
Make sure that you comply with the changed codes license, different projects are differently licensed.

The release process is not only runnable with GitHub actions, but also locally or anywhere else.
The default `--dry-run` mode creates all artifacts without publishing them.
The release logic is written in code and unit tested.
Validations prevent accidentally publishing wrong releases.

## Creating a Release

In order to push artifacts, the env variable GITHUB_TOKEN needs to be set.
The token has to be a classic PAT with `package:write` scope.

```bash
export GITHUB_TOKEN=$(cat /tmp/my-classic-pat)
# Creating a release on a fork in another organization
pnpm release --github-repo my-org/zitadel --no-dry-run
```

## The Standard Release Process

- The release script is triggered on manual workflow dispatch events as well as on every commit to main.
- If the release script is triggered on a major or minor maintenance branch, the version is bumped according to conventional commits.
  Release branches for regular or maintenance releases are determined by the patterns `v[0-9].x` and `v[0-9].[0-9].x`, like `v4.x` or `v4.4.x`.
  The following artifacts are published when a version is bumped:
  - A new GitHub release is created with tarballs for the API and the Login with a checksums.txt file.
  - NPM packages are pushed:
    - @${ZITADEL_RELEASE_GITHUB_ORG}/client@${ZITADEL_RELEASE_VERSION}
    - @${ZITADEL_RELEASE_GITHUB_ORG}/proto@${ZITADEL_RELEASE_VERSION}
  - Docker images are pushed:
    - ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api:${ZITADEL_RELEASE_VERSION}
    - ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api-debug:${ZITADEL_RELEASE_VERSION}
    - ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/login:${ZITADEL_RELEASE_VERSION}
  - If the bumped version is the highest regular semantic version in the repository, the Docker images are additionally pushed with the `latest` tag:
    - ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api:latest
    - ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api-debug:latest
    - ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/login:latest
  - If the release script is not triggered on a major or minor maintenance branch, only SHA tagged Docker images are pushed.
    Apart from this, nothing else is released.
    On every commit to main, the release script is triggered in CI and SHA Docker images are pushed.
    - ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api:${ZITADEL_RELEASE_REVISION}
    - ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/login:${ZITADEL_RELEASE_REVISION}

## Customizing the Release Process

Releasing is implemented using the [Nx Release programmatic API](https://nx.dev/docs/guides/nx-release/programmatic-api).
The repo-wide release implementations are done in the @zitadel/release project,
the project-specific implementations are done in special Nx targets called `nx-release-publish`.
If in accordance with the [LICENSE](LICENSE), the release process can be changed, linted and tested.

### Linting

Run TypeScript type checking:
```bash
pnpm nx run @zitadel/release:lint
```

### Testing

Install dependencies first:
```bash
pnpm install
```

Run the tests:
```bash
pnpm nx run @zitadel/release:test
```
