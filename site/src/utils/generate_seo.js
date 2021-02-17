import fs from 'fs';

export default function generate_seo(dirpath, dir) {
    try {
        return fs.readFileSync(`${dirpath}${dir}/seo.html`, 'utf-8');
    } catch (error) {
        return '';
    }
};
