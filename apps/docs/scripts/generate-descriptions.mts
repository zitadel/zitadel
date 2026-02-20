/**
 * generate-descriptions.mts
 *
 * Generates SEO meta description frontmatter for MDX docs pages that are missing one,
 * using a local OpenAI-compatible LLM endpoint (e.g. LM Studio).
 *
 * Usage:
 *   node --experimental-strip-types apps/docs/scripts/generate-descriptions.mts [options]
 *
 * Options:
 *   --base-url <url>    LLM server base URL  (default: http://127.0.0.1:1234)
 *   --model <id>        Model identifier     (default: google/gemma-3-12b)
 *   --dry-run           Preview without writing files
 *   --file <path>       Process a single file only (relative to content/)
 *   --delay <ms>        Delay between API calls in ms (default: 200)
 *   --start-from <n>    Skip the first N eligible files (resume after interruption)
 *
 * Excluded directories (generated / versioned content that should not be edited):
 *   - content/reference/  (auto-generated gRPC / protobuf API reference)
 *   - content/v4.x/       (versioned doc snapshots downloaded at build time)
 */

import { readFileSync, writeFileSync } from 'node:fs';
import { readdir } from 'node:fs/promises';
import { join, relative, basename } from 'node:path';

// ---------------------------------------------------------------------------
// Directories (relative to content/) that contain generated or versioned files
// and must not be edited.
// ---------------------------------------------------------------------------
const EXCLUDED_DIR_PREFIXES = [
  'reference', // auto-generated API reference (proto/gRPC docs)
];
const EXCLUDED_DIR_PATTERNS = [
  /^v\d/, // versioned snapshots: v4.10, v4.11, etc.
];

function isExcluded(contentRelPath: string): boolean {
  const first = contentRelPath.split('/')[0];
  if (EXCLUDED_DIR_PREFIXES.includes(first)) return true;
  if (EXCLUDED_DIR_PATTERNS.some((re) => re.test(first))) return true;
  return false;
}

// ---------------------------------------------------------------------------
// CLI argument parsing
// ---------------------------------------------------------------------------

const args = process.argv.slice(2);

function getArg(flag: string): string | undefined {
  const idx = args.indexOf(flag);
  return idx !== -1 ? args[idx + 1] : undefined;
}

const BASE_URL = getArg('--base-url') ?? 'http://127.0.0.1:1234';
const MODEL = getArg('--model') ?? 'google/gemma-3-12b';
const DRY_RUN = args.includes('--dry-run');
const SINGLE_FILE = getArg('--file');
const DELAY_MS = Number(getArg('--delay') ?? '200');
const START_FROM = Number(getArg('--start-from') ?? '0');

const CONTENT_DIR = new URL('../content', import.meta.url).pathname;

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/** Recursively collect all .mdx file paths under a directory. */
async function globMdx(dir: string): Promise<string[]> {
  const entries = await readdir(dir, { withFileTypes: true, recursive: true });
  return entries
    .filter((e) => e.isFile() && e.name.endsWith('.mdx'))
    .map((e) => join(e.parentPath ?? join(dir, e.name), e.name));
}

/** Parse YAML frontmatter from an MDX string. Returns raw frontmatter block and body. */
function parseFrontmatter(content: string): {
  fm: Record<string, string>;
  rawFm: string;
  body: string;
} {
  const match = content.match(/^---\r?\n([\s\S]*?)\r?\n---\r?\n([\s\S]*)$/);
  if (!match) return { fm: {}, rawFm: '', body: content };

  const rawFm = match[1];
  const body = match[2];
  const fm: Record<string, string> = {};

  for (const line of rawFm.split('\n')) {
    const kv = line.match(/^(\w[\w_-]*):\s*(.*)$/);
    if (kv) fm[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, '');
  }

  return { fm, rawFm, body };
}

/** Strip MDX-specific syntax to get readable plain-text context for the LLM. */
function extractContext(body: string, maxChars = 900): string {
  return (
    body
      // Remove import statements
      .replace(/^import\s.*$/gm, '')
      // Remove JSX/HTML tags
      .replace(/<[^>]+>/g, ' ')
      // Replace code fences: keep only the first line (often the most informative command/type)
      // so the model can pick up concrete keywords (e.g. endpoint names, action types)
      .replace(/```(?:\w+)?\n([^\n]*)(?:[\s\S]*?)```/g, '[$1]')
      // Remove any remaining bare ``` delimiters
      .replace(/```/g, '')
      // Collapse extra whitespace
      .replace(/\n{3,}/g, '\n\n')
      .trim()
      .slice(0, maxChars)
  );
}

/**
 * Sanitize LLM output to prevent frontmatter corruption and ensure quality.
 * Strips markdown, conversational prefixes, and enforces a hard character cap.
 */
function sanitizeDescription(text: string, maxLen = 160): string {
  let clean = text.trim();

  // Strip surrounding quotes
  clean = clean.replace(/^["']|["']$/g, '').trim();

  // Strip common conversational/label prefixes the model might add
  clean = clean
    .replace(/^(meta\s+description:|description:|here is (the )?meta description:|sure[,!]?\s)/i, '')
    .trim();

  // Strip rogue markdown that could corrupt YAML (bold, italic, inline code)
  clean = clean.replace(/[*`_]/g, '');

  // Collapse double spaces left by markdown stripping
  clean = clean.replace(/  +/g, ' ').trim();

  // Remove trailing period for consistency
  clean = clean.replace(/\.$/, '');

  // Hard cap at maxLen: cut to last whole word, never append ellipses
  // (ellipses can push the string back over the cap and look bad in SERPs)
  if (clean.length > maxLen) {
    const truncated = clean.slice(0, maxLen);
    const lastSpace = truncated.lastIndexOf(' ');
    clean = lastSpace > 40 ? truncated.slice(0, lastSpace) : truncated.trim();
  }

  return clean;
}

/** Inject a description field into a frontmatter block, right after the title line. */
function injectDescription(content: string, description: string): string {
  // Escape any double quotes in the generated description
  const escaped = description.replace(/"/g, '\\"');

  // Insert after the title: line inside frontmatter
  if (/^title:/m.test(content)) {
    return content.replace(/^(title:.*)$/m, `$1\ndescription: "${escaped}"`);
  }

  // Fallback: insert as first field inside frontmatter block
  return content.replace(/^---\r?\n/, `---\ndescription: "${escaped}"\n`);
}

/** Perform a single chat completion call and return raw content. */
async function callLLM(messages: Array<{ role: string; content: string }>): Promise<string> {
  const response = await fetch(`${BASE_URL}/v1/chat/completions`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      model: MODEL,
      messages,
      max_tokens: 120,
      temperature: 0.3,
      stop: ['\n'], // prevent second lines or trailing commentary
    }),
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`LLM API error ${response.status}: ${text}`);
  }

  const json = (await response.json()) as {
    choices: Array<{ message: { content: string } }>;
  };

  return json.choices?.[0]?.message?.content ?? '';
}

/** Call the local LLM and return a sanitized, length-validated description.
 *  Automatically retries once if the first result is out of the 130–160 char range. */
async function generateDescription(
  title: string,
  context: string,
): Promise<string> {
  const systemPrompt =
    'You write SEO meta descriptions for ZITADEL technical documentation. ' +
    'Return ONLY the meta description text (no labels, no JSON, no quotes). ' +
    'Write exactly ONE sentence in plain text. ' +
    'Length MUST be 130–160 characters INCLUDING spaces. Do NOT exceed 160. ' +
    'No markdown, no emojis, no exclamation marks, no ellipses. Do not end with a period. ' +
    'Avoid semicolons; one sentence, you may use one comma clause. ' +
    'Avoid generic openings like "This guide", "This page", "Learn how", "In this article". ' +
    'Include "ZITADEL" if it fits naturally. ' +
    'Prefer one concrete technical keyword from the content (e.g. Actions V2, webhook, OIDC, PAT). ' +
    'If the page covers deprecation or migration, mention it neutrally.';

  const userPrompt =
    `Title: ${title}\n` +
    `Excerpt:\n${context}\n\n` +
    `Meta description:`;

  const messages = [
    { role: 'system', content: systemPrompt },
    { role: 'user', content: userPrompt },
  ];

  const rawFirst = await callLLM(messages);
  let desc = sanitizeDescription(rawFirst);

  // Score the RAW output length, not the sanitized length.
  // Sanitizing a 200-char response to 158 chars would make distToRange return 0,
  // masking the fact that the model didn't follow the length constraint at all.
  const distToRange = (len: number) => Math.max(0, 130 - len, len - 160);

  // Retry whenever the raw output missed the target range (too short or too long).
  if (distToRange(rawFirst.length) > 0) {
    const retryMessages = [
      ...messages,
      { role: 'assistant', content: rawFirst },
      {
        role: 'user',
        content:
          `That was ${rawFirst.length} characters. Rewrite to be 130–160 characters (inclusive), including spaces. ` +
          `Make sure it is a complete sentence that does not need to be cut off. Keep the technical meaning.\n\nMeta description:`,
      },
    ];
    const rawRetry = await callLLM(retryMessages);
    // Only replace if the retry obeyed the rules better (judge both by raw length)
    if (distToRange(rawRetry.length) < distToRange(rawFirst.length)) {
      desc = sanitizeDescription(rawRetry);
    }
  }

  return desc;
}

/** Sleep helper. */
const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

async function main() {
  console.log(`LLM endpoint : ${BASE_URL}`);
  console.log(`Model        : ${MODEL}`);
  console.log(`Dry run      : ${DRY_RUN}`);
  console.log(`Delay        : ${DELAY_MS}ms`);
  console.log(`Excluded     : reference/, v4.x/\n`);

  // Build file list
  let files: string[];
  if (SINGLE_FILE) {
    files = [join(CONTENT_DIR, SINGLE_FILE)];
  } else {
    const all = await globMdx(CONTENT_DIR);
    files = all
      .filter((f) => {
        const rel = relative(CONTENT_DIR, f);
        // Skip partial files (filenames starting with _)
        if (basename(f).startsWith('_')) return false;
        // Skip generated/versioned directories
        if (isExcluded(rel)) return false;
        return true;
      });
    files.sort();
  }

  // Filter to files missing a description
  const eligible = files.filter((f) => {
    const content = readFileSync(f, 'utf-8');
    const { fm } = parseFrontmatter(content);
    return !fm['description'];
  });

  console.log(`Total non-generated MDX files : ${files.length}`);
  console.log(`Missing description           : ${eligible.length}`);
  if (START_FROM > 0) console.log(`Skipping first                : ${START_FROM}`);
  console.log('');

  const toProcess = eligible.slice(START_FROM);
  let succeeded = 0;
  let failed = 0;

  for (let i = 0; i < toProcess.length; i++) {
    const filePath = toProcess[i];
    const relPath = relative(CONTENT_DIR, filePath);
    const globalIndex = START_FROM + i + 1;
    const progress = `[${globalIndex}/${eligible.length}]`;

    const content = readFileSync(filePath, 'utf-8');
    const { fm, body } = parseFrontmatter(content);
    const title = fm['title'] ?? basename(filePath, '.mdx');
    const context = extractContext(body);

    process.stdout.write(`${progress} ${relPath} … `);

    try {
      const description = await generateDescription(title, context);
      const charCount = description.length;

      if (DRY_RUN) {
        console.log(`\n  → (dry-run) ${description} [${charCount} chars]`);
      } else {
        const updated = injectDescription(content, description);
        writeFileSync(filePath, updated, 'utf-8');
        console.log(`done [${charCount} chars]`);
      }

      succeeded++;
    } catch (err) {
      console.log(`FAILED — ${(err as Error).message}`);
      failed++;
    }

    if (i < toProcess.length - 1) await sleep(DELAY_MS);
  }

  console.log(`\n✓ Succeeded: ${succeeded}  ✗ Failed: ${failed}`);
  if (failed > 0) process.exit(1);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
