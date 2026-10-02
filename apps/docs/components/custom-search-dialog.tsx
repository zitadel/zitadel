'use client';

import { useEffect, useState, ReactNode } from 'react';
import { useParams } from 'next/navigation';
import {
    SearchDialog,
    SearchDialogClose,
    SearchDialogContent,
    SearchDialogHeader,
    SearchDialogIcon,
    SearchDialogInput,
    SearchDialogList,
    SearchDialogOverlay,
    type SharedProps,
} from 'fumadocs-ui/components/dialog/search';
import { getVersionFromSlug } from '@/lib/versions';

type SearchResultItem = {
    id: string;
    content: ReactNode;
    url: string;
    type: 'page';
};

// --- RegExp Escape Utility ---
function escapeRegExp(string: string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

// --- MARKDOWN STRIPPER ---
// Converts raw MDX/Markdown syntax into clean prose
function stripMarkdown(markdown: string): string {
    if (!markdown) return '';
    return markdown
        .replace(/```[\s\S]*?```/g, '')             // Remove multi-line code blocks
        .replace(/`([^`]+)`/g, '$1')                 // Remove inline code ticks
        .replace(/#{1,6}\s+/g, '')                  // Remove headers (###)
        .replace(/!\[.*?\]\(.*?\)/g, '')            // Remove images
        .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')     // Convert links [text](url) -> text
        .replace(/[*_~]{1,3}([^*_~]+)[*_~]{1,3}/g, '$1') // Remove bold / italic / strikethrough
        .replace(/^\s*>\s+/gm, '')                  // Remove blockquotes
        .replace(/<[^>]*>/g, '')                    // Remove HTML / MDX tags
        .replace(/\s+/g, ' ')                       // Collapse multiple whitespace/newlines
        .trim();
}

// --- POLISHED HIGHLIGHT COMPONENT ---
const HighlightMatch = ({ text, query }: { text: string; query: string }) => {
    const cleanText = stripMarkdown(text);

    if (!query || !cleanText) {
        return <span className="text-xs text-muted-foreground line-clamp-2">{cleanText}</span>;
    }

    const safeRegex = new RegExp(`(${escapeRegExp(query)})`, 'gi');
    const parts = cleanText.split(safeRegex);

    return (
        <span className="text-xs text-muted-foreground line-clamp-2">
            {parts.map((part, i) =>
                part.toLowerCase() === query.toLowerCase() ? (
                    <mark
                        key={i}
                        className="bg-primary/20 text-primary font-semibold rounded-sm px-1 py-0.5"
                    >
                        {part}
                    </mark>
                ) : (
                    <span key={i}>{part}</span>
                )
            )}
        </span>
    );
};

export default function CustomSearchDialog(props: SharedProps) {
    const [query, setQuery] = useState('');
    const [debouncedQuery, setDebouncedQuery] = useState('');
    const [results, setResults] = useState<SearchResultItem[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    const params = useParams();
    const slug = params?.slug as string[] | undefined;
    const currentVersion = getVersionFromSlug(slug);

    useEffect(() => {
        const handler = setTimeout(() => {
            setDebouncedQuery(query);
        }, 300);
        return () => clearTimeout(handler);
    }, [query]);

    useEffect(() => {
        if (!debouncedQuery.trim()) {
            setResults([]);
            setIsLoading(false);
            return;
        }

        setIsLoading(true);
        const abortController = new AbortController();

        async function fetchResults() {
            try {
                const searchUrl = `/docs/api/search?q=${encodeURIComponent(debouncedQuery)}&tag=${encodeURIComponent(currentVersion)}`;
                const res = await fetch(searchUrl, { signal: abortController.signal });

                if (!res.ok) throw new Error('Search failed');

                const data = await res.json();

                const seenUrls = new Set<string>();
                const formattedResults: SearchResultItem[] = [];

                for (const item of data) {
                    const targetUrl = item.url;

                    if (!seenUrls.has(targetUrl)) {
                        seenUrls.add(targetUrl);

                        formattedResults.push({
                            id: targetUrl,
                            content: (
                                <div className="flex flex-col gap-1 py-0.5">
                                    <span className="font-medium text-foreground">
                                        {item.title || 'Documentation Page'}
                                    </span>
                                    <HighlightMatch text={item.description || ''} query={debouncedQuery} />
                                </div>
                            ),
                            url: targetUrl,
                            type: 'page'
                        });

                        if (formattedResults.length >= 15) break;
                    }
                }

                setResults(formattedResults);

            } catch (err: any) {
                if (err.name !== 'AbortError') {
                    console.error('Search error:', err);
                    setResults([]);
                }
            } finally {
                setIsLoading(false);
            }
        }

        fetchResults();

        return () => abortController.abort();
    }, [debouncedQuery, currentVersion]);

    return (
        <SearchDialog
            {...props}
            search={query}
            onSearchChange={setQuery}
        >
            <SearchDialogOverlay />
            <SearchDialogContent>
                <SearchDialogHeader>
                    <SearchDialogIcon />
                    <SearchDialogInput placeholder="Search documentation..." />

                    {isLoading && (
                        <div className="flex items-center justify-center px-2">
                            <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                        </div>
                    )}

                    <SearchDialogClose />
                </SearchDialogHeader>

                <SearchDialogList items={results} />

                {!isLoading && debouncedQuery && results.length === 0 && (
                    <div className="p-6 text-center text-sm text-muted-foreground">
                        No results found for "<span className="font-semibold text-foreground">{debouncedQuery}</span>".
                    </div>
                )}
            </SearchDialogContent>
        </SearchDialog>
    );
}