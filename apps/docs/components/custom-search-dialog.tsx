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

function formatBreadcrumb(url: string): string[] {
    if (!url) return [];

    let pathname = url;
    try {
        pathname = new URL(url).pathname;
    } catch {
        // If the URL is relative, we can still process it
    }

    const basePath = pathname.split('#')[0];
    const segments = basePath
        .split('/')
        .filter(segment => segment && segment.toLowerCase() !== 'docs');

    return segments.map(segment =>
        segment.split('-')
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
                        const breadcrumbs = formatBreadcrumb(targetUrl);

                        formattedResults.push({
                            id: targetUrl,
                            content: (
                                <div className="flex flex-col gap-1.5 py-1">
                                    {breadcrumbs.length > 0 && (
                                        <div className="flex items-center gap-1.5 text-[11px] font-semibold text-muted-foreground/80 uppercase tracking-wider">
                                            {breadcrumbs.map((crumb, idx) => (
                                                <span key={idx} className="flex items-center gap-1.5">
                                                    {idx > 0 && <span className="text-muted-foreground/40">›</span>}
                                                    <span className="truncate max-w-[120px]">{crumb}</span>
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
                            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" /></svg>
                        </span>
                        <p className="text-sm text-foreground font-medium">{errorMsg}</p>
                    </div>
                )}

                {!debouncedQuery && !errorMsg && (
                    <div className="flex flex-col items-center justify-center py-12 px-6 text-center text-muted-foreground">
                        <svg className="mb-4 text-muted-foreground/30" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8" /><line x1="21" y1="21" x2="16.65" y2="16.65" /></svg>
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