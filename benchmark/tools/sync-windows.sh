#!/usr/bin/env bash
# Rewrite the embedded WINDOWS block in the fetch scripts from the published pages.
# Run after every republish; the fetch scripts must work standalone in Cloud Shell,
# where this repository is not checked out, so the windows have to be baked in.
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
W="$("${HERE}/gen-windows.sh" "$@")"
for f in "${HERE}/fetch-gcp-metrics.sh" "${HERE}/fetch-gcp-queries.sh"; do
  WINDOWS="$W" python3 - "$f" <<'PY'
import os,re,sys
p=sys.argv[1]; s=open(p).read()
new='WINDOWS="${WINDOWS_DEFAULT:-\n'+os.environ['WINDOWS']+'\n}"\n'
# Count substitutions, not content changes: a correct re-run is a no-op and must not fail.
s2,n=re.subn(r'WINDOWS="\$\{WINDOWS_DEFAULT:-\n.*?\n\}"\n', new, s, count=1, flags=re.S)
if n==0:
    s2,n=re.subn(r'WINDOWS="\n.*?\n"\n', new, s, count=1, flags=re.S)
assert n==1, f"{p}: no WINDOWS block matched"
if s2!=s:
    open(p,'w').write(s2)
    print("  updated", file=sys.stderr)
PY
  echo "synced $(basename "$f")"
done
