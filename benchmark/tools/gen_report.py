#!/usr/bin/env python3
"""Generate a benchmark docs page from a k6 run.

Usage: gen_report.py <target> <summary.log> <output.json> <version>
Writes apps/docs/content/apis/benchmarks/<version>/<target>/index.mdx
(and mirrors output.json into apps/docs/src/data/benchmarks/... like every
existing version does).
"""
import json, re, sys, pathlib, shutil

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

# --- failed checks -> "Observed errors" ---
cf = re.search(r'^\s*checks_failed\.+:\s*([\d.]+)%\s+(\d+) out of (\d+)', summary, re.M)
failed_names = re.findall(r'^\s*✗ (.+?)\n\s*↳\s+(\d+)% — ✓ ([\d,]+) / ✗ ([\d,]+)', summary, re.M)
if cf and int(cf.group(2)) > 0:
    detail = ', '.join(f'{n.strip()} ({f} failed)' for n, _, _, f in failed_names) or 'see k6 output'
    observed = f'{cf.group(2)} of {cf.group(3)} checks failed ({cf.group(1)}%): {detail}'
else:
    observed = 'none'

# --- test window from the converted per-second data ---
rows = json.loads(pathlib.Path(out_json).read_text())
start_ts, end_ts = rows[0]['timestamp'], rows[-1]['timestamp']
start_hm = start_ts[11:16]
run = re.search(r'running \((\d+)m([\d.]+)s\)', summary)
duration = f'{run.group(1)}min' if run else ''

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
| Test start                            | {start_hm} UTC                                                                      |
| Test duration                         | {duration}                                                                          |
| Executed test                         | {esc(target)}                                                                    |
| k6 version                            | v2.1.0                                                                           |
| VUs                                   | 600                                                                              |
| Client location                       | US1                                                                              |
| ZITADEL location                      | US1                                                                              |
| ZITADEL container specification       | vCPU: 6<br/> Memory: 6 Gi <br/>Container min scale: 7<br/>Container max scale: 7 |
| ZITADEL Version                       | {version}                                                                           |
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
shutil.copyfile(out_json, content_dir / 'output.json')
shutil.copyfile(out_json, data_dir / 'output.json')
print(f'{target}: start={start_ts} end={end_ts} duration={duration} '
      f'iterations={iterations} ({per_sec:.0f}/s) rows={len(rows)}')
print(f'  observed errors: {observed}')
print(f'  wrote {content_dir}/index.mdx + output.json (mirrored to src/data)')
