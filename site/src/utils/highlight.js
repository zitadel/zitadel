import 'prismjs/components/prism-bash';

import PrismJS from 'prismjs';

import { langs } from './markdown.js';

export function highlight(source, lang) {
    const plang = langs[lang] || '';
    const highlighted = plang ? PrismJS.highlight(
        source,
        PrismJS.languages[plang],
        lang
    ) : source.replace(/[&<>]/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;' })[c]);

    return `<pre class='language-${plang}'><code>${highlighted}</code></pre>`;
}