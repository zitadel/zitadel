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
 *   --limit <n>         Process at most N files then stop
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
const MODEL = getArg('--model') ?? 'qwen2.5-coder-14b-instruct';
const DRY_RUN = args.includes('--dry-run');
const SINGLE_FILE = getArg('--file');
const DELAY_MS = Number(getArg('--delay') ?? '200');
const START_FROM = Number(getArg('--start-from') ?? '0');
const LIMIT = getArg('--limit') !== undefined ? Number(getArg('--limit')) : Infinity;

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

/**
 * Build a compact, representative context string for the LLM.
 *
 * Strategy:
 *   1. Extract all ## / ### headings as a page-scope outline (the model sees
 *      every major topic even on long pages).
 *   2. Append the first ~900 chars of cleaned body prose so specific keywords
 *      and phrasing are available.
 *
 * This avoids both the "only summarised the first subsection" problem (fixed by
 * the outline) and the "fed 3000 chars of repetitive code fences" problem (fixed
 * by the char cap on the prose portion).
 */
function extractContext(body: string, maxBodyChars = 900): string {
  // --- 1. Strip code fences, imports, and JSX before doing anything else ------
  const cleaned = body
    .replace(/^import\s.*$/gm, '')
    .replace(/<[^>]+>/g, ' ')
    // Keep only the first (most informative) line of each code fence
    .replace(/```(?:\w+)?\n([^\n]*)(?:[\s\S]*?)```/g, '[$1]')
    .replace(/```/g, '')
    .replace(/\n{3,}/g, '\n\n')
    .trim();

  // --- 2. Build heading outline (## and ###, deduplicated) --------------------
  const headings: string[] = [];
  for (const line of cleaned.split('\n')) {
    const m = line.match(/^(#{2,3})\s+(.+)/);
    if (m) headings.push(`${m[1]} ${m[2].trim()}`);
  }
  const outline =
    headings.length > 0 ? `Page sections:\n${headings.join('\n')}\n\n` : '';

  // --- 3. Prose excerpt (non-heading lines, up to maxBodyChars) ---------------
  const prose = cleaned
    .split('\n')
    .filter((l) => !l.match(/^#{1,6}\s/))
    .join('\n')
    .slice(0, maxBodyChars)
    .trim();

  return (outline + prose).trim();
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

  // Guard against YAML-breaking newlines (belt-and-suspenders: stop seqs should
  // prevent these, but a model that prefixes a blank line can still sneak one in)
  clean = clean.replace(/\r?\n+/g, ' ');

  // Collapse double spaces left by markdown stripping or newline removal
  clean = clean.replace(/  +/g, ' ').trim();

  // Strip leading/trailing dashes, em-dashes, and colons the model occasionally adds
  clean = clean.replace(/^[\s\-–—:]+/, '').trim();
  clean = clean.replace(/[\s\-–—:]+$/, '').trim();

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
async function callLLM(
  messages: Array<{ role: string; content: string }>,
  stopSeqs: string[] = ['\n', '\r\n'],
): Promise<string> {
  const response = await fetch(`${BASE_URL}/v1/chat/completions`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      model: MODEL,
      messages,
      max_tokens: 120,
      temperature: 0.3,
      stop: stopSeqs, // prevent second lines or trailing commentary
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

/**
 * Returns true when the raw LLM output has obvious format problems that the
 * sanitizer can fix but that also indicate the model ignored the instructions.
 * Used to decide whether to retry even when the length happens to be fine.
 */
function looksBad(raw: string): boolean {
  const t = raw.trim();
  if (!t) return true; // empty — model hit stop on a leading newline
  if (/^["'].*["']$/.test(t)) return true; // fully quoted
  // Label prefixes: "Meta description:", "Meta description -", "- item", "* item"
  if (/^(meta\s*description\s*[:\-]|description\s*[:\-])/i.test(t)) return true;
  if (/^[-*]\s+/.test(t)) return true; // bullet list start
  if (/[`*_]/.test(t)) return true; // markdown artifacts
  return false;
}

interface GenerateResult {
  description: string;
  rawFirstLen: number;
  retried: boolean;
}

/** Call the local LLM and return a sanitized, length-validated description.
 *  Retries once if the first result is out of the 130–160 char range OR has bad formatting. */
async function generateDescription(
  title: string,
  context: string,
): Promise<GenerateResult> {
  const systemPrompt =
    'You write SEO meta descriptions for ZITADEL technical documentation. ' +
    'Return ONLY the meta description text (no labels, no JSON, no quotes). ' +
    'Write exactly ONE sentence in plain text. ' +
    'Length MUST be 130–160 characters INCLUDING spaces. Do NOT exceed 160. ' +
    'No markdown, no emojis, no exclamation marks, no ellipses. Do not end with a period. ' +
    'Avoid semicolons; one sentence, you may use one comma clause. ' +
    'Avoid generic openings like "This guide", "This page", "Learn how", "In this article". ' +
    'Include "ZITADEL" when it improves clarity. ' +
    'Summarize the overall scope of the page (top 2–3 topics), not a single subsection. ' +
    'If the excerpt contains multiple major sections (e.g., OIDC, SAML, users, IDPs), reflect 2–3 of them rather than focusing on the first. ' +
    'Extract 1–2 relevant technical keywords directly from the content. Do NOT invent version numbers or features not present in the excerpt. ' +
    'If the page covers deprecation or migration, mention it neutrally.'

  const userPrompt =
    `Title: ${title}\n` +
    `Excerpt:\n${context}\n\n` +
    `Meta description:`;

  const messages = [
    { role: 'system', content: systemPrompt },
    { role: 'user', content: userPrompt },
  ];

  // If the model starts with a leading newline, stop:[LF/CRLF] yields empty output.
  // Retry once with double-newline stops to recover.
  let rawFirst = await callLLM(messages);
  if (!rawFirst.trim()) {
    rawFirst = await callLLM(messages, ['\n\n', '\r\n\r\n']);
  }

  let desc = sanitizeDescription(rawFirst);

  // Score the RAW output length, not the sanitized length.
  // Sanitizing a 200-char response to 158 chars would make distToRange return 0,
  // masking the fact that the model didn't follow the length constraint at all.
  const distToRange = (len: number) => Math.max(0, 130 - len, len - 160);

  // Retry when the raw output missed the length target OR had obvious format problems
  // (quotes, label prefix, markdown) — so the model gets corrective feedback.
  const needsRetry = distToRange(rawFirst.trim().length) > 0 || looksBad(rawFirst);
  let retried = false;

  if (needsRetry) {
    retried = true;
    const retryMessages = [
      ...messages,
      { role: 'assistant', content: rawFirst },
      {
        role: 'user',
        content:
          `That was ${rawFirst.trim().length} characters. Rewrite to be 130–160 characters (inclusive), including spaces. ` +
          `Return ONLY plain text — no quotes, no labels, no markdown. ` +
          `Make sure it is a complete sentence that does not need to be cut off. Keep the technical meaning.\n\nMeta description:`,
      },
    ];
    const rawRetry = await callLLM(retryMessages);
    // Accept the retry if it is better-formatted AND no further from the target range
    if (!looksBad(rawRetry) && distToRange(rawRetry.trim().length) <= distToRange(rawFirst.trim().length)) {
      desc = sanitizeDescription(rawRetry);
    }
  }

  return { description: desc, rawFirstLen: rawFirst.trim().length, retried };
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

  const toProcess = eligible.slice(START_FROM, START_FROM + (isFinite(LIMIT) ? LIMIT : eligible.length));
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
      const { description, rawFirstLen, retried } = await generateDescription(title, context);
      const charCount = description.length;

      if (DRY_RUN) {
        const retryTag = retried ? ` retry=yes raw1=${rawFirstLen}` : ` raw1=${rawFirstLen}`;
        console.log(`\n  → (dry-run) ${description} [${charCount} chars${retryTag}]`);
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
