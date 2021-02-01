import send from '@polka/send';

import get_guides from './_guides.js';

let json;

export function get(req, res) {
    if (!json || process.env.NODE_ENV !== 'production') {
        const guides = get_guides()
            .map(guide => {
                return {
                    fragment: guide.fragment,
                    answer: guide.answer,
                    metadata: guide.metadata
                };
            });

        json = JSON.stringify(guides);
    }

    send(res, 200, json, {
        'Content-Type': 'application/json',
        'Cache-Control': `max-age=${5 * 60 * 1e3}` // 5 minutes
    });
}
