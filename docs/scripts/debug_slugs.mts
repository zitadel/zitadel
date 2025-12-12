import { source } from '../lib/source';

console.log("Printing slugs for user service:");
source.getPages().forEach(page => {
    if (page.url.includes('zitadel.user.v2.UserService.CreateUser')) {
        console.log(`URL: ${page.url}`);
        console.log(`Slug: ${JSON.stringify(page.slugs)}`);
    }
});
