import { glob } from 'glob';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT_DIR = path.join(__dirname, '..');
const OLD_XML = path.join(ROOT_DIR, 'old.xml');
const CONTENT_DIR = path.join(ROOT_DIR, 'content');

function extractUrlsFromXml(xmlContent) {
    const regex = /<loc>(https:\/\/zitadel\.com\/docs\/[^<]+)<\/loc>/g;
    const urls = [];
    let match;
    while ((match = regex.exec(xmlContent)) !== null) {
        urls.push(match[1]);
    }
    return urls;
}

async function getNewUrlsFromContent() {
    const mdxFiles = await glob('**/reference/api/**/*.mdx', { cwd: CONTENT_DIR });
    const urls = mdxFiles.map(file => {
        let slug = file.replace(/\.mdx$/, '');
        return `https://zitadel.com/docs/${slug}`;
    });
    return urls;
}

function parseOldUrl(url) {
    const prefix = 'https://zitadel.com/docs/apis/resources/';
    if (!url.startsWith(prefix)) return null;

    const rest = url.slice(prefix.length);
    const [serviceDir, methodSlug] = rest.split('/');
    
    if (!methodSlug) return { serviceDir, isIndex: true, full: url };

    return { serviceDir, methodSlug, full: url };
}

function parseNewUrl(url) {
    let rest;
    if (url.includes('/reference/api/')) {
        rest = url.split('/reference/api/')[1];
    } else {
        return null;
    }

    const parts = rest.split('/');
    if (parts.length < 2) return null;

    const serviceDir = parts[0]; 
    const fileSlug = parts[1]; 

    const methodParts = fileSlug.split('.');
    const method = methodParts[methodParts.length - 1];

    return { serviceDir, method, full: url, fileSlug };
}

const serviceMapping = {
    'action_service_v2': 'action',
    'admin': 'admin',
    'mgmt': 'management',
    'auth': 'auth',
    'system': 'system',
    'user_service_v2': 'user',
    'session_service_v2': 'session',
    'oidc_service_v2': 'oidc',
    'settings_service_v2': 'settings',
    'org_service_v2': 'org',
    'org_service/v2': 'org',
    'project_service_v2': 'project',
    'feature_service_v2': 'feature',
    'idp_service_v2': 'idp',
    'instance_service_v2': 'instance',
    'saml_service_v2': 'saml',
    'internal_permission_service_v2': 'internal_permission',
    'application_service_v2': 'application',
    'authorization_service_v2': 'authorization',
    'webkey_service_v2': 'webkey',
};

async function run() {
    if (!fs.existsSync(OLD_XML)) {
        console.error(`Old sitemap not found at ${OLD_XML}`);
        process.exit(1);
    }
    const oldContent = fs.readFileSync(OLD_XML, 'utf8');
    const oldUrls = extractUrlsFromXml(oldContent).filter(u => u.includes('/apis/resources/'));
    
    const newUrls = (await getNewUrlsFromContent()).filter(u => !/\/v\d+\.\d+(\.\d+)?\//.test(u));
    console.log(`Scanning content directory: found ${newUrls.length} unversioned new URLs.`);

    const newServiceMap = {}; 
    
    for (const u of newUrls) {
        const parsed = parseNewUrl(u);
        if (!parsed) continue;
        
        const isVersioned = /\/v\d+\.\d+\.\d+\//.test(u);
        
        if (!newServiceMap[parsed.serviceDir]) newServiceMap[parsed.serviceDir] = [];
        
        if (isVersioned) {
             const existing = newServiceMap[parsed.serviceDir].find(n => n.method === parsed.method && n.fileSlug === parsed.fileSlug && !/\/v\d+\.\d+\.\d+\//.test(n.full));
             if (existing) continue;
        }

        newServiceMap[parsed.serviceDir].push(parsed);
    }

    const redirects = [];
    const missing = [];

    for (const u of oldUrls) {
        const parsed = parseOldUrl(u);
        if (!parsed) continue;

        if (parsed.isIndex) {
            const rawService = parsed.serviceDir;
            const mappedService = serviceMapping[rawService] || rawService;
             if (newServiceMap[mappedService]) {
                 redirects.push({
                     source: u.replace('https://zitadel.com/docs', ''),
                     destination: `/reference/api/${mappedService}`,
                     permanent: true
                 });
             } else {
                 missing.push(u);
             }
            continue;
        }

        const mappedService = serviceMapping[parsed.serviceDir];
        if (!mappedService || !newServiceMap[mappedService]) {
            missing.push(u);
            continue;
        }

        const candidates = newServiceMap[mappedService];
        
        candidates.sort((a, b) => {
            const aVersioned = /\/v\d+\.\d+\.\d+\//.test(a.full);
            const bVersioned = /\/v\d+\.\d+\.\d+\//.test(b.full);
            if (!aVersioned && bVersioned) return -1;
            if (aVersioned && !bVersioned) return 1;

            const lengthDiff = b.method.length - a.method.length;
            if (lengthDiff !== 0) return lengthDiff;
            
            if (a.fileSlug.includes('.v2.') && !b.fileSlug.includes('.v2.')) return -1;
            if (!a.fileSlug.includes('.v2.') && b.fileSlug.includes('.v2.')) return 1;
            
            return 0;
        });

        let found = null;

        // Manual overrides for UserService "ghost" methods
        if (mappedService === 'user') {
            if (parsed.methodSlug === 'user-service-create-user') {
                found = candidates.find(c => c.method === 'AddHumanUser');
            } else if (parsed.methodSlug === 'user-service-update-user') {
                found = candidates.find(c => c.method === 'UpdateHumanUser');
            }
        }

        if (!found) {
            for (const cand of candidates) {
                const candKebab = cand.method.replace(/([a-z0-9])([A-Z])/g, '$1-$2').toLowerCase();
                
                if (parsed.methodSlug === candKebab || parsed.methodSlug.endsWith('-' + candKebab)) {
                    found = cand;
                    break;
                }
            }
        }

        if (found) {
            redirects.push({
                source: u.replace('https://zitadel.com/docs', ''),
                destination: found.full.replace('https://zitadel.com/docs', ''),
                permanent: true
            });
        } else {
             missing.push(u);
        }
    }

    console.log(`Generated ${redirects.length} redirects.`);
    console.log(`Missing ${missing.length} URLs.`);
    
    fs.writeFileSync(path.join(ROOT_DIR, 'redirects.json'), JSON.stringify(redirects, null, 2));
    fs.writeFileSync(path.join(ROOT_DIR, 'missing.txt'), missing.join('\n'));
}

run();
