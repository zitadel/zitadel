export let combinedSlugs = {};
export let dirpath = 'docs/';

export default function generate_search_indexes(dirpath) {
    // return fs
    //     .readdirSync(`${dirpath}`)
    //     .filter((file) => {
    //         return file[0] !== '.' && path.extname(file) === '.md' && file.endsWith(`.${lang}.md`);
    //     })
    //     .map((file) => {
    //         const markdown = fs.readFileSync(`${dirpath}/${file}`, 'utf-8');
    //         const { content, metadata } = extract_frontmatter(markdown);
    //         const section_slug = make_slug(metadata.title);
    //         const subsections = [];

    //         // const slugger = new marked.Slugger();
    //         renderer.heading = (text, level, rawtext) => {
    //             const slug = level <= 4 && make_slug(rawtext);

    //             if (level === 3 || level === 4) {
    //                 const title = text.replace(/<\/?code>/g, '').replace(/\.(\w+)(\((.+)?\))?/, (m, $1, $2, $3) => {
    //                     if ($3) return `.${$1}(...)`;
    //                     if ($2) return `.${$1}()`;
    //                     return `.${$1}`;
    //                 });

    //                 subsections.push({ slug, title, level });
    //             }

    //             return `
    // 				<h${level}>
    // 					<span id="${slug}" class="offset-anchor" ${level > 4 ? 'data-scrollignore' : ''}></span>
    // 					<a href="${dir}#${slug}" class="anchor" aria-hidden="true"> <i class="las la-link"></i> </a>
    // 					${text}
    // 				</h${level}>`;
    //         };

    //         block_types.forEach((type) => {
    //             const fn = renderer[type];
    //             renderer[type] = function () {
    //                 return fn.apply(this, arguments);
    //             };
    //         });

    //         const html = marked(content, { renderer });
    //         const hashes = {};

    //         return {
    //             html: html.replace(/@@(\d+)/g, (m, id) => hashes[id] || m),
    //             metadata,
    //             subsections,
    //             slug: section_slug,
    //             file
    //         };
    //     });
}
