#!/usr/bin/env node

import { program } from 'commander';
import fs from 'fs';
import * as envfile from 'envfile';
import http from 'node:http';
import open from 'open';
import url from 'url';
import { listFrameworks } from '@netlify/framework-info';

class ApplicationType {
    static WEB = "web";
    static NATIVE = "native";
    static USER_AGENT = "user_agent";
    static API = "api";
    static SAML = "saml";
}
class Scope {
    static OPENID = "openid";
    static PROFILE = "profile";
    static EMAIL = "email";
    static OFFLINE_ACCESS = "offline_access";
}

class GrantType {
    static AUTHORIZATION_CODE = "authorization_code";
    static IMPLICIT = "implicit";
    static REFRESH_TOKEN = "refresh_token";
    static DEVICE_CODE = "device_code";
}

class ResponseType {
    static CODE = "code";
    static ID_TOKEN = "id_token";
    static ID_TOKEN_TOKEN = "id_token token"
}

class AuthenticationType {
    static BASIC = "basic";
    static NONE = "none";
    static POST = "post";
    static PRIVATE_KEY_JWT = "private_key_jwt";
}

class AppTypes {
    static ANGULAR = {
        type: "angular",
        applicationType: ApplicationType.USER_AGENT,
        redirectUris: ["localhost:4200"],
        postLogoutUris: ["localhost:4200"],
        scope: [Scope.OPENID, Scope.EMAIL, Scope.PROFILE, Scope.OFFLINE_ACCESS],
        grantTypes: [GrantType.AUTHORIZATION_CODE],
        responseTypes: [ResponseType.CODE],
        authenticationType: AuthenticationType.NONE
    }
    static NEXT = {
        type: "next",
        applicationType: ApplicationType.WEB,
        redirectUris: ["http://localhost:3000/api/auth/callback/zitadel"],
        postLogoutUris: ["http://localhost:3000"],
        scope: [Scope.OPENID, Scope.EMAIL, Scope.PROFILE],
        grantTypes: [GrantType.AUTHORIZATION_CODE],
        responseTypes: [ResponseType.CODE],
        authenticationType: AuthenticationType.BASIC
    }
    static NUXT = {
        type: "nuxt",
        applicationType: "",
        redirectUris: [""],
        postLogoutUris: [""],
        scope: [Scope.OPENID, Scope.EMAIL, Scope.PROFILE],
        grantTypes: [""],
        responseTypes: [""],
        authenticationType: ""
    }
    static REACT = {
        type: "react",
        applicationType: ApplicationType.USER_AGENT,
        redirectUris: ["http://localhost:3000/callback"],
        postLogoutUris: ["http://localhost:3000"],
        scope: [Scope.OPENID, Scope.EMAIL, Scope.PROFILE],
        grantTypes: [GrantType.AUTHORIZATION_CODE],
        responseTypes: [ResponseType.CODE],
        authenticationType: AuthenticationType.NONE
    }
    static SVELTE = {
        type: "svelte",
        applicationType: ApplicationType.WEB,
        redirectUris: ["http://localhost:3000/api/auth/callback/zitadel"],
        postLogoutUris: ["http://localhost:3000"],
        scope: [Scope.OPENID, Scope.EMAIL, Scope.PROFILE],
        grantTypes: [GrantType.AUTHORIZATION_CODE],
        responseTypes: [ResponseType.CODE],
        authenticationType: AuthenticationType.BASIC
    }
    static VUE = {
        type: "vue",
        applicationType: ApplicationType.USER_AGENT,
        redirectUris: ['http://localhost:5173/auth/signinwin/zitadel'],
        postLogoutUris: ['http://localhost:5173'],
        scope: [Scope.OPENID, Scope.EMAIL, Scope.PROFILE],
        grantTypes: [GrantType.AUTHORIZATION_CODE],
        responseTypes: [ResponseType.CODE],
        authenticationType: AuthenticationType.NONE
    }
    static VITE = {
        type: "vite",
        applicationType: ApplicationType.USER_AGENT,
        redirectUris: ["http://localhost:3000/callback"],
        postLogoutUris: ["http://localhost:3000"],
        scope: [Scope.OPENID, Scope.EMAIL, Scope.PROFILE],
        grantTypes: [GrantType.AUTHORIZATION_CODE],
        responseTypes: [ResponseType.CODE],
        authenticationType: AuthenticationType.NONE
    }
    static UNKNOWN = {
        type: "unknown",
        applicationType: ApplicationType.USER_AGENT,
        redirectUris: ["http://localhost:3000"],
        postLogoutUris: ["http://localhost:3000"],
        scope: [Scope.OPENID, Scope.EMAIL, Scope.PROFILE, Scope.OFFLINE_ACCESS],
        grantTypes: [GrantType.AUTHORIZATION_CODE],
        responseTypes: [ResponseType.CODE],
        authenticationType: AuthenticationType.NONE
    }
}

program
    .version('1.0.0', '-v, --version')
    .usage('[OPTIONS]...')
    .option('--env-file <PATH>', 'path to .env file, if not set it is created in the current directory', './.env')
    .option('-p, --project-id <PROJECT_ID>', 'id of the project you want the app to be created in')
    .option('-n, --app-name <APP_NAME>', 'name of the app you want to create')
    .option('-o, --org-id <ORG_ID>', 'id of the org you want the app ot be created in, by default it uses the organization of the signed in user')
    .requiredOption('-u, --url <URL>', 'url to your zitadel instance (e.g. https://myprefix.zitadel.cloud)')
    .parse(process.argv);

const options = program.opts();
const port = 3333;

let zitadelUrl = await generateURL(options, port);
open(zitadelUrl.toString());

awaitCredentials(port);

function awaitCredentials(port) {
    new Promise((resolve) => {
        const server = http.createServer();
        
        server.on('request', (req, res) => {
            res.writeHead(200);
            resolve({data: url.parse(req.url, true).query, server});

            res.end('.env is configured, continue in your app');
        })

        server.listen(port, () => {
            console.log(`Awaiting response from browser`);
        });
    }).then((res) => {
        writeCredentials(res.data.clientId, res.data.clientSecret);
        res.server.close();
    })
}

function writeCredentials(clientID, clientSecret) {
    let parsedEnv = envfile.parse(options.envFile);
    parsedEnv.CLIENT_ID = clientID ? clientID : parsedEnv.CLIENT_ID ? parsedEnv.CLIENT_ID : '';
    parsedEnv.CLIENT_SECRET = clientSecret ? clientSecret : parsedEnv.CLIENT_SECRET ? parsedEnv.CLIENT_SECRET : '';
    fs.writeFileSync(options.envFile, envfile.stringify(parsedEnv));

    console.log(`credentials were written to ${options.envFile}`)
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
                config = AppTypes.ANGULAR;
                break;
            case 'next-nx':
            case 'next':
                config = AppTypes.NEXT;
                break;
            case 'nuxt':
            case 'nuxt3':
                config = AppTypes.NUXT;
                break;
            case 'react-static':
            case 'create-react-app':
                config = AppTypes.REACT;
                break;
            case 'svelte':
            case 'svelte-kit':
                config = AppTypes.SVELTE;
                break;
            case 'vue':
                config = AppTypes.VUE;
                break;
            case 'vite':
                config = AppTypes.VITE;
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
                config = AppTypes.UNKNOWN;
            break;
        }
        config.frameworks = frameworks.map(framework => framework.id);
        return config;
    });
}

async function generateURL(options, port) {
    let config = await getConfigFromFramework();

    if (options?.projectId != undefined) {
        config.projectId = options.projectId
    }
    if (options?.orgId != undefined) {
        config.orgId = options.orgId
    }
    if (options?.appName != undefined) {
        config.appName = options.appName
    }

    config.redirectTo = new URL(`http://localhost${port}`).toString();

    let url = new URL(`${options.url}/ui/console/projects/app-create`)
    url.search = new URLSearchParams(config);

    return url;
}