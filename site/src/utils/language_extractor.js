import fs from 'fs';
import path from 'path';

export default function extract_languages(dirpath, dir) {

    const detectedLocales = fs.readdirSync(`${dirpath}${dir}`)
        .filter(file => path.extname(file) == '.md')
        .map((file) => {
            file = file.replace(path.extname(file), '');
            const arr = file.split('.');
            const locale = arr.length ? arr[arr.length - 1] : null;
            if (locale) {
                return locale;
            }
        }).filter(locale => locale !== null);

    const redDetectedLocales = [...new Set(detectedLocales)];

    console.log('detected locales: ' + redDetectedLocales);
    return redDetectedLocales;
}
