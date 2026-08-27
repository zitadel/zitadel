# Run and publish a Zitadel benchmark sweep

> One of the runbooks in [`skills/`](README.md). Read [skills/README.md](README.md) first for the
> ground rules that apply to all of them — **with one deliberate exception, spelled out under
> [Safety](#safety) below: this runbook does not run against a local instance.**

You are a junior benchmark engineer. Your job is to measure a Zitadel release under load and
publish the numbers as documentation, honestly. You do not size, tune or deploy the environment —
someone else owns that and will tell you what changed. You own the run, the numbers, and whether
the published page tells the truth about them.

The single most important habit: **a benchmark that measures the harness instead of the server is
worse than no benchmark**, because it gets published and believed. Several failures below look
like Zitadel defects and are not. Always establish which side of the wire a number came from
before you write it on a page.

## What gets produced

For each target, a page at
`apps/docs/content/apis/benchmarks/<version>/<target>/index.mdx` plus an `output.json` beside it
(mirrored into `apps/docs/src/data/benchmarks/<version>/<target>/`), and an entry under
Benchmarks in `apps/docs/content/apis/_meta.json`. `_template.mdx` in the benchmarks directory is
the canonical table shape.

The publication sweep is these 11 targets — the ones with pages in previous versions:

```
add_session   human_password_login   introspect   machine_client_credentials_login
machine_jwt_profile_grant   machine_pat_login   manipulate_user   oidc_session
otp_session   password_session   user_info
```

`make` in `benchmark/` has more targets than these; the extra ones are not published.

## Prerequisites

- `benchmark/.env` with `ZITADEL_HOST`, `ADMIN_LOGIN_NAME`, `ADMIN_PASSWORD` and `ADMIN_PAT`.
  `ADMIN_PAT` belongs to a machine user with IAM owner viewer / org manager / user manager and is
  what every cleanup path uses. `.env` is gitignored — keep it that way.
- `k6`, `jq`, `curl`, `python3`, and DuckDB (vendored at `benchmark/tools/.bin/duckdb`).
- Disk. A 30-minute run at 600 VUs writes a **2–3 GB** CSV; eleven of them will not fit. The
  tooling converts and deletes each CSV as it lands. Do not disable that.

## Safety

**This runbook targets a shared QA instance, not localhost.** That is the exception to the
skills-wide "local only" rule, and it is why the destructive rules here are absolute rather than
advisory. That instance also hosts real working organizations.

- **Only ever delete organizations named exactly `load-test-<ISO8601 timestamp>`**
  (e.g. `load-test-2026-08-26T20:11:56.440Z`). Match on an anchored regex, never a `startswith`.
  Re-check the name immediately before each DELETE. Never delete an org that does not follow the
  convention, whatever else is true.
- **Never clean up while a run is in flight.** The live run's own org matches the pattern too.
  Both cleanup scripts refuse when a k6 process is running.
- Cleanup scripts are **dry-run by default** and need `--confirm`.

### The `pgrep -f` self-match trap

An in-flight guard written as `pgrep -f k6` matches **your own shell's command line** whenever the
pattern appears in the command you are running — including inside a heredoc. This produced a
cleanup that refused to run with no k6 anywhere on the box, and a `pkill -f` that killed the wrong
thing. Match the process name (`pgrep -x k6`) or use a PID.

## Run

From `benchmark/`:

| Command | Purpose |
| --- | --- |
| `./run-bench.sh <target> <vus> <duration>` | One target. Tees k6's terminal output — including the `TOTAL RESULTS` block the docs page embeds — to `summaries/<target>_<stamp>.log`, and stamps `start_utc=` / `end_utc=` lines around the run. Runs the orphan-org preflight first. Skips immediately if a `.hold` file exists. |
| `./run-11.sh` | The full sweep: 11 targets, 600 VUs, 1800s each, published as each one lands. ~6.5h. |
| `./retry-failed.sh <sweep.log>` | Re-runs every target the sweep logged as failed. |

Create `.hold` to stop a sweep after the current target finishes without killing the running
test — far better than killing it, which leaves an orphan org and a multi-GB partial CSV.

**Run long sweeps detached and logged to files, never only in a session's scrollback.** A lost
session cost a run mid-flight once; the logs under `summaries/` are what made it recoverable.

## Publish

`./postprocess.sh <target> [version]` does the whole chain and then deletes the CSV:

1. `tools/csv2json.sh <in.csv> <out.json>` — DuckDB, per-second `quantile_cont` p50/p95/p99 over
   the custom trend metrics. **Setup rows carry `group='::setup'`, so filtering on an empty
   `group` is what isolates the test phase.** This is the step nobody remembered: the conversion
   is DuckDB, and it is not documented anywhere else in the repo.
2. `tools/status_breakdown.sh` — captures the HTTP status/error breakdown *while the CSV still
   exists*. Without this you cannot diagnose a failure rate after the CSV is gone.
3. `tools/gen_report.py <target> <log> <output.json> <version>` — renders the page and mirrors
   `output.json`.

The page is fully generated, including the deployment-spec rows, which are **hardcoded in
`gen_report.py`** — if the deployment changes, change them there. Republishing a page from its own
`output.json` is safe and idempotent.

### Time window

`Test start` and `Test end` carry the full date and time in UTC, derived from the first and last
per-second row in `output.json`. That per-second data is the only absolute record of when a test
ran, so a page must never report a bare `HH:MM` — without the date the metrics cannot be traced
back in monitoring the next day. (Pages before v4.17.1 have the bare form.)

### Rows only a human can fill

`gen_report.py` leaves these blank on purpose. Ask, and leave them empty rather than guessing —
an empty cell is honest, a copied one from the previous version is not:

- ZITADEL metrics during test, Database metrics during test
- Top 3 most expensive database queries
- ZITADEL feature flags
- flowchart outcome

## Metric cardinality

k6's `url` system tag is one time series per distinct URL. Targets that put ids in paths produced
**800k+ series** against k6's suggested limit of 100k, driving k6 to 6.8 GB and an OOM that took
down a run.

The fix is at the call site, not in the environment: requests whose path carries a user/session/
project id set a grouped `name` tag (`/v2beta/users/{userId}`). With those in place `name` stays
low-cardinality and worth keeping, while the raw `url` tag stays out via `K6_SYSTEM_TAGS` in
`run-bench.sh`. Watch k6's RSS during a run; if it climbs past a couple of GB, suspect a new
id-bearing URL without a grouped tag.

## Failure modes

Learn to tell these apart before touching anything.

### Orphan org from a failed run — `409 User already exists`

A failed run leaves its org and users behind. Their machine users (`zitachine-N`) collide
instance-wide with the next run's setup. `run-bench.sh` runs a preflight that deletes stale
benchmark orgs and waits for deletion (which is asynchronous) before starting. If you must remove
one while another run is legitimately in flight, use `./delete-one-org.sh <id> <exact-name>`.

### Read-after-write race in user creation

`createHuman`/`createMachine` POST the user, then GET `/v2/users/{id}` and read `res.json('user')`.
When the read beats the user projection, the field is missing. This originally resolved
`undefined` and surfaced much later as `TypeError: Cannot read property '0' of undefined` on
`loginNames[0]`, killing setup ~45s in — intermittent and load-dependent, and it looked exactly
like the collision above but was not. All the users existed; only the read-back failed.

`readUserAfterCreate()` now retries up to 10 times with backoff and rejects with the user id and
username instead of resolving `undefined`.

**Known remaining hole:** that loop calls `res.json('user')` unguarded, so a GET that *times out*
(null body) throws `GoError: the body is null so we can't transform it to JSON` instead of
consuming a retry. Guard the body and status before parsing, and keep looping on a timeout.

This is also a real signal, not only a harness flaw: projection lag under load is Zitadel
behaviour worth reporting.

### Check names that collide

k6 aggregates checks **by name**. Three library functions once registered checks under the same
name `update user is status ok`, two asserting `status === 201` against endpoints returning 200.
The result was a permanent, meaningless 28% failure rate on a published page — the most alarming
number on it, and entirely fictional. Give every check a distinct name, assert `2xx` rather than
an exact code, and do not duplicate in the use case a check the library already performs.

If a check failure rate looks structural (a suspiciously round split, or exactly one check's worth
of failures per iteration), suspect the harness before the server.

### Org deletion fails with 500 after heavy user churn

Deleting an org sends every child object to the pusher in one call, which exceeds the maximum
argument size once the org has churned enough users. Use
`./purge-org.sh <org-id> <exact-name> --confirm`: it deletes users a page (1000) at a time, then
the org. An org that had been returning a 16s 500 deleted in 1.3s afterwards.

Two things that loop forever if you do not handle them:

- **Ghost users.** The query side still lists a user the command side returns 404 for — the
  projection is stale. Treat 404 as success.
- Stop when a page removes nothing, not only when the count reaches zero.

## Reporting

Publish what happened, including the parts that make the harness look bad. When you correct a
generated row by hand — as with the fictional 28% above — say on the page that you did, and why.

Compare against the goals and the previous version's pages for the same target; that like-for-like
history is the whole point of publishing by version. Zitadel tracks throughput goals in GitHub
issues (e.g. #8352 for machine JWT profile grants, #4424 for logins) — state plainly whether each
is met.

Distinguish, always: **a rerun is needed when the numbers are untrustworthy**, not merely when a
run was ugly. A setup crash that aborts before the test phase costs nothing — the completed runs
around it are still valid. A cardinality bug that had k6 swapping during the measurement window
invalidates that target's numbers, and it must be run again.
