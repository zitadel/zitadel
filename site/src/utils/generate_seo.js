import fs from 'fs';

export default function generate_docs(dirpath, dir, lang) {
    try {
        return fs.readFileSync(`${dirpath}${dir}/seo_${lang}.html`, 'utf-8');
    } catch (error) {
        return '';
    }
};
