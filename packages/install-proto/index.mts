import { finished } from "node:stream/promises";
import * as childProcess from "node:child_process";
import { createHash, Hash } from "node:crypto";
import { Config, type Package } from "./config.mts";
import { download, getAsset } from "./download.mts";
import { type Readable } from "node:stream";
import { createWriteStream } from "node:fs";
import { chmod, mkdir } from "node:fs/promises";
import { once } from "node:events";
import path from "node:path";

const outputDir = path.resolve(
  process.env.NX_WORKSPACE_ROOT ?? "",
  ".artifacts/bin",
);
const inputDir = path.resolve(process.argv[2]);

console.log(`starting install-proto input:${inputDir} output:${outputDir}`);

function writeToFile(pkg: Package, readable: Readable) {
  const writeStream = createWriteStream(path.join(outputDir, pkg.extract));
  return finished(readable.pipe(writeStream), { cleanup: true });
}

// maybe change this to outputting to stdout so we can reuse the file create code above
// and have a seperate file name
async function extractToFile(pkg: Package, readable: Readable) {
  const flags = pkg.file.endsWith(".tar.gz") ? "-xzf" : "-xf";
  const stripComponents = pkg.extract.split("/").length - 1;
  const tarArgs = [
    flags,
    "-",
    "--strip-components",
    stripComponents.toString(),
    "-C",
    outputDir,
    pkg.extract,
  ];

  const tar = childProcess.spawn("tar", tarArgs, {
    stdio: ["pipe", "inherit", "inherit"],
  });

  await finished(readable.pipe(tar.stdin), { cleanup: true });
  const [code] = await once(tar, "close");

  if (code !== 0) {
    throw new Error(`tar failed extracting ${pkg.file} with exit code ${code}`);
  }
}

async function makeExecutable(pkg: Package) {
  if (process.platform === "win32") {
    return;
  }

  const extract = pkg.extract.split("/").at(-1);
  if (!extract) {
    throw new Error(`Could not determine executable name for ${pkg.file}`);
  }

  const executable = path.join(outputDir, extract);
  await chmod(executable, 0o755);

  console.log(`made ${executable} executable`);
}

class CheckHash {
  #hash: Hash;
  #digest: string | undefined;

  constructor(readable: Readable) {
    this.#hash = createHash("sha256");
    readable.on("data", (chunk) => this.#hash.update(chunk));
  }

  get digest() {
    return (this.#digest ??= `sha256:${this.#hash.digest("hex")}`);
  }
}

async function processPackage(pkg: Package) {
  console.log(`processing ${pkg.owner}/${pkg.repo}/${pkg.file}`);

  const asset = await getAsset(pkg);
  const { browser_download_url, digest } = asset;
  const readStream = await download(browser_download_url);

  const checkHash = new CheckHash(readStream);

  if (pkg.file.endsWith(".tar.gz") || pkg.file.endsWith(".zip")) {
    await extractToFile(pkg as Package, readStream);
  } else {
    await writeToFile(pkg, readStream);
  }

  if (checkHash.digest !== digest) {
    throw new Error(
      `SHA256 mismatch for ${pkg.file}: expected ${digest}, got ${checkHash.digest}`,
    );
  }

  console.log(`successfully downloaded ${pkg.owner}/${pkg.repo}/${pkg.file}`);

  await makeExecutable(pkg);
}

await mkdir(outputDir, { recursive: true });
const config = await Config.read(inputDir);
const promises = config.packages.map((pkg) => processPackage(pkg));
await Promise.all(promises);

console.log("success");
