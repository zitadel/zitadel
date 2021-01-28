import send from '@polka/send';

import generate_docs from '../utils/generate_docs.js';
import generate_seo from '../utils/generate_seo.js';

let json;

export function get(req, res) {
    if (!json || process.env.NODE_ENV !== 'production') {
        const { slug } = req.params;
        const localecode = 'en';
        // TODO replace with localecode from url param later
        const seo = generate_seo(`docs/${localecode}/`, slug);
        const docs = generate_docs(`docs/${localecode}/`, slug);
        json = JSON.stringify({ docs, seo }); // TODO it errors if I send the non-stringified value
    }

    send(res, 200, json, {
        'Content-Type': 'application/json'
    });
}
