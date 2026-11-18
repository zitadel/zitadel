#!/usr/bin/env node
// Fails the build when a published Server Action lives outside src/lib/actions/, or when the set of
// published Server Actions differs from server-actions.allowlist.json.
//
// Next.js publishes every export of a "use server" module that the build graph reaches as a public,
// unauthenticated POST endpoint, whether or not a client component calls it. Internal helpers that
// trust their arguments must therefore never end up in that set. Any change to it has to be a
// deliberate, reviewed edit of the allowlist.
//
// Usage: node checks/server-actions.mjs [path/to/server-reference-manifest.json]

import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const manifestPath = process.argv[2] ?? join(here, "..", ".next", "server", "server-reference-manifest.json");
const allowlistPath = join(here, "server-actions.allowlist.json");
const ACTIONS_DIR = "src/lib/actions/";

function normalizeFilename(filename) {
  // the manifest stores paths relative to the workspace root (e.g. apps/login/src/lib/actions/flow.ts)
  const index = filename.indexOf("src/");
  return index === -1 ? filename : filename.slice(index);
}

function published(manifest) {
  const actions = new Set();
  for (const runtime of ["node", "edge"]) {
    for (const entry of Object.values(manifest[runtime] ?? {})) {
      actions.add(`${normalizeFilename(entry.filename)}#${entry.exportedName}`);
    }
  }
  return actions;
}

function allowed(allowlist) {
  const actions = new Set();
  for (const [filename, names] of Object.entries(allowlist)) {
    for (const name of names) {
      actions.add(`${filename}#${name}`);
    }
  }
  return actions;
}

const actual = published(JSON.parse(readFileSync(manifestPath, "utf8")));
const expected = allowed(JSON.parse(readFileSync(allowlistPath, "utf8")));

const outsideActionsDir = [...actual].filter((action) => !action.startsWith(ACTIONS_DIR)).sort();
const unexpected = [...actual].filter((action) => !expected.has(action)).sort();
const missing = [...expected].filter((action) => !actual.has(action)).sort();

if (outsideActionsDir.length || unexpected.length || missing.length) {
  if (outsideActionsDir.length) {
    console.error(`These Server Actions are published from outside ${ACTIONS_DIR}:`);
    outsideActionsDir.forEach((action) => console.error(`  ! ${action}`));
    console.error(
      `Only modules in ${ACTIONS_DIR} may use "use server". Move the implementation to src/lib/server/ (import "server-only").`,
    );
  }
  if (unexpected.length) {
    console.error("These exports are published as Server Actions but are not in checks/server-actions.allowlist.json:");
    unexpected.forEach((action) => console.error(`  + ${action}`));
    console.error(
      'Every argument of a Server Action is attacker-controlled. Move internal helpers to a module without "use server"',
      '(add import "server-only"), or add the action to the allowlist after reviewing that it derives identity and policy on the server.',
    );
  }
  if (missing.length) {
    console.error("These allowlisted Server Actions are no longer published, remove them from the allowlist:");
    missing.forEach((action) => console.error(`  - ${action}`));
  }
  process.exit(1);
}

console.log(`Server Actions match the allowlist (${actual.size} published).`);
