#!/usr/bin/env node
// Summarizes `pnpm audit --json` output grouped by vulnerable package, showing
// the installed versions, the fixed range and the dependency paths that pull it in.
//
// Usage (from the repository root):
//   pnpm audit --json > /tmp/audit.json; node skills/npm-dependency-audit-summary.mjs /tmp/audit.json
//   pnpm audit --json | node skills/npm-dependency-audit-summary.mjs
import { readFileSync } from "node:fs";

const input = readFileSync(process.argv[2] ?? 0, "utf8");
const { advisories = {}, metadata = {} } = JSON.parse(input);

const byPackage = new Map();
for (const advisory of Object.values(advisories)) {
  const entry = byPackage.get(advisory.module_name) ?? { severities: new Set(), patched: new Set(), versions: new Set(), paths: new Set() };
  entry.severities.add(advisory.severity);
  entry.patched.add(advisory.patched_versions);
  for (const finding of advisory.findings) {
    entry.versions.add(finding.version);
    finding.paths.forEach((path) => entry.paths.add(path));
  }
  byPackage.set(advisory.module_name, entry);
}

const MAX_PATHS = 5;
for (const [name, entry] of [...byPackage].sort(([a], [b]) => a.localeCompare(b))) {
  console.log(`${name}  installed: ${[...entry.versions].join(", ")}  fixed: ${[...entry.patched].join(" and ")}  severity: ${[...entry.severities].join("/")}`);
  const paths = [...entry.paths].sort();
  // The first segment is the workspace project (`apps__login` = apps/login, `.` = root),
  // the second is the dependency declared in that project's package.json.
  paths.slice(0, MAX_PATHS).forEach((path) => console.log(`    ${path}`));
  if (paths.length > MAX_PATHS) console.log(`    ... ${paths.length - MAX_PATHS} more`);
}

const counts = metadata.vulnerabilities ?? {};
console.log(`\n${byPackage.size} packages affected`, counts);
