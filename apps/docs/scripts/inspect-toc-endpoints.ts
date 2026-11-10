import { register } from 'node:module';
import { pathToFileURL } from 'node:url';

register('./scripts/loader.mjs', pathToFileURL('./'));
register('fumadocs-mdx/node/loader', import.meta.url);

const { source } = await import('../lib/source');
const pages = source.getPages();

const endpointsPage = pages.find(p => p.url.endsWith('/apis/openidoauth/endpoints'));

if (endpointsPage) {
    console.log('Endpoints Page URL:', endpointsPage.url);
    console.log('TOC URLs:', endpointsPage.data.toc.map(i => i.url));
} else {
    console.log('Endpoints page not found.');
}
