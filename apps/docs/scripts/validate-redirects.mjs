import fs from 'fs';
import path from 'path';

const REDIRECTS_PATH = '/home/ffo/git/zitadel/zitadel/apps/docs/redirects.json';

function validateRedirects() {
    const content = fs.readFileSync(REDIRECTS_PATH, 'utf8');
    const redirects = JSON.parse(content);

    const sources = new Map();
    const inconsistencies = [];

    // 1. Basic Format & Duplicates & Self-Redirects
    redirects.forEach((r, index) => {
        if (!r.source || !r.destination) {
            inconsistencies.push({ type: 'MISSING_FIELDS', index, entry: r });
            return;
        }

        if (!r.source.startsWith('/')) {
            inconsistencies.push({ type: 'SOURCE_MISSING_SLASH', index, source: r.source });
        }

        if (r.destination.startsWith('/docs/') && r.basePath !== false) {
             inconsistencies.push({ type: 'DESTINATION_HAS_DOCS_PREFIX', index, source: r.source, destination: r.destination });
        }

        if (r.source === r.destination) {
            inconsistencies.push({ type: 'SELF_REDIRECT', index, source: r.source });
        }

        if (sources.has(r.source)) {
            inconsistencies.push({ type: 'DUPLICATE_SOURCE', index, source: r.source, previousIndex: sources.get(r.source).index });
        } else {
            sources.set(r.source, { destination: r.destination, index });
        }
    });

    // 2. Circular & Chains
    const checked = new Set();
    sources.forEach((value, source) => {
        let currentSource = source;
        const path = [currentSource];
        const visitedInCurrentPath = new Set([currentSource]);

        while (sources.has(sources.get(currentSource)?.destination)) {
            const nextDestination = sources.get(currentSource).destination;
            
            // Normalize for comparison if one has /docs and other doesn't
            // But for now let's just use exact match since they are all consistent
            
            if (visitedInCurrentPath.has(nextDestination)) {
                inconsistencies.push({ type: 'CIRCULAR_REDIRECT', path: [...path, nextDestination] });
                break;
            }

            path.push(nextDestination);
            visitedInCurrentPath.add(nextDestination);
            currentSource = nextDestination;
        }
        
        if (path.length > 2) {
             inconsistencies.push({ type: 'REDIRECT_CHAIN', path });
        }
    });

    // Deduplicate inconsistencies
    const uniqueInconsistencies = [];
    const seen = new Set();
    const counts = {};

    inconsistencies.forEach(inc => {
        counts[inc.type] = (counts[inc.type] || 0) + 1;
        const key = JSON.stringify(inc);
        if (!seen.has(key)) {
            uniqueInconsistencies.push(inc);
            seen.add(key);
        }
    });

    console.log('Inconsistency Counts:', JSON.stringify(counts, null, 2));
    
    const others = uniqueInconsistencies.filter(inc => inc.type !== 'DESTINATION_HAS_DOCS_PREFIX');
    if (others.length > 0) {
        console.log('Other Inconsistencies:', JSON.stringify(others, null, 2));
    } else {
        console.log('No other inconsistencies found.');
    }
}

validateRedirects();
