import { glob } from 'glob';
import { resolve } from 'node:path';

const CONTENT_ROOT = resolve('content');
async function test() {
    const files = await glob('**/*.mdx', { cwd: CONTENT_ROOT });
    console.log(`Found ${files.length} files`);
    if (files.length > 0) {
        console.log(`Sample: ${files[0]}`);
    }
}
test();
