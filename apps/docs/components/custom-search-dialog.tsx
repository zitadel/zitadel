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
import { CircleAlert, Search } from 'lucide-react';

type SearchResultItem = {
    id: string;
    content: ReactNode;
    url: string;
    type: 'page';
};

function escapeRegExp(string: string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function formatBreadcrumb(url: string, searchType: 'docs' | 'api'): string[] {
    if (!url) return [];

    const base = typeof window !== 'undefined' ? window.location.origin : 'http://localhost';
    const basePath = new URL(url, base).pathname;

    // 1. Get all valid segments
    const segments = basePath.split('/').filter(Boolean);

    // 2. Remove leading 'docs' if present
    if (segments.length > 0 && segments[0].toLowerCase() === 'docs') {
        segments.shift();
    }

    // 3. Remove version segment (e.g., 'v4.17') if it immediately follows 'docs'
    if (segments.length > 0 && /^v\d+\.\d+$/i.test(segments[0])) {
        segments.shift();
    }

    // 4. For API searches, remove leading 'reference' or 'api'
    if (searchType === 'api') {
        while (
            segments.length > 0 &&
            (segments[0].toLowerCase() === 'reference' || segments[0].toLowerCase() === 'api')
        ) {
            segments.shift();
        }
    }

    // Remove the last segment (the page itself), decode, and split by hyphens or underscores
    return segments.slice(0, -1).map(segment =>
        decodeURIComponent(segment)
            .split(/[-_]/)
            .map(word => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ')
    );
}

const HighlightMatch = ({ text, query }: { text: string; query: string }) => {
    if (!query || !text) {
        return <span className="text-sm text-muted-foreground line-clamp-2 leading-relaxed">{text}</span>;
    }

    const safeRegex = new RegExp(`(${escapeRegExp(query)})`, 'gi');
    const parts = text.split(safeRegex);

    return (
        <span className="text-sm text-muted-foreground line-clamp-2 leading-relaxed">
            {parts.map((part, i) =>
                part.toLowerCase() === query.toLowerCase() ? (
                    <mark
                        key={i}
                        className="bg-primary/20 text-primary font-medium rounded-sm px-0.5 py-0.5"
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
    const [searchType, setSearchType] = useState<'docs' | 'api'>('docs');
    const [results, setResults] = useState<SearchResultItem[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);

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
            setErrorMsg(null);
            return;
        }

        setIsLoading(true);
        setErrorMsg(null);
        const abortController = new AbortController();

        async function fetchResults() {
            try {
                const searchUrl = `/docs/api/search?q=${encodeURIComponent(debouncedQuery)}&type=${searchType}&tag=${encodeURIComponent(currentVersion)}`;
                const res = await fetch(searchUrl, { signal: abortController.signal });

                if (!res.ok) {
                    if (res.status === 429) throw new Error('Too many requests. Please slow down.');
                    throw new Error('Search service is currently unavailable.');
                }

                const data = await res.json();
                const seenUrls = new Set<string>();
                const formattedResults: SearchResultItem[] = [];

                for (const item of data) {
                    const targetUrl = item.url;

                    if (!seenUrls.has(targetUrl)) {
                        seenUrls.add(targetUrl);
                        const breadcrumbs = formatBreadcrumb(targetUrl, searchType);

                        formattedResults.push({
                            id: targetUrl,
                            content: (
                                <div className="flex flex-col gap-1.5 py-1">
                                    {breadcrumbs.length > 0 && (
                                        <div className="flex items-center gap-1.5 text-[11px] font-semibold text-muted-foreground/80 uppercase tracking-wider">
                                            {breadcrumbs.map((crumb, idx) => (
                                                <span key={idx} className="flex items-center gap-1.5">
                                                    {idx > 0 && <span className="text-muted-foreground/40">›</span>}
                                                    <span>{crumb}</span>
                                                </span>
                                            ))}
                                        </div>
                                    )}
                                    <span className="font-semibold text-foreground text-base">
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
                    setErrorMsg(err.message);
                    setResults([]);
                }
            } finally {
                // Fix: Only clear loading if this specific request wasn't superseded/aborted
                if (!abortController.signal.aborted) {
                    setIsLoading(false);
                }
            }
        }

        fetchResults();

        return () => abortController.abort();
    }, [debouncedQuery, searchType, currentVersion]);

    return (
        <SearchDialog
            {...props}
            search={query}
            onSearchChange={setQuery}
        >
            <SearchDialogOverlay />
            <SearchDialogContent>
                <SearchDialogHeader className="flex flex-col gap-3 border-b pb-3">
                    <div className="flex items-center gap-2 w-full">
                        <SearchDialogIcon />
                        <SearchDialogInput placeholder={searchType === 'docs' ? 'Search documentation...' : 'Search API endpoints...'} />

                        {isLoading && (
                            <div className="flex items-center justify-center px-2">
                                <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                            </div>
                        )}

                        <SearchDialogClose />
                    </div>

                    <div className="flex justify-start w-full gap-2 px-3">
                        <button
                            type="button"
                            onClick={() => setSearchType('docs')}
                            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition-all duration-200 ${searchType === 'docs'
                                ? 'bg-primary text-primary-foreground shadow-sm'
                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                                }`}
                        >
                            Guides & Docs
                        </button>
                        <button
                            type="button"
                            onClick={() => setSearchType('api')}
                            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition-all duration-200 ${searchType === 'api'
                                ? 'bg-primary text-primary-foreground shadow-sm'
                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                                }`}
                        >
                            API Endpoints
                        </button>
                    </div>
                </SearchDialogHeader>

                {results.length > 0 && <SearchDialogList items={results} />}

                {errorMsg && !isLoading && (
                    <div className="flex flex-col items-center justify-center py-10 px-6 text-center">
                        <span className="text-red-500 bg-red-500/10 p-3 rounded-full mb-3">
                            <CircleAlert className="w-6 h-6" />
                        </span>
                        <p className="text-sm text-foreground font-medium">{errorMsg}</p>
                    </div>
                )}

                {!debouncedQuery && !errorMsg && (
                    <div className="flex flex-col items-center justify-center py-12 px-6 text-center text-muted-foreground">
                        <Search className="w-12 h-12 mb-4 text-muted-foreground/30" />
                        <p className="text-sm font-medium">Type a query to search {searchType === 'docs' ? 'Guides & Docs' : 'API Endpoints'}.</p>
                    </div>
                )}

                {!isLoading && debouncedQuery && results.length === 0 && !errorMsg && (
                    <div className="flex flex-col items-center justify-center py-12 px-6 text-center text-muted-foreground">
                        <p className="text-sm font-medium">
                            No results found in {searchType === 'docs' ? 'Docs' : 'API Endpoints'} for "<span className="font-semibold text-foreground">{debouncedQuery}</span>".
                        </p>
                        <p className="text-xs mt-2">Try adjusting your search terms or checking for typos.</p>
                    </div>
                )}
            </SearchDialogContent>
        </SearchDialog>
    );
}
