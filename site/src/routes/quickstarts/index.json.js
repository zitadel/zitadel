import send from '@polka/send';

import get_quickstarts from './_quickstarts.js';

let json;

export function get(req, res) {
    if (!json || process.env.NODE_ENV !== 'production') {
        const qss = get_quickstarts()
            .map(qs => {
                return {
                    fragment: qs.fragment,
                    answer: qs.answer,
                    metadata: qs.metadata
                };
            });

        json = JSON.stringify(qss);
    }

    send(res, 200, json, {
        'Content-Type': 'application/json',
        'Cache-Control': `max-age=${5 * 60 * 1e3}` // 5 minutes
    });
}