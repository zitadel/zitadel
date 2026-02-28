import { spawnSync } from "node:child_process";
import { existsSync } from "node:fs";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import path from "node:path";

const SDK_PACKAGE = "@zitadel/nextjs";
const WORKSPACE_SDK_PACKAGES = [
  "@zitadel/zitadel-js@workspace:*",
  "@zitadel/react@workspace:*",
  "@zitadel/nextjs@workspace:*",
] as const;

const OIDC_ROUTE_TEMPLATES: Record<string, string> = {
  "api/auth/signin/route": `import { signIn } from "@zitadel/nextjs/auth/oidc";

export async function GET() {
  await signIn();
}
`,
  "api/auth/callback/route": `import { handleCallback } from "@zitadel/nextjs/auth/oidc";

export async function GET(request: Request) {
  await handleCallback(request);
  const postLoginPath = process.env.ZITADEL_POST_LOGIN_URL ?? "/auth";
  return Response.redirect(new URL(postLoginPath, request.url));
}
`,
  "api/auth/signout/route": `import { signOut } from "@zitadel/nextjs/auth/oidc";

export async function GET(request: Request) {
  const postLogoutRedirectUri = process.env.ZITADEL_POST_LOGOUT_URL
    ? new URL(process.env.ZITADEL_POST_LOGOUT_URL, request.url).toString()
    : undefined;

  await signOut({ postLogoutRedirectUri });

  const fallbackPath = process.env.ZITADEL_POST_LOGOUT_URL ?? "/auth";
  return Response.redirect(new URL(fallbackPath, request.url));
}
`,
};

const OIDC_TEST_PAGE_TEMPLATES: Record<string, string> = {
  "auth/page": `import { getSession } from "@zitadel/nextjs";

function decodeIdTokenClaims(idToken) {
  if (!idToken) {
    return null;
  }

  const [, payload] = idToken.split(".");
  if (!payload) {
    return null;
  }

  try {
    const json = Buffer.from(payload, "base64url").toString("utf8");
    return JSON.parse(json);
  } catch {
    return null;
  }
}

export default async function AuthPage() {
  const session = await getSession();
  const idTokenClaims = decodeIdTokenClaims(session?.idToken);

  const result = session
    ? {
        authenticated: true,
        expiresAt: session.expiresAt,
        expiresAtIso: new Date(session.expiresAt * 1000).toISOString(),
        hasIdToken: Boolean(session.idToken),
        hasRefreshToken: Boolean(session.refreshToken),
        idTokenClaims,
      }
    : {
        authenticated: false,
      };

  return (
    <main style={{ padding: "2rem", fontFamily: "sans-serif" }}>
      <h1>ZITADEL OIDC test page</h1>
      <p>Use this page to verify that the redirect login flow works.</p>
      <ul>
        <li>
          <a href="/api/auth/signin">Sign in with ZITADEL</a>
        </li>
        <li>
          <a href="/api/auth/signout">Sign out</a>
        </li>
      </ul>
      <h2>OIDC result</h2>
      <pre>{JSON.stringify(result, null, 2)}</pre>
    </main>
  );
}
`,
};

const SESSION_ROUTE_TEMPLATES: Record<string, string> = {
  "api/auth/session/create/route": `import { createSession } from "@zitadel/nextjs/auth/session";

export async function POST(request: Request) {
  const body = (await request.json()) as Parameters<typeof createSession>[0];
  const session = await createSession(body);
  return Response.json(session);
}
`,
  "api/auth/session/callback/route": `import { createCallback } from "@zitadel/nextjs/auth/session";

export async function POST(request: Request) {
  const body = (await request.json()) as Parameters<typeof createCallback>[0];
  const callbackUrl = await createCallback(body);
  return Response.json({ callbackUrl });
}
`,
};

const API_ROUTE_TEMPLATES: Record<string, string> = {
  "api/zitadel/user/route": `import { createZitadelApiClient } from "@zitadel/nextjs/api";

export async function GET(request: Request) {
  const userId = new URL(request.url).searchParams.get("userId");
  if (!userId) {
    return Response.json({ error: "Missing required query parameter: userId" }, { status: 400 });
  }

  const api = await createZitadelApiClient();
  const user = await api.userService.getUser({ userId });
  return Response.json(user);
}
`,
};

const WEBHOOK_ROUTE_TEMPLATES: Record<string, string> = {
  "api/zitadel/events/route": `import { createZitadelWebhookHandler } from "@zitadel/nextjs/webhook";

export const POST = createZitadelWebhookHandler({
  payloadType: "json",
  onEvent: async (_event) => {},
});
`,
};

const OIDC_ENV_VARS: Array<[string, string]> = [
  ["ZITADEL_ISSUER_URL", "https://my-instance.zitadel.cloud"],
  ["ZITADEL_CLIENT_ID", ""],
  ["ZITADEL_CALLBACK_URL", "http://localhost:3000/api/auth/callback"],
  ["ZITADEL_COOKIE_SECRET", ""],
  ["ZITADEL_POST_LOGIN_URL", "/auth"],
  ["ZITADEL_POST_LOGOUT_URL", "/auth"],
];

const SESSION_ENV_VARS: Array<[string, string]> = [
  ["ZITADEL_API_URL", "https://my-instance.zitadel.cloud"],
  ["ZITADEL_SERVICE_USER_TOKEN", ""],
];

const API_ENV_VARS: Array<[string, string]> = [
  ["ZITADEL_API_URL", "https://my-instance.zitadel.cloud"],
  ["ZITADEL_SERVICE_USER_TOKEN", ""],
  ["ZITADEL_SERVICE_USER_KEY_ID", ""],
  ["ZITADEL_SERVICE_USER_ID", ""],
  ["ZITADEL_SERVICE_USER_PRIVATE_KEY", ""],
];

const WEBHOOK_ENV_VARS: Array<[string, string]> = [
  ["ZITADEL_WEBHOOK_PAYLOAD_TYPE", "json"],
  ["ZITADEL_WEBHOOK_SECRET", ""],
  ["ZITADEL_WEBHOOK_JWKS_ENDPOINT", ""],
  ["ZITADEL_WEBHOOK_JWE_PRIVATE_KEY", ""],
];

type Logger = Pick<Console, "log" | "warn">;

export type PackageManager = "pnpm" | "npm" | "yarn";
export type AuthMode = "oidc" | "session";
export type DependencySource = "npm" | "workspace";

export interface AddCommandOptions {
  cwd?: string;
  dryRun?: boolean;
  skipInstall?: boolean;
  source?: DependencySource;
  auth?: AuthMode;
  withApi?: boolean;
  withWebhook?: boolean;
  logger?: Logger;
}

export interface AddCommandResult {
  projectRoot: string;
  appDirectory: string;
  packageManager: PackageManager;
  authMode: AuthMode;
  withApi: boolean;
  withWebhook: boolean;
  dependencySource: DependencySource;
  createdFiles: string[];
  skippedFiles: string[];
  envUpdated: boolean;
  dependencyInstalled: boolean;
}

interface ProjectPackageJson {
  dependencies?: Record<string, string>;
  devDependencies?: Record<string, string>;
}

function isTypeScriptProject(
  projectRoot: string,
  pkg: ProjectPackageJson,
): boolean {
  if (existsSync(path.join(projectRoot, "tsconfig.json"))) {
    return true;
  }
  return Boolean(pkg.dependencies?.typescript || pkg.devDependencies?.typescript);
}

export function detectPackageManager(projectRoot: string): PackageManager {
  let currentDir = projectRoot;
  let previousDir = "";
  while (currentDir !== previousDir) {
    if (existsSync(path.join(currentDir, "pnpm-lock.yaml"))) {
      return "pnpm";
    }
    if (existsSync(path.join(currentDir, "yarn.lock"))) {
      return "yarn";
    }
    if (existsSync(path.join(currentDir, "package-lock.json"))) {
      return "npm";
    }
    previousDir = currentDir;
    currentDir = path.dirname(currentDir);
  }
  return "npm";
}

function dedupeEnvVars(envVars: Array<[string, string]>): Array<[string, string]> {
  const map = new Map<string, string>();
  for (const [key, value] of envVars) {
    if (!map.has(key)) {
      map.set(key, value);
    }
  }
  return Array.from(map.entries());
}

function findPnpmWorkspaceRoot(projectRoot: string): string | null {
  let currentDir = projectRoot;
  let previousDir = "";
  while (currentDir !== previousDir) {
    if (existsSync(path.join(currentDir, "pnpm-workspace.yaml"))) {
      return currentDir;
    }
    previousDir = currentDir;
    currentDir = path.dirname(currentDir);
  }
  return null;
}

function resolveTemplates(options: {
  authMode: AuthMode;
  withApi: boolean;
  withWebhook: boolean;
}): Record<string, string> {
  const templates: Record<string, string> = {};
  if (options.authMode === "session") {
    Object.assign(templates, SESSION_ROUTE_TEMPLATES);
  } else {
    Object.assign(templates, OIDC_ROUTE_TEMPLATES, OIDC_TEST_PAGE_TEMPLATES);
  }
  if (options.withApi) {
    Object.assign(templates, API_ROUTE_TEMPLATES);
  }
  if (options.withWebhook) {
    Object.assign(templates, WEBHOOK_ROUTE_TEMPLATES);
  }
  return templates;
}

function resolveEnvVars(options: {
  authMode: AuthMode;
  withApi: boolean;
  withWebhook: boolean;
}): Array<[string, string]> {
  const vars: Array<[string, string]> =
    options.authMode === "session" ? [...SESSION_ENV_VARS] : [...OIDC_ENV_VARS];
  if (options.withApi) {
    vars.push(...API_ENV_VARS);
  }
  if (options.withWebhook) {
    vars.push(...WEBHOOK_ENV_VARS);
  }
  return dedupeEnvVars(vars);
}

export function mergeEnvExample(
  content: string,
  requiredEnvVars: Array<[string, string]> = OIDC_ENV_VARS,
): string {
  const missing = requiredEnvVars.filter(
    ([key]) => !new RegExp(`^${key}=`, "m").test(content),
  );
  if (!missing.length) {
    return content;
  }

  const lines = [
    "# ZITADEL Next.js SDK bootstrap",
    ...missing.map(([key, value]) => `${key}=${value}`),
  ];

  const prefix = content.trimEnd();
  return `${prefix}${prefix ? "\n\n" : ""}${lines.join("\n")}\n`;
}

function resolveAppDirectory(projectRoot: string): string {
  const appRoot = path.join(projectRoot, "app");
  if (existsSync(appRoot)) {
    return appRoot;
  }

  const srcAppRoot = path.join(projectRoot, "src", "app");
  if (existsSync(srcAppRoot)) {
    return srcAppRoot;
  }

  throw new Error(
    "Could not find App Router directory. Expected either ./app or ./src/app.",
  );
}

async function readProjectPackageJson(
  projectRoot: string,
): Promise<ProjectPackageJson> {
  const packageJsonPath = path.join(projectRoot, "package.json");
  if (!existsSync(packageJsonPath)) {
    throw new Error(`No package.json found in ${projectRoot}`);
  }

  const content = await readFile(packageJsonPath, "utf8");
  try {
    return JSON.parse(content) as ProjectPackageJson;
  } catch {
    throw new Error(`Invalid JSON in ${packageJsonPath}`);
  }
}

function assertNextProject(pkg: ProjectPackageJson): void {
  const hasNextDependency = Boolean(
    pkg.dependencies?.next || pkg.devDependencies?.next,
  );
  if (!hasNextDependency) {
    throw new Error(
      "Target project is not a Next.js project (missing `next` dependency).",
    );
  }
}

function installSdkDependency(
  packageManager: PackageManager,
  projectRoot: string,
  source: DependencySource,
  dryRun: boolean,
  logger: Logger,
) {
  if (source === "workspace") {
    if (packageManager !== "pnpm") {
      throw new Error(
        "The --source workspace mode currently requires pnpm. Use pnpm in a workspace project or switch to --source npm.",
      );
    }
    if (!findPnpmWorkspaceRoot(projectRoot)) {
      throw new Error(
        "Could not find pnpm-workspace.yaml above the target project. The --source workspace mode requires the app to be inside a pnpm workspace.",
      );
    }
  }

  const commandByManagerBySource: Record<
    DependencySource,
    Record<PackageManager, [string, string[]]>
  > = {
    npm: {
      pnpm: ["pnpm", ["add", SDK_PACKAGE]],
      npm: ["npm", ["install", SDK_PACKAGE]],
      yarn: ["yarn", ["add", SDK_PACKAGE]],
    },
    workspace: {
      pnpm: ["pnpm", ["add", ...WORKSPACE_SDK_PACKAGES]],
      npm: ["npm", ["install", SDK_PACKAGE]],
      yarn: ["yarn", ["add", SDK_PACKAGE]],
    },
  };
  const [command, args] = commandByManagerBySource[source][packageManager];

  if (dryRun) {
    logger.log(
      `[dry-run] run ${[command, ...args].join(" ")} in ${projectRoot}`,
    );
    return;
  }

  const result = spawnSync(command, args, {
    cwd: projectRoot,
    stdio: "inherit",
    env: process.env,
  });
  if (result.status !== 0) {
    if (source === "workspace") {
      throw new Error(
        "Failed to install workspace-linked SDK dependencies. Ensure the target app is part of the same pnpm workspace, or use --source npm.",
      );
    }
    throw new Error(
      `Failed to install ${SDK_PACKAGE} with ${packageManager}. If you're testing before publish, rerun with --skip-install and install local tarballs for @zitadel/zitadel-js, @zitadel/react, and @zitadel/nextjs.`,
    );
  }
}

export async function runAddCommand(
  options: AddCommandOptions = {},
): Promise<AddCommandResult> {
  const logger = options.logger ?? console;
  const dryRun = Boolean(options.dryRun);
  const skipInstall = Boolean(options.skipInstall);
  const dependencySource = options.source ?? "npm";
  const authMode = options.auth ?? "oidc";
  if (dependencySource !== "npm" && dependencySource !== "workspace") {
    throw new Error("Invalid source mode. Expected one of: npm, workspace.");
  }
  if (authMode !== "oidc" && authMode !== "session") {
    throw new Error("Invalid auth mode. Expected one of: oidc, session.");
  }
  const withApi = Boolean(options.withApi);
  const withWebhook = Boolean(options.withWebhook);
  const projectRoot = path.resolve(options.cwd ?? process.cwd());

  const packageJson = await readProjectPackageJson(projectRoot);
  assertNextProject(packageJson);
  const fileExtension = isTypeScriptProject(projectRoot, packageJson)
    ? "ts"
    : "js";
  const appDirectory = resolveAppDirectory(projectRoot);
  const packageManager = detectPackageManager(projectRoot);
  const templates = resolveTemplates({ authMode, withApi, withWebhook });
  const requiredEnvVars = resolveEnvVars({ authMode, withApi, withWebhook });

  const createdFiles: string[] = [];
  const skippedFiles: string[] = [];

  for (const [routePath, template] of Object.entries(templates)) {
    const templateExtension = routePath.endsWith("/page")
      ? fileExtension === "ts"
        ? "tsx"
        : "jsx"
      : fileExtension;
    const relativePath = `${routePath}.${templateExtension}`;
    const targetPath = path.join(appDirectory, relativePath);
    const displayPath = path.relative(projectRoot, targetPath);

    if (existsSync(targetPath)) {
      skippedFiles.push(displayPath);
      logger.warn(`skip existing file: ${displayPath}`);
      continue;
    }

    createdFiles.push(displayPath);
    if (dryRun) {
      logger.log(`[dry-run] create ${displayPath}`);
      continue;
    }

    await mkdir(path.dirname(targetPath), { recursive: true });
    await writeFile(targetPath, template, "utf8");
    logger.log(`created ${displayPath}`);
  }

  const envExamplePath = path.join(projectRoot, ".env.example");
  const envExampleExists = existsSync(envExamplePath);
  const envExampleCurrent = envExampleExists
    ? await readFile(envExamplePath, "utf8")
    : "";
  const envExampleNext = mergeEnvExample(envExampleCurrent, requiredEnvVars);
  const envUpdated = envExampleNext !== envExampleCurrent;

  if (envUpdated) {
    const envDisplayPath = path.relative(projectRoot, envExamplePath);
    if (dryRun) {
      logger.log(`[dry-run] update ${envDisplayPath}`);
    } else {
      await writeFile(envExamplePath, envExampleNext, "utf8");
      logger.log(`updated ${envDisplayPath}`);
    }
  }

  if (!skipInstall) {
    installSdkDependency(
      packageManager,
      projectRoot,
      dependencySource,
      dryRun,
      logger,
    );
  }

  return {
    projectRoot,
    appDirectory: path.relative(projectRoot, appDirectory),
    packageManager,
    authMode,
    withApi,
    withWebhook,
    dependencySource,
    createdFiles,
    skippedFiles,
    envUpdated,
    dependencyInstalled: !skipInstall && !dryRun,
  };
}
