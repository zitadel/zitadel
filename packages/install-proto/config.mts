import { readFile } from "node:fs/promises";
import { join } from "node:path";

type InstallProtoPackage = {
  version: string;
  fileTemplate: string;
  extractTemplate: string;
  overrides?: Record<string, string>;
};

type InstallProtoConfig = {
  packages: Record<string, InstallProtoPackage>;
};

export interface Package {
  get owner(): string;
  get repo(): string;
  get version(): string;
  get file(): string;
  get extract(): string;
}

class PackageImpl implements Package {
  readonly #owner: string;
  readonly #repo: string;
  #pkg: InstallProtoPackage;

  #platform = process.platform;
  #arch = process.arch;

  constructor(name: string, pkg: InstallProtoPackage) {
    const [owner, repo] = name.split("/");
    this.#owner = owner;
    this.#repo = repo;
    this.#pkg = pkg;
  }

  private template(template: string) {
    const defaultOverrides = Object.entries({
      "{VERSION}": this.version.replace("v", ""),
      "{PLATFORM}": this.#platform,
      "{ARCH}": this.#arch,
      "{EXE}": this.#platform === "win32" ? ".exe" : "",
    });

    const overrides = Object.entries(this.#pkg.overrides ?? {});

    return [...defaultOverrides, ...overrides].reduce(
      (acc, [from, to]) => acc.replaceAll(from, to),
      template,
    );
  }

  get version(): string {
    return this.#pkg.version;
  }

  get owner() {
    return this.#owner;
  }

  get repo() {
    return this.#repo;
  }

  get file(): string {
    return this.template(this.#pkg.fileTemplate);
  }

  get extract() {
    return this.template(this.#pkg.extractTemplate);
  }
}

export class Config {
  static async read(inputDir: string) {
    const path = join(inputDir, "package.json");

    const content = await readFile(path, "utf8");
    const packageJson = JSON.parse(content);

    if (
      !("installProto" in packageJson) ||
      typeof packageJson.installProto !== "object" ||
      packageJson.installProto === null
    ) {
      throw new Error("installProto config not found in package.json");
    }

    return new Config(packageJson.installProto);
  }

  readonly #packages: Package[];

  private constructor(config: InstallProtoConfig) {
    this.#packages = Object.entries(config.packages).map(
      ([name, pkg]) => new PackageImpl(name, pkg),
    );
  }

  get packages() {
    return this.#packages;
  }
}
