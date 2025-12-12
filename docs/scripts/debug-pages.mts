
import { register } from 'node:module';
import { pathToFileURL } from 'node:url';

register('./scripts/loader.mjs', pathToFileURL('./'));
register('fumadocs-mdx/node/loader', import.meta.url);

const { source } = await import('../lib/source');

console.log('Pages found:', source.getPages().length);
const managementPage = source.getPage(['references', 'api-v1', 'management']);
console.log('Management page found:', !!managementPage);

if (managementPage) {
    console.log('Management page slug:', managementPage.slugs);
    console.log('Management page url:', managementPage.url);
}

const allSlugs = source.getPages().map(p => p.slugs.join('/'));
console.log('Sample slugs:', allSlugs.slice(0, 10));
console.log('Contains references/api-v1/management:', allSlugs.includes('references/api-v1/management'));
