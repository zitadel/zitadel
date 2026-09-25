---
name: npm-dependency-audit
description: Investigate and fix npm dependency vulnerabilities reported by `pnpm audit` in this pnpm/Nx monorepo. Fixes each finding at the package that introduces it (direct bump, in-range lockfile refresh or parent bump), keeps `pnpm.overrides` as a documented last resort, prunes overrides that are no longer needed and prepares a version-bump PR with manual test notes. Use when asked to fix npm/pnpm audit findings, Dependabot or GitHub security alerts for JavaScript packages, or to clean up pnpm overrides.
---

# npm dependency audit

Use this workflow to fix `pnpm audit` findings in this repository. Fix each finding at the package that introduces it. Treat an override as a documented exception, not the default fix.

Run every command from the repository root. Workspaces are listed in `pnpm-workspace.yaml`, and overrides live under `overrides:` in the same file.

## Principles

1. **Fix at the origin.** Prefer these fixes, in this order:
   1. bump the direct dependency,
   2. refresh the transitive dependency inside its existing semver range,
   3. bump the parent package that pins the vulnerable version,
   4. add an override.
2. **Overrides are temporary.** Every override needs a comment that names the parent that pins the old version. Anyone can then check later whether it can be removed.
3. **Review existing overrides every time.** Remove all of them, reinstall and re-add only the ones that audit still needs.
4. **Keep the lockfile diff small.** Don't delete `pnpm-lock.yaml` to regenerate it. Update only the packages you are fixing, so reviewers can follow the change.
5. **Keep the lockfile in sync.** Run `pnpm install` after every edit to a `package.json` or to `overrides`, and commit `pnpm-lock.yaml` together with it. CI installs with `--frozen-lockfile` and fails when they disagree.
6. **Write the PR as a dependency update.** This is a public repository. Describe version bumps, behavior changes and what to test. Leave out exploit details, attack scenarios and "how this could be abused" commentary. Linking to a public advisory is fine but not needed.

## 1. Collect the findings

```bash
git switch main && git pull --ff-only
pnpm install
pnpm audit --json > /tmp/audit.json
node .agents/skills/npm-dependency-audit/scripts/summarize-audit.mjs /tmp/audit.json
```

The summary shows each vulnerable package, its installed and fixed versions and its dependency paths. In a path like `apps__docs>raw-loader>webpack>schema-utils>ajv>fast-uri`:

- the first segment is the workspace project (`apps__docs` = `apps/docs`, `.` = root `package.json`),
- the second segment is the **direct dependency** in that project's `package.json`,
- the last segment is the vulnerable package.

## 2. Classify every finding

For each vulnerable package, find out which package in the path controls its version:

```bash
pnpm why -r <vulnerable-package>                          # who pulls it in, at which versions
npm view <parent>@<installed-version> dependencies.<pkg>   # the range the parent allows
npm view <parent> version                                  # latest parent release
npm view <parent>@latest dependencies.<pkg>                # does the latest parent allow the fix?
```

| Situation | Fix |
|---|---|
| The direct dependency at the top of the path is unused (`git grep -l "<pkg>" -- <project> ':!**/package.json'` finds nothing, and it isn't a peer dependency of another package) | Remove it with `pnpm --filter <project> remove <pkg>`. This removes its whole subtree. |
| The vulnerable package is a direct dependency | Raise its version in `package.json`, and raise the lower bound as well (`^4.1.6` → `^4.1.11`) so the fix is explicit. Then run `pnpm install`. |
| Transitive, and the parent's range already allows the fixed version | Refresh it in the lockfile (see below). No `package.json` change needed. |
| Transitive, the parent pins an old version, but a newer parent release allows the fix | Bump the parent. If that's a major version or 0.x minor bump, check its breaking changes (step 4). |
| The parent pins it exactly, even in its latest release (also check `npm view <parent> dist-tags` for a `next`/beta that fixes it) | Add an override with a comment (step 3). |
| The fix only exists in a major version of the vulnerable package | Check whether the parent works with that major. If unsure, keep the finding open and mention it in the PR instead of forcing it. |

### Refreshing a transitive dependency in the lockfile

```bash
pnpm update -r --depth Infinity <package>
pnpm why -r <package>        # confirm the old version is gone
```

Pitfalls:

- **Run one package per `pnpm update` call.** Passing several names in one call doesn't reliably update all of them. Verify each one with `pnpm why`.
- **Never pass a version range** (`pnpm update js-yaml@3`) when the package is also a direct dependency anywhere. pnpm then rewrites that project's `package.json` spec to the given range. Check `git diff -- '**/package.json'` after every update.
- If an old copy survives, `pnpm why -r <package>@<old-version>` shows which parent still holds it. That parent often has a patch release that widens the range; update it the same way.
- Don't edit `pnpm-lock.yaml` by hand.

## 3. Review and prune overrides

1. Note the current `overrides:` block in `pnpm-workspace.yaml` (and any `pnpm.overrides` in the root `package.json`).
2. Check `git log -S'"<package>@' -- pnpm-workspace.yaml` and the related PR descriptions to learn why each entry was added.
3. Remove **all** overrides, apply your direct bumps and run `pnpm install`.
4. Run `pnpm audit` again. Everything that was fixed at the origin is gone. What remains needs either a parent bump or an override.
5. Re-add only the overrides that are still needed. Scope each one to the vulnerable range, not the whole package, and name the pinning parent:

```yaml
overrides:
  # nx pins smol-toml 1.6.1 (latest nx 22.x and 23.x)
  "smol-toml@<1.7.1": "^1.7.1"
```

6. Run `pnpm install` so the lockfile records the overrides, then `pnpm audit` again. It must report no findings.
7. For each override you removed, confirm the resolved version is still at or above the floor the override enforced (`pnpm why -r <package>`). List the removed overrides in the PR.

An override forces a version the parent never tested with. Stay inside the same major version, and run the parent's build or tests after adding one.

A range override also applies to **direct** dependencies whose declared range overlaps the selector. pnpm then writes the overridden range into that project's `importers` entry in `pnpm-lock.yaml`. For example, with the `esbuild@>=0.27.3 <0.28.1` override, `"esbuild": "^0.28.0"` is recorded as `specifier: ^0.28.1`. Frozen installs accept this, but the manifest and lockfile no longer match, which confuses reviewers and review bots. Raise the lower bound of each such direct dependency to the override's floor:

```bash
git grep -nE '"<package>": "' -- '**/package.json'     # direct users of the overridden package
pnpm install                                            # after raising the ranges
```

## 4. Assess breaking changes

For every bump that is a major, a 0.x minor or a large minor version, read the changes between the old and new versions:

```bash
gh release list -R <owner>/<repo> --limit 30
gh release view <tag> -R <owner>/<repo>
curl -sL https://raw.githubusercontent.com/<owner>/<repo>/main/CHANGELOG.md | less
npm view <package> repository.url          # find the repository
```

Then check whether the repository uses the affected APIs:

```bash
grep -rn "<api-or-config-key>" apps console packages tests --exclude-dir=node_modules --exclude-dir=.next --exclude-dir=dist
```

Watch for:

- **engine requirements.** pnpm only warns when a package's `engines` doesn't match the running Node version (this repository doesn't set `engine-strict`). Installs succeed, and the failure only shows up at runtime. Check the new version and any new dependencies it brings in, for example `npm view <package>@<new> engines.node` and `npm view <new-dependency> engines.node`. Compare against **every** place that declares a Node version:
  - `.nvmrc`
  - the Node.js prerequisite in `CONTRIBUTING.md`
  - `node-version` in `.github/workflows/*.yml`
  - `.devcontainer/docker-compose.yaml`
  - the `FROM` lines in `apps/login/Dockerfile*`

  If a bump raises the minimum above one of them, update the stale declaration in the same PR and explain why in the description. Or keep the older release if the declared version must stay supported. If the declarations already disagree with each other, point that out in the PR,
- **ESM-only packages** consumed from CommonJS code,
- **changed defaults**, for example a test container's wait strategy or a logger's flush interval,
- **anything operators can observe.** This includes telemetry attribute names, log format, HTTP headers, config keys, the Docker image and `output: "standalone"`. These changes affect self-hosters. Ask the maintainer whether to ship them in the same PR or separately; if it's separate, keep the override for now,
- framework changes affecting **routing, middleware/proxy, image optimization, server actions, CSP/nonce or the build output** (Next.js, Angular).

## 5. Verify

Use only Nx targets that exist (see the root `AGENTS.md`, or run `pnpm nx show project <project>`):

```bash
pnpm audit                                   # expect: No known vulnerabilities found
pnpm install --frozen-lockfile               # what CI runs; fails if a manifest and the lockfile disagree
pnpm nx run-many -t build test-unit lint -p @zitadel/login @zitadel/client @zitadel/console @zitadel/docs
```

`@zitadel/console` has no `test-unit` target; Nx skips it. For bumps that touch the login's Docker image, containers or telemetry, also run the dockerized login tests (`pnpm --filter @zitadel/login test-dockerized`, needs Docker). If some of these checks can't run locally, say so in the PR.

## 6. Open the PR

- Title: `chore(deps): <short summary>`. `deps` is an allowed scope in `.github/semantic.yml`.
- Body: use `.github/pull_request_template.md` and keep every section.
  - **Which Problems Are Solved**: dependencies fell behind their fixed releases, and some overrides were no longer needed or were missing context.
  - **How the Problems Are Solved**: a table of bumps (`package | from | to | where | type`), then the in-range lockfile refreshes, the removed overrides and the remaining overrides with their pinning parent.
  - **Additional Changes**: a section on breaking or behavior-changing bumps. For each one: what changed upstream, whether the repository is affected (with evidence, such as "all containers set an explicit wait strategy"), and **manual checks** with concrete paths, for example `/ui/v2/login/loginname` or `/docs/apis/...`.
  - **Additional Context**: follow-ups, for example a major bump that was deliberately deferred.
- End the description with the agent attribution line if your tooling requires one.
