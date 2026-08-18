import { execFile as execFileCallback } from "node:child_process";
import { promisify } from "node:util";

const execFile = promisify(execFileCallback);

export const goos = (
  await execFile("go", ["env", "GOOS"], { encoding: "utf8" })
).stdout.trim();
export const goarch = (
  await execFile("go", ["env", "GOARCH"], { encoding: "utf8" })
).stdout.trim();
