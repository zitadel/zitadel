import { mkdtemp, mkdir, readFile, rm, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { afterEach, describe, expect, test } from "vitest";
import {
  detectPackageManager,
  mergeEnvExample,
  runAddCommand,
} from "./add.js";

const noopLogger = { log: () => {}, warn: () => {} };

async function createTempProject(options: { typescript?: boolean } = {}): Promise<string> {
  const typescript = options.typescript ?? true;
  const dir = await mkdtemp(path.join(os.tmpdir(), "zitadel-nextjs-add-"));
  await writeFile(
    path.join(dir, "package.json"),
    JSON.stringify(
      {
        name: "test-app",
        private: true,
        dependencies: {
          next: "15.0.0",
        },
      },
      null,
      2,
    ),
    "utf8",
  );
  if (typescript) {
    await writeFile(
      path.join(dir, "tsconfig.json"),
      JSON.stringify({ compilerOptions: {} }, null, 2),
      "utf8",
    );
  }
  await mkdir(path.join(dir, "app"), { recursive: true });
  return dir;
}

const cleanupDirs: string[] = [];

afterEach(async () => {
  await Promise.all(
    cleanupDirs.splice(0).map((dir) =>
      rm(dir, { recursive: true, force: true }),
    ),
  );
});

describe("add command helpers", () => {
  test("detectPackageManager prefers pnpm lockfile", async () => {
    const dir = await createTempProject();
    cleanupDirs.push(dir);
    await writeFile(path.join(dir, "pnpm-lock.yaml"), "lockfile", "utf8");

    expect(detectPackageManager(dir)).toBe("pnpm");
  });

  test("detectPackageManager resolves lockfile from parent directories", async () => {
    const root = await mkdtemp(path.join(os.tmpdir(), "zitadel-nextjs-add-parent-"));
    cleanupDirs.push(root);
    const projectDir = path.join(root, "apps", "web");
    await mkdir(path.join(projectDir, "app"), { recursive: true });
    await writeFile(
      path.join(projectDir, "package.json"),
      JSON.stringify(
        {
          name: "test-app",
          private: true,
          dependencies: { next: "15.0.0" },
        },
        null,
        2,
      ),
      "utf8",
    );
    await writeFile(path.join(root, "pnpm-lock.yaml"), "lockfile", "utf8");

    expect(detectPackageManager(projectDir)).toBe("pnpm");
  });

  test("mergeEnvExample is idempotent", () => {
    const once = mergeEnvExample("FOO=bar\n");
    const twice = mergeEnvExample(once);

    expect(once).toContain("ZITADEL_ISSUER_URL=");
    expect(once).toContain("ZITADEL_CALLBACK_URL=");
    expect(once).toContain("ZITADEL_POST_LOGIN_URL=");
    expect(twice).toBe(once);
  });
});

describe("runAddCommand", () => {
  test("creates route files and env example", async () => {
    const dir = await createTempProject();
    cleanupDirs.push(dir);

    const result = await runAddCommand({
      cwd: dir,
      skipInstall: true,
      logger: noopLogger,
    });

    expect(result.authMode).toBe("oidc");
    expect(result.dependencySource).toBe("npm");
    expect(result.withApi).toBe(false);
    expect(result.withWebhook).toBe(false);
    expect(result.createdFiles).toHaveLength(4);
    expect(result.skippedFiles).toHaveLength(0);
    expect(result.envUpdated).toBe(true);
    await expect(
      readFile(path.join(dir, "app/api/auth/signin/route.ts"), "utf8"),
    ).resolves.toContain('from "@zitadel/nextjs/auth/oidc"');
    await expect(
      readFile(path.join(dir, "app/api/auth/callback/route.ts"), "utf8"),
    ).resolves.toContain('?? "/auth"');
    await expect(
      readFile(path.join(dir, "app/auth/page.tsx"), "utf8"),
    ).resolves.toContain("Sign in with ZITADEL");
    await expect(
      readFile(path.join(dir, "app/auth/page.tsx"), "utf8"),
    ).resolves.toContain("idTokenClaims");
    await expect(readFile(path.join(dir, ".env.example"), "utf8")).resolves.toContain(
      "ZITADEL_ISSUER_URL=",
    );
    await expect(readFile(path.join(dir, ".env.example"), "utf8")).resolves.toContain(
      "ZITADEL_POST_LOGIN_URL=/auth",
    );
  });

  test("does not write files in dry-run mode", async () => {
    const dir = await createTempProject();
    cleanupDirs.push(dir);

    const result = await runAddCommand({
      cwd: dir,
      dryRun: true,
      skipInstall: true,
      logger: noopLogger,
    });

    expect(result.createdFiles).toHaveLength(4);
    await expect(
      readFile(path.join(dir, "app/api/auth/signin/route.ts"), "utf8"),
    ).rejects.toThrow();
    await expect(
      readFile(path.join(dir, "app/auth/page.tsx"), "utf8"),
    ).rejects.toThrow();
    await expect(readFile(path.join(dir, ".env.example"), "utf8")).rejects.toThrow();
  });

  test("skips existing files without overwriting", async () => {
    const dir = await createTempProject();
    cleanupDirs.push(dir);
    const existingFile = path.join(dir, "app/api/auth/signin/route.ts");
    await mkdir(path.dirname(existingFile), { recursive: true });
    await writeFile(existingFile, "custom-content", "utf8");

    const result = await runAddCommand({
      cwd: dir,
      skipInstall: true,
      logger: noopLogger,
    });

    expect(result.skippedFiles).toContain("app/api/auth/signin/route.ts");
    expect(result.createdFiles).toContain("app/auth/page.tsx");
    await expect(readFile(existingFile, "utf8")).resolves.toBe("custom-content");
  });

  test("creates JavaScript route files when project is not TypeScript", async () => {
    const dir = await createTempProject({ typescript: false });
    cleanupDirs.push(dir);

    const result = await runAddCommand({
      cwd: dir,
      skipInstall: true,
      logger: noopLogger,
    });

    expect(result.createdFiles).toContain("app/api/auth/signin/route.js");
    expect(result.createdFiles).toContain("app/auth/page.jsx");
    await expect(
      readFile(path.join(dir, "app/api/auth/signin/route.js"), "utf8"),
    ).resolves.toContain("signIn");
  });

  test("supports explicit session auth and optional api/webhook scaffolds", async () => {
    const dir = await createTempProject();
    cleanupDirs.push(dir);

    const result = await runAddCommand({
      cwd: dir,
      auth: "session",
      withApi: true,
      withWebhook: true,
      skipInstall: true,
      logger: noopLogger,
    });

    expect(result.authMode).toBe("session");
    expect(result.dependencySource).toBe("npm");
    expect(result.withApi).toBe(true);
    expect(result.withWebhook).toBe(true);
    expect(result.createdFiles).toEqual(
      expect.arrayContaining([
        "app/api/auth/session/create/route.ts",
        "app/api/auth/session/callback/route.ts",
        "app/api/zitadel/user/route.ts",
        "app/api/zitadel/events/route.ts",
      ]),
    );
    expect(result.createdFiles).not.toContain("app/auth/page.tsx");
    await expect(readFile(path.join(dir, ".env.example"), "utf8")).resolves.toContain(
      "ZITADEL_WEBHOOK_SECRET=",
    );
    await expect(readFile(path.join(dir, ".env.example"), "utf8")).resolves.toContain(
      "ZITADEL_API_URL=",
    );
    await expect(readFile(path.join(dir, ".env.example"), "utf8")).resolves.not.toContain(
      "ZITADEL_ISSUER_URL=",
    );
  });

  test("supports workspace dependency source", async () => {
    const dir = await createTempProject();
    cleanupDirs.push(dir);
    await writeFile(path.join(dir, "pnpm-lock.yaml"), "lockfile", "utf8");
    await writeFile(path.join(dir, "pnpm-workspace.yaml"), "packages:\n  - .\n", "utf8");
    const logs: string[] = [];

    const result = await runAddCommand({
      cwd: dir,
      source: "workspace",
      dryRun: true,
      logger: { log: (msg) => logs.push(msg), warn: () => {} },
    });

    expect(result.dependencySource).toBe("workspace");
    expect(logs.some((line) => line.includes("@zitadel/nextjs@workspace:*"))).toBe(
      true,
    );
    expect(logs.some((line) => line.includes("@zitadel/react@workspace:*"))).toBe(
      true,
    );
    expect(
      logs.some((line) => line.includes("@zitadel/zitadel-js@workspace:*")),
    ).toBe(true);
  });
});
