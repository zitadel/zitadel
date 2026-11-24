// Keeps the Nx task cache of CI in workflow artifacts instead of Nx Cloud.
//
// Nx only restores cache entries it recorded itself, so a copied .nx/cache folder
// is not enough. Instead, a job downloads earlier artifacts into a folder and
// serves it to Nx over Nx's self-hosted remote cache API
// (https://nx.dev/docs/kb/self-hosted-caching).
//
// Artifacts are bound to the run that uploaded them, so a job restores only from
// runs on main of this repository and from earlier runs of its own branch. The
// Actions cache can't give that guarantee on Depot runners: all of this
// repository's branches share one Depot cache namespace, so any branch's job
// can write entries main restores.
//
//   node nx-task-cache.mjs find <artifact>          write main-run=<id> and branch-run=<id> to $GITHUB_OUTPUT
//   node nx-task-cache.mjs serve <store>            serve <store> on a free 127.0.0.1 port, written to <store>.url
//   node nx-task-cache.mjs prune <store> [<skip>]   delete entries this job did not use, and those listed in <skip>
//
// serve logs each request to <store>.log as "<event> <hash> <bytes>".
import { createServer } from "node:http";
import {
  appendFileSync,
  createReadStream,
  createWriteStream,
  existsSync,
  readFileSync,
  readdirSync,
  renameSync,
  rmSync,
  statSync,
  writeFileSync,
} from "node:fs";
import { join } from "node:path";
import { pipeline } from "node:stream/promises";

const [command, arg, skipList] = process.argv.slice(2);

if (command === "find" && arg) {
  // Without a cache the job is slower, not wrong.
  await find(arg).catch((error) => console.log(`::warning::Running without the Nx task cache: ${error.message}`));
} else if (command === "serve" && arg) {
  serve(arg);
} else if (command === "prune" && arg) {
  prune(arg, skipList);
} else {
  console.error("usage: nx-task-cache.mjs find <artifact> | serve <store> | prune <store> [<skip>]");
  process.exit(2);
}

// Finds the newest unexpired artifact from main of this repository, and from the
// branch this job runs for (HEAD_BRANCH in HEAD_REPOSITORY_ID, which is a fork for
// pull requests from forks).
async function find(artifact) {
  const { GITHUB_API_URL, GITHUB_REPOSITORY, GITHUB_REPOSITORY_ID, GITHUB_TOKEN, HEAD_BRANCH, HEAD_REPOSITORY_ID } =
    process.env;
  const response = await fetch(
    `${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/actions/artifacts?name=${artifact}&per_page=100`,
    { headers: { Accept: "application/vnd.github+json", Authorization: `Bearer ${GITHUB_TOKEN}` } },
  );
  if (!response.ok) {
    throw new Error(`listing ${artifact} artifacts returned ${response.status}`);
  }
  const artifacts = (await response.json()).artifacts
    .filter((a) => !a.expired && a.workflow_run)
    .sort((a, b) => b.created_at.localeCompare(a.created_at));
  const newest = (branch, repositoryId) =>
    artifacts.find(
      (a) => a.workflow_run.head_branch === branch && String(a.workflow_run.head_repository_id) === String(repositoryId),
    )?.workflow_run.id ?? "";
  const onMain = HEAD_BRANCH === "main" && HEAD_REPOSITORY_ID === GITHUB_REPOSITORY_ID;
  const outputs = `main-run=${newest("main", GITHUB_REPOSITORY_ID)}\nbranch-run=${onMain ? "" : newest(HEAD_BRANCH, HEAD_REPOSITORY_ID)}\n`;
  console.log(outputs);
  if (process.env.GITHUB_OUTPUT) {
    appendFileSync(process.env.GITHUB_OUTPUT, outputs);
  }
}

function serve(store) {
  const log = `${store}.log`;
  const record = (event, hash, bytes) => appendFileSync(log, `${event} ${hash} ${bytes}\n`);
  const server = createServer(async (req, res) => {
    const hash = /^\/v1\/cache\/(\w+)$/.exec(req.url ?? "")?.[1];
    if (!hash) {
      return end(res, 404);
    }
    const file = join(store, hash);
    // Written beside the store, so neither prune nor the upload sees a partial entry.
    const partial = `${store}.${hash}.${process.pid}.partial`;
    try {
      if (req.method === "GET") {
        if (!existsSync(file)) {
          record("miss", hash, 0);
          return end(res, 404);
        }
        const size = statSync(file).size;
        record("hit", hash, size);
        res.writeHead(200, { "Content-Type": "application/octet-stream", "Content-Length": size });
        await pipeline(createReadStream(file), res);
      } else if (req.method === "PUT") {
        if (existsSync(file)) {
          req.resume();
          record("exists", hash, 0);
          return end(res, 409);
        }
        await pipeline(req, createWriteStream(partial));
        const size = statSync(partial).size;
        const expected = req.headers["content-length"];
        if (expected !== undefined && String(size) !== expected) {
          rmSync(partial, { force: true });
          return end(res, 400, `expected ${expected} bytes, got ${size}`);
        }
        renameSync(partial, file);
        record("store", hash, size);
        end(res, 200);
      } else {
        end(res, 405);
      }
    } catch (error) {
      console.error(`${req.method} ${hash}: ${error}`);
      rmSync(partial, { force: true });
      if (res.headersSent) {
        res.destroy();
      } else {
        end(res, 500, String(error));
      }
    }
  });
  server.listen(0, "127.0.0.1", () => {
    writeFileSync(`${store}.url`, `http://127.0.0.1:${server.address().port}`);
  });
}

function prune(store, skipList) {
  const log = `${store}.log`;
  const events = existsSync(log)
    ? readFileSync(log, "utf8").split("\n").filter(Boolean).map((line) => line.split(" "))
    : [];
  const used = new Set(events.filter(([event]) => event !== "miss").map(([, hash]) => hash));
  const skip = new Set(skipList && existsSync(skipList) ? readFileSync(skipList, "utf8").split("\n") : []);
  let kept = 0;
  let bytes = 0;
  for (const name of readdirSync(store)) {
    if (used.has(name) && !skip.has(name)) {
      kept++;
      bytes += statSync(join(store, name)).size;
    } else {
      rmSync(join(store, name), { force: true });
    }
  }
  const count = (event) => events.filter(([e]) => e === event).length;
  const summary = `Nx task cache: ${count("hit")} hits, ${count("miss")} misses, ${count("store")} stored; uploading ${kept} entries (${Math.round(bytes / 1e6)} MB)`;
  console.log(summary);
  if (process.env.GITHUB_STEP_SUMMARY) {
    appendFileSync(process.env.GITHUB_STEP_SUMMARY, `${summary}\n`);
  }
}

function end(res, status, message) {
  res.writeHead(status, { "Content-Type": "text/plain" });
  res.end(message);
}
