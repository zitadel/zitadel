#!/usr/bin/env node

import { type AuthMode, runAddCommand } from "./add.js";

function printHelp() {
  console.log(`zitadel-nextjs <command>

Commands:
  add                 Add @zitadel/nextjs bootstrap files to an existing Next.js App Router project

Options:
  --auth <mode>       Auth scaffold mode: oidc (default) or session
  --with-api          Add a sample authenticated ZITADEL API route
  --with-webhook      Add a sample Actions v2 webhook route
  --with-events       Alias for --with-webhook
  --cwd <path>        Project path (defaults to current working directory)
  --dry-run           Print planned changes without writing files
  --skip-install      Do not run package manager install command
  --yes               Reserved for non-interactive mode compatibility
  -h, --help          Show this help message
`);
}

function parseAddArgs(args: string[]) {
  const options: {
    cwd?: string;
    dryRun?: boolean;
    skipInstall?: boolean;
    auth?: AuthMode;
    withApi?: boolean;
    withWebhook?: boolean;
  } = {};

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === "--cwd") {
      const nextValue = args[i + 1];
      if (!nextValue || nextValue.startsWith("-")) {
        throw new Error("Missing value for --cwd");
      }
      options.cwd = nextValue;
      i++;
      continue;
    }
    if (arg === "--dry-run") {
      options.dryRun = true;
      continue;
    }
    if (arg === "--auth") {
      const nextValue = args[i + 1];
      if (!nextValue || nextValue.startsWith("-")) {
        throw new Error("Missing value for --auth");
      }
      if (nextValue !== "oidc" && nextValue !== "session") {
        throw new Error("Invalid --auth value. Expected: oidc or session.");
      }
      options.auth = nextValue;
      i++;
      continue;
    }
    if (arg === "--with-api") {
      options.withApi = true;
      continue;
    }
    if (arg === "--with-webhook" || arg === "--with-events") {
      options.withWebhook = true;
      continue;
    }
    if (arg === "--skip-install") {
      options.skipInstall = true;
      continue;
    }
    if (arg === "--yes") {
      continue;
    }
    if (arg === "-h" || arg === "--help") {
      printHelp();
      process.exit(0);
    }
    throw new Error(`Unknown option: ${arg}`);
  }

  return options;
}

async function main() {
  const args = process.argv.slice(2);
  const command = args[0];

  if (!command || command === "-h" || command === "--help") {
    printHelp();
    return;
  }

  if (command !== "add") {
    throw new Error(`Unknown command: ${command}`);
  }

  const result = await runAddCommand(parseAddArgs(args.slice(1)));
  console.log("");
  console.log("ZITADEL Next.js bootstrap complete.");
  console.log(`Project: ${result.projectRoot}`);
  console.log(`App dir: ${result.appDirectory}`);
  console.log(`Package manager: ${result.packageManager}`);
  console.log(`Auth scaffold: ${result.authMode}`);
  console.log(`Include API sample: ${result.withApi ? "yes" : "no"}`);
  console.log(`Include webhook sample: ${result.withWebhook ? "yes" : "no"}`);
  if (result.createdFiles.length) {
    console.log(`Created files (${result.createdFiles.length}):`);
    result.createdFiles.forEach((file) => console.log(`  - ${file}`));
  } else {
    console.log("No new route files created.");
  }
  if (result.skippedFiles.length) {
    console.log(`Skipped existing files (${result.skippedFiles.length}):`);
    result.skippedFiles.forEach((file) => console.log(`  - ${file}`));
  }
  console.log(
    `Environment file updated: ${result.envUpdated ? "yes" : "no"}`,
  );
  if (!result.dependencyInstalled) {
    console.log(
      "Dependency install skipped. Install SDK packages before running your app.",
    );
  }
  if (result.authMode === "oidc") {
    console.log("Test page: /auth");
    console.log("Next: set ZITADEL_* values in .env.local and open that page.");
  }
}

main().catch((error) => {
  console.error(error instanceof Error ? error.message : String(error));
  process.exit(1);
});
