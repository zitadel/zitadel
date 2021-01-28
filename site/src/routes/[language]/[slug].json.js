import send from '@polka/send';

import generate_docs from '../../utils/generate_docs.js';
import generate_seo from '../../utils/generate_seo.js';

let json;

export function get(req, res) {
    if (!json || process.env.NODE_ENV !== 'production') {
        const { language, slug } = req.params;
        console.log('lang', language);
        const localecode = language ? language : 'en';
        const seo = generate_seo(`docs/${localecode}/`, slug);
        const docs = generate_docs(`docs/${localecode}/`, slug, localecode);
        json = JSON.stringify({ docs, seo }); // TODO it errors if I send the non-stringified value
    }

    send(res, 200, json, {
        'Content-Type': 'application/json'
    });
}
