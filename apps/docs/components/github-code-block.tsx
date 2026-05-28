import { DynamicCodeBlock } from 'fumadocs-ui/components/dynamic-codeblock';

const REVALIDATE_SECONDS = 60 * 60 * 24;

function toRawGithubUrl(url: string): string {
    const parsed = new URL(url);
    if (parsed.protocol !== 'https:' || parsed.hostname !== 'github.com') {
        throw new Error(`GithubCodeBlock only accepts https://github.com URLs, got: ${url}`);
    }
    const segments = parsed.pathname.split('/').filter(Boolean);
    if (segments.length < 5 || segments[2] !== 'blob') {
        throw new Error(`Unexpected GitHub URL shape, expected /<owner>/<repo>/blob/<ref>/<path>, got: ${url}`);
    }
    const [owner, repo, , ref, ...path] = segments;
    return `https://raw.githubusercontent.com/${owner}/${repo}/${ref}/${path.join('/')}`;
}

function parseLineRange(hash: string): { start: number; end: number } | null {
    const match = hash.match(/L(\d+)(?:-L(\d+))?/);
    if (!match) return null;
    const start = parseInt(match[1], 10);
    const end = match[2] ? parseInt(match[2], 10) : start;
    return { start, end };
}

function sliceAndDedent(code: string, range: { start: number; end: number }): string {
    const selected = code.split('\n').slice(range.start - 1, range.end);
    const minIndent = selected.reduce((min, line) => {
        if (line.trim().length === 0) return min;
        const indent = line.match(/^\s*/)?.[0].length ?? 0;
        return Math.min(min, indent);
    }, Infinity);
    const dedent = minIndent === Infinity ? 0 : minIndent;
    return selected.map(line => line.slice(dedent)).join('\n');
}

export async function GithubCodeBlock({ url }: { url: string }) {
    const rawUrl = toRawGithubUrl(url);
    const response = await fetch(rawUrl, { next: { revalidate: REVALIDATE_SECONDS } });
    if (!response.ok) {
        throw new Error(`GithubCodeBlock failed to fetch ${rawUrl}: ${response.status} ${response.statusText}`);
    }
    let code = await response.text();

    const range = parseLineRange(new URL(url).hash);
    if (range) {
        code = sliceAndDedent(code, range);
    }

    const lang = new URL(url).pathname.split('.').pop() || 'text';

    return <DynamicCodeBlock lang={lang} code={code} />;
}
