// components/github-code-block.tsx
import { CodeBlock, Pre } from 'fumadocs-ui/components/codeblock';

export async function GithubCodeBlock({ url }: { url: string }) {
    // 1. Parse the URL to get the raw content URL and line fragments
    const urlObj = new URL(url);
    const rawUrl = url
        .replace('github.com', 'raw.githubusercontent.com')
        .replace('/blob/', '/');

    // Extract line numbers from hash (e.g., #L10-L20 or #L5)
    const lineMatch = urlObj.hash.match(/L(\d+)(?:-L(\d+))?/);
    const startLine = lineMatch ? parseInt(lineMatch[1], 10) : null;
    const endLine = lineMatch ? (lineMatch[2] ? parseInt(lineMatch[2], 10) : startLine) : null;

    // 2. Fetch the content
    const response = await fetch(rawUrl);
    let code = await response.text();

    // 3. Optional: Extract specific lines
    if (startLine !== null && endLine !== null) {
        const lines = code.split('\n');
        const selectedLines = lines.slice(startLine - 1, endLine);

        // Remove common indentation (dedent)
        const minIndent = selectedLines.reduce((min, line) => {
            if (line.trim().length === 0) return min;
            const indent = line.match(/^\s*/)?.[0].length ?? 0;
            return Math.min(min, indent);
        }, Infinity);

        code = selectedLines
            .map(line => line.slice(minIndent === Infinity ? 0 : minIndent))
            .join('\n');
    }

    const lang = url.split('.').pop()?.split('#')[0] || 'text';

    return (
        <CodeBlock lang={lang}>
            <Pre>{code}</Pre>
        </CodeBlock>
    );
}