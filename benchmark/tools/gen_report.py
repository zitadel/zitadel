#!/usr/bin/env python3
"""Generate a benchmark docs page from a k6 run.

Usage: gen_report.py <target> <summary.log> <output.json> <version>
Writes apps/docs/content/apis/benchmarks/<version>/<target>/index.mdx
(and mirrors output.json into apps/docs/src/data/benchmarks/... like every
existing version does).
"""
import collections, json, re, sys, pathlib, shutil

ANSI = re.compile(r'\x1b\[[0-9;]*m')
DOCS = pathlib.Path(__file__).resolve().parents[2] / 'apps/docs'

target, log_path, out_json, version = sys.argv[1:5]
log = ANSI.sub('', pathlib.Path(log_path).read_text())

# --- k6 terminal summary: the whole TOTAL RESULTS block ---
m = re.search(r'( *█ TOTAL RESULTS.*?)\n\ndefault ✓ \[[^\]]*\][^\n]*\n', log, re.S)
if not m:
    m = re.search(r'( *█ TOTAL RESULTS.*?running \([^\)]*\),[^\n]*\ndefault [^\n]*)\n', log, re.S)
summary = m.group(1).rstrip() if m else log
summary = re.sub(r'\n{3,}', '\n\n', summary)

# --- iterations per second ---
it = re.search(r'^\s*iterations\.+:\s*(\d+)\s+([\d.]+)/s', summary, re.M)
iterations, per_sec = (it.group(1), float(it.group(2))) if it else ('', 0.0)

# --- "Observed errors" ---
# k6 counts a failed check and a failed request separately. A run can have zero
# failed checks and still have thousands of failed requests, so reporting only
# checks publishes "none" over real errors. Report both.
cf = re.search(r'^\s*checks_failed\.+:\s*([\d.]+)%\s+(\d+) out of (\d+)', summary, re.M)
failed_names = re.findall(r'^\s*✗ (.+?)\n\s*↳\s+(\d+)% — ✓ ([\d,]+) / ✗ ([\d,]+)', summary, re.M)
rf = re.search(r'^\s*http_req_failed\.+:\s*([\d.]+)%\s+(\d+) out of (\d+)', summary, re.M)

# Classify what k6 logged, so the row says what went wrong and not only how much.
# Counted over error/warning lines only, to avoid matching unrelated prose.
# Network-level causes, which partition the failed requests.
REQUEST_SIGNATURES = [
    ('request timeout', 'request timeout'),
    ('GOAWAY', 'connection closed by server (GOAWAY)'),
    ('connection reset by peer', 'connection reset'),
    ('read tcp', 'connection reset'),
    ('context deadline exceeded', 'context deadline exceeded'),
    ('no such host', 'DNS resolution failure'),
]
# Script-level fallout. These are consequences of the failures above, not extra
# failed requests, so they are reported separately and never summed with them.
SCRIPT_SIGNATURES = [
    ('the body is null', 'iteration{s} aborted on a null response body'),
]
err_lines = [l for l in log.splitlines()
             if 'level=error' in l or 'Request Failed' in l or 'level=warning' in l]
blob = '\n'.join(err_lines)

# Classify each logged request failure once, by its first matching signature.
# Counting signature occurrences across the whole blob instead would double-count
# any line two signatures both match, and can report more causes than failures.
counts = collections.Counter()
for line in err_lines:
    if 'Request Failed' not in line:
        continue
    for needle, label in REQUEST_SIGNATURES:
        if needle in line:
            counts[label] += 1
            break
    else:
        m = re.search(r'status: (\d{3})', line)
        counts[f'HTTP {m.group(1)}' if m else 'unclassified'] += 1
breakdown = [f'{n:,}x {label}' for label, n in counts.most_common()]

# k6 counts a request as failed whether or not it logged a line for it (a 503
# that fails expected_response logs nothing), so the classified causes can add
# up to less than the total. Say so rather than implying the list is complete.
if rf:
    unexplained = int(rf.group(2)) - sum(counts.values())
    if unexplained > 0:
        breakdown.append(f'{unexplained:,}x not individually logged (see k6 output)'
                         if breakdown else
                         f'none of the {unexplained:,} individually logged (see k6 output)')

script_errs = [f'{blob.count(needle):,} ' + label.format(s='' if blob.count(needle) == 1 else 's')
               for needle, label in SCRIPT_SIGNATURES if blob.count(needle)]

parts = []
if rf and int(rf.group(2)) > 0:
    part = f'{int(rf.group(2)):,} of {int(rf.group(3)):,} requests failed ({rf.group(1)}%)'
    if breakdown:
        part += ': ' + ', '.join(breakdown)
    parts.append(part)
elif breakdown:
    parts.append('logged errors: ' + ', '.join(breakdown))

if script_errs:
    parts.append('; '.join(script_errs))

if cf and int(cf.group(2)) > 0:
    detail = ', '.join(f'{n.strip()} ({f} failed)' for n, _, _, f in failed_names) or 'see k6 output'
    parts.append(f'{int(cf.group(2)):,} of {int(cf.group(3)):,} checks failed ({cf.group(1)}%): {detail}')

observed = '. '.join(parts) if parts else 'none'

# --- test window from the converted per-second data ---
rows = json.loads(pathlib.Path(out_json).read_text())
# The per-second rows are the only absolute record of when the test ran, so the
# page reports the full window rather than just a time of day: without the date
# the metrics cannot be traced back in monitoring the next day.
timestamps = sorted(r['timestamp'] for r in rows)
start_ts, end_ts = timestamps[0], timestamps[-1]
fmt = lambda ts: f'{ts[:10]} {ts[11:19]} UTC'
start_pretty, end_pretty = fmt(start_ts), fmt(end_ts)
run = re.search(r'running \((\d+)m([\d.]+)s\)', summary)
duration = f'{run.group(1)}min' if run else ''

# VU count comes from the run-bench.sh stamp rather than a constant: targets are
# not all run at the same concurrency, and a hardcoded 600 would silently
# misreport any run that deviates.
vu = re.search(r'^start_utc=.*\bvus=(\d+)', log, re.M)
vus = vu.group(1) if vu else '600'

esc = lambda s: s.replace('_', r'\_')
pretty = target.replace('_', ' ')

mdx = f'''---
title: {pretty} benchmark of zitadel {version}
sidebar_label: {pretty}
---

import OutputSource from "./output.json";
import {{ BenchmarkChart }} from '@/components/benchmark_chart';

Benchmark results of the {version} release of Zitadel.

## Performance test results

| Metric                                | Value                                                                            |
|:--------------------------------------|:---------------------------------------------------------------------------------|
| Baseline                              | none                                                                             |
| Purpose                               | Test current performance                                                         |
| Test start                            | {start_pretty}                                                                      |
| Test end                              | {end_pretty}                                                                        |
| Test duration                         | {duration}                                                                          |
| Executed test                         | {esc(target)}                                                                    |
| k6 version                            | v2.1.0                                                                           |
| VUs                                   | {vus}                                                                              |
| Client location                       | US1                                                                              |
| ZITADEL location                      | US1                                                                              |
| ZITADEL container specification       | vCPU: 6<br/> Memory: 6 Gi <br/>Container min scale: 7<br/>Container max scale: 7 |
| ZITADEL Version                       | {version}                                                                           |
| ZITADEL Settings                      | Eventstore autovacuum (new in v4.17): Enabled: true<br/>VacuumThreshold: 1000000<br/>AnalyzeThreshold: 1000000<br/>Tuned for the ~4,000 events/s expected during the benchmark, per the [tuning guide](https://zitadel.com/docs/self-hosting/manage/tuning#reading-events). |
| ZITADEL feature flags                 |                                                                                  |
| Database                              | type: psql<br />version: v17.4                                                   |
| Database location                     | US1                                                                              |
| Database specification                | vCPU: 8<br/> memory: 32Gib                                                       |
| ZITADEL metrics during test           |                                                                                  |
| Observed errors                       | {observed}                                                                          |
| Top 3 most expensive database queries |                                                                                  |
| Database metrics during test          |                                                                                  |
| k6 Iterations per second              | {per_sec:.0f}                                                                       |
| k6 output                             | [output](#k6-output)                                                             |
| flowchart outcome                     |                                                                                  |

## Endpoint latencies

<BenchmarkChart testResults={{OutputSource}} />

## k6 output
```bash
{summary}
```
'''

content_dir = DOCS / 'content/apis/benchmarks' / version / target
data_dir = DOCS / 'src/data/benchmarks' / version / target
content_dir.mkdir(parents=True, exist_ok=True)
data_dir.mkdir(parents=True, exist_ok=True)
(content_dir / 'index.mdx').write_text(mdx)
# republishing a page from its own output.json is a legitimate no-op copy
for dest in (content_dir / 'output.json', data_dir / 'output.json'):
    if pathlib.Path(out_json).resolve() != dest.resolve():
        shutil.copyfile(out_json, dest)
print(f'{target}: start={start_ts} end={end_ts} duration={duration} '
      f'iterations={iterations} ({per_sec:.0f}/s) rows={len(rows)}')
print(f'  observed errors: {observed}')
print(f'  wrote {content_dir}/index.mdx + output.json (mirrored to src/data)')
