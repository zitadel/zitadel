import { DynamicCodeBlock } from 'fumadocs-ui/components/dynamic-codeblock';

export async function GithubCodeBlock({ url }: { url: string }) {
    const urlObj = new URL(url);
    const rawUrl = url
        .replace('github.com', 'raw.githubusercontent.com')
        .replace('/blob/', '/');

    const lineMatch = urlObj.hash.match(/L(\d+)(?:-L(\d+))?/);
    const startLine = lineMatch ? parseInt(lineMatch[1], 10) : null;
    const endLine = lineMatch ? (lineMatch[2] ? parseInt(lineMatch[2], 10) : startLine) : null;

    const response = await fetch(rawUrl);
    let code = await response.text();

    if (startLine !== null && endLine !== null) {
        const lines = code.split('\n');
        const selectedLines = lines.slice(startLine - 1, endLine);

        const minIndent = selectedLines.reduce((min, line) => {
            if (line.trim().length === 0) return min;
            const indent = line.match(/^\s*/)?.[0].length ?? 0;
            return Math.min(min, indent);
        }, Infinity);

        code = selectedLines
            .map(line => line.slice(minIndent === Infinity ? 0 : minIndent))
            .join('\n');
    }

    const pathname = urlObj.pathname.split('#')[0];
    const lang = pathname.split('.').pop() || 'text';

    return <DynamicCodeBlock lang={lang} code={code} />;
}
