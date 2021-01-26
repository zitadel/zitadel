import send from '@polka/send';
import { locale } from 'svelte-i18n';

import { LANGUAGES } from '../../config.js';
import { INIT_OPTIONS } from '../i18n.js';
import generate_docs from '../utils/generate_docs.js';
import generate_seo from '../utils/generate_seo.js';

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

            const seo = generate_seo('docs/', slug, localecode);

            // import(`../../../docs/${slug}/seo_${localecode}.jsonld`).then((module) => {
            //     const seo = module.meta;
            //     const docs = generate_docs('docs/', slug, localecode);
            //     json = JSON.stringify({ docs, seo }); // TODO it errors if I send the non-stringified value
            // }).catch(error => {
            //     console.error(error);

            const docs = generate_docs('docs/', slug, localecode);
            json = JSON.stringify({ docs, seo }); // TODO it errors if I send the non-stringified value

            // });
        });
    }

    send(res, 200, json, {
        'Content-Type': 'application/json'
    });
}
