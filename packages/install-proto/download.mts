import { Readable } from "node:stream";
import type { ReadableStream as NodeReadableStream } from "node:stream/web";
import type { Package } from "./config.mts";
import { Octokit } from "octokit";

// could be replaced by just calling the api directly to avoid depdency
// https://docs.github.com/en/rest/releases/releases?apiVersion=2026-03-10#get-a-release-by-tag-name
const octokit = new Octokit({
  auth: process.env.GH_TOKEN,
});

export async function getAsset(pkg: Package) {
  const res = await octokit.rest.repos.getReleaseByTag({
    owner: pkg.owner,
    repo: pkg.repo,
    tag: pkg.version,
  });

  if (!res.data) {
    throw new Error(
      `Release not found for ${pkg.owner}/${pkg.repo}/v${pkg.version}`,
    );
  }

  const asset = res.data.assets.find((asset) => asset.name === pkg.file);
  if (!asset) {
    throw new Error(`Asset not found for ${pkg.file}`);
  }

  console.log(`found asset for ${pkg.owner}/${pkg.repo}/${pkg.file}`);

  return asset;
}

export async function download(url: string) {
  const { ok, body, statusText } = await fetch(url, {
    redirect: "follow",
  });
  if (!ok) {
    throw new Error(`Failed to download ${url}: ${statusText}`);
  }

  if (!body) {
    throw new Error(`Missing response body for ${url}`);
  }

  return Readable.fromWeb(body as NodeReadableStream);
}
