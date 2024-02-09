#!/usr/bin/env node

import { program } from 'commander';
import fs from 'fs';
import * as envfile from 'envfile';
import http from 'node:http';
import open from 'open';
import url from 'url';
import { listFrameworks } from '@netlify/framework-info'

program
    .version('1.0.0', '-v, --version')
    .usage('[OPTIONS]...')
    .requiredOption('-u, --url <URL>', 'url to your zitadel instance (e.g. https://myprefix.zitadel.cloud)')
    .option('--env-file <PATH>', 'path to .env file, if not set it is created in the current directory', './.env')
    .parse(process.argv);

const options = program.opts();

let config = await getConfigFromFramework();
let params = new URLSearchParams(config).toString()

console.log(params);
console.log(JSON.stringify(config));

open(`${options.url}/ui/console/projects/app-create?${params}`);

awaitCredentials(3000);

function awaitCredentials(port) {
    new Promise((resolve) => {
        const server = http.createServer();
        
        server.on('request', (req, res) => {
            res.writeHead(200);
            resolve({data: url.parse(req.url, true).query, server});

            res.end('.env is configured, continue in your app');
        })

        server.listen(port, () => {
            console.log(`Server listening on port ${port}`);
        });
    }).then((res) => {
        console.log(`got request ${JSON.stringify(res.data)}`);
        writeCredentials(res.data.clientId, res.data.clientSecret);
        res.server.close();
    })
}

function writeCredentials(clientID, clientSecret) {
    let parsedEnv = envfile.parse(options.envFile);
    parsedEnv.CLIENT_ID = clientID ? clientID : parsedEnv.CLIENT_ID ? parsedEnv.CLIENT_ID : '';
    parsedEnv.CLIENT_SECRET = clientSecret ? clientSecret : parsedEnv.CLIENT_SECRET ? parsedEnv.CLIENT_SECRET : '';
    fs.writeFileSync(options.envFile, envfile.stringify(parsedEnv))
}

function getConfigFromFramework() {
    return listFrameworks().then((frameworks)=> {
        if (frameworks.length == 0) {
            console.log('no framework detected');
        }
        let config;
        console.log(frameworks[0].id);
        switch (frameworks[0].id) {
            case 'angular':
                config = readConfig('angular');
                break;
            case 'next-nx':
            case 'next':
                config = readConfig('next');
                break;
            case 'nuxt':
            case 'nuxt3':
                config = readConfig('nuxt');
                break;
            case 'react-static':
            case 'create-react-app':
                config = readConfig('react');
                break;
            case 'svelte':
            case 'svelte-kit':
                config = readConfig('svelte');
                break;
            case 'vue':
                config = readConfig('vue');
                break;
            case 'vite':
                config = readConfig('vite');
                break;
            case 'astro':
            case 'docusaurus':
            case 'docusaurus-v2':
            case 'eleventy':
            case 'gatsby':
            case 'gridsome':
            case 'hexo':
            case 'hugo':
            case 'hydrogen':
            case 'jekyll':
            case 'middleman':
            case 'blitz':
            case 'phenomic':
            case 'qwik':
            case 'redwoodjs':
            case 'remix':
            case 'solid-js':
            case 'solid-start':
            case 'stencil':
            case 'vuepress':
            case 'assemble':
            case 'docpad':
            case 'harp':
            case 'metalsmith':
            case 'roots':
            case 'wintersmith':
            case 'cecil':
            case 'zola':
            case 'ember':
            case 'expo':
            case 'quasar':
            case 'quasar-v0.17':
            case 'sapper':
            case 'brunch':
            case 'parcel':
            case 'grunt':
            case 'gulp':
            case 'wmr':
            default: 
                config = readConfig('angular');
            break;
        }
        return config;
    });
}

function readConfig(name) {
    return {
        "type": "angular",
        "redirectUris": ["localhost:4200"],
        "scope": ["openid", "profile", "email", "offline_access"],
        "grantTypes": ["authorizationCode"],
        "responseTypes": ["code"],
        "authenticationType": "none"
    };
    // return JSON.parse(fs.readFileSync(`./frameworks/${name}.json`));
}