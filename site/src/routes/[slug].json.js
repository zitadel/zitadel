import send from '@polka/send';
import { locale } from 'svelte-i18n';

import { LANGUAGES } from '../../config.js';
import { INIT_OPTIONS } from '../i18n.js';
import generate_docs from '../utils/generate_docs.js';

let json;

export function get(req, res) {
    if (!json || process.env.NODE_ENV !== 'production') {
        const { slug } = req.params;
        locale.subscribe(localecode => {
            console.log('sublocale: ' + localecode, LANGUAGES);
            if (!LANGUAGES.includes(localecode)) {
                console.log(INIT_OPTIONS);
                localecode = INIT_OPTIONS.initialLocale || 'en';
            }
            json = JSON.stringify(generate_docs('docs/', slug, localecode)); // TODO it errors if I send the non-stringified value
        });
    }

    send(res, 200, json, {
        'Content-Type': 'application/json'
    });
}
