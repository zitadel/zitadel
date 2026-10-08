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

function escapeRegExp(string: string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

const HighlightMatch = ({ text, query }: { text: string; query: string }) => {
    if (!query || !text) {
        return <span className="text-xs text-muted-foreground line-clamp-2">{text}</span>;
    }

    const safeRegex = new RegExp(`(${escapeRegExp(query)})`, 'gi');
    const parts = text.split(safeRegex);

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
    const [searchType, setSearchType] = useState<'docs' | 'api'>('docs');
    const [results, setResults] = useState<SearchResultItem[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    // New error state to track backend failures
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
                <SearchDialogHeader className="flex flex-col gap-2 border-b pb-2">
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

                    <div className="flex justify-start w-full gap-1 px-3 pt-1">
                        <button
                            type="button"
                            onClick={() => setSearchType('docs')}
                            className={`px-3 py-1 text-xs font-medium rounded-md transition-colors ${searchType === 'docs'
                                ? 'bg-primary/10 text-primary border border-primary/20'
                                : 'text-muted-foreground hover:bg-muted'
                                }`}
                        >
                            Guides & Docs
                        </button>
                        <button
                            type="button"
                            onClick={() => setSearchType('api')}
                            className={`px-3 py-1 text-xs font-medium rounded-md transition-colors ${searchType === 'api'
                                ? 'bg-primary/10 text-primary border border-primary/20'
                                : 'text-muted-foreground hover:bg-muted'
                                }`}
                        >
                            API Endpoints
                        </button>
                    </div>
                </SearchDialogHeader>

                {results.length > 0 && <SearchDialogList items={results} />}

                {/* Explicit Error State */}
                {errorMsg && !isLoading && (
                    <div className="p-6 text-center text-sm text-red-500 font-medium">
                        {errorMsg}
                    </div>
                )}

                {!debouncedQuery && !errorMsg && (
                    <div className="p-6 text-center text-sm text-muted-foreground">
                        Type a query to search {searchType === 'docs' ? 'Guides & Docs' : 'API Endpoints'}.
                    </div>
                )}

                {!isLoading && debouncedQuery && results.length === 0 && !errorMsg && (
                    <div className="p-6 text-center text-sm text-muted-foreground">
                        No results found in {searchType === 'docs' ? 'Docs' : 'API Endpoints'} for "<span className="font-semibold text-foreground">{debouncedQuery}</span>".
                    </div>
                )}
            </SearchDialogContent>
        </SearchDialog>
    );
}
