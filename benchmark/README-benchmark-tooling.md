# Running and publishing a benchmark sweep

Helper scripts around the `make` targets in this directory. Credentials and the
target host come from `.env` (`ZITADEL_HOST`, `ADMIN_LOGIN_NAME`, `ADMIN_PASSWORD`,
`ADMIN_PAT`).

## Run

- `./run-bench.sh <target> <vus> <duration>` — one target. Tees k6's terminal
  output, including the `TOTAL RESULTS` block that the docs page embeds, to
  `summaries/<target>_<stamp>.log`. Skips immediately if a `.hold` file exists,
  which is how an in-flight sweep is stopped without killing the running test.
- `./run-11.sh` — the full publication sweep: the 11 targets that have docs pages,
  600 VUs / 1800s each, converted and published as each one lands. ~6.5h.
- `./retry-failed.sh <sweep.log>` — re-runs every target the sweep logged as failed.

## Publish

- `./postprocess.sh <target> [version]` — converts the newest CSV in `output/`,
  writes the docs page, then deletes the CSV (a 30min run is 2-3 GB, and 11 of
  them will not fit on the runner).
- `tools/csv2json.sh <in.csv> <out.json>` — DuckDB per-second p50/p95/p99 over the
  test-phase custom trend metrics. Setup rows carry `group='::setup'`, so filtering
  on an empty group is what isolates the test phase.
- `tools/gen_report.py <target> <log> <output.json> <version>` — renders
  `apps/docs/content/apis/benchmarks/<version>/<target>/index.mdx` from the k6 log
  and mirrors `output.json` into `apps/docs/src/data/benchmarks/`.

Deployment facts in the results table (container/DB spec, locations, feature flags)
cannot be read from the host and must be filled in by hand. New pages also need an
entry under Benchmarks in `apps/docs/content/apis/_meta.json`.

## Cleanup

A failed run leaves an orphan org behind, which can break later runs.

- `./cleanup-stale-orgs.sh [--confirm]` — dry-run by default. Refuses while a k6 run
  is in flight, because the live run's own org matches the naming convention.
- `./delete-one-org.sh <id> <exact-name>` — removes a single orphan while another
  run is legitimately running.

Both only ever touch organizations named exactly `load-test-<ISO timestamp>`
(e.g. `load-test-2026-08-26T20:11:56.440Z`). Never delete an organization that does
not follow that convention.
