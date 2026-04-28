import {
  createNodesFromFiles,
  type CreateNodesV2,
  type CreateNodesResult,
} from "@nx/devkit";
import { dirname } from "path";
import { Config } from "./config.mts";

export const name = "install-proto";

export const createNodesV2: CreateNodesV2<{}> = [
  "**/package.json",
  async (configFiles, options, context) => {
    return await createNodesFromFiles(
      (configFile, _, _context) => createNodesInternal(configFile),
      configFiles,
      options,
      context,
    );
  },
];

async function createNodesInternal(
  configFile: string,
): Promise<CreateNodesResult> {
  const projectRoot = dirname(configFile);

  let config: Awaited<ReturnType<typeof Config.read>>;
  try {
    config = await Config.read(projectRoot);
  } catch {
    return {};
  }

  const outputs = config.packages
    .map((pkg) => pkg.extract.split("/").at(-1))
    .map((file) => `{workspaceRoot}/.artifacts/bin/${file}`);

  return {
    projects: {
      [projectRoot]: {
        targets: {
          ["install-proto"]: {
            executor: "nx:run-commands",
            options: {
              command: "node packages/install-proto/index.mts {projectRoot}",
            },
            cache: true,
            inputs: [
              "{projectRoot}/package.json",
              { runtime: "node -p 'process.platform'" },
              { runtime: "node -p 'process.arch'" },
              "{workspaceRoot}/packages/install-proto/**/*",
            ],
            outputs,
          },
        },
      },
    },
  };
}
