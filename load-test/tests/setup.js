import url from './url.js';
import http from 'k6/http';
import { check } from 'k6';
import { Trend } from 'k6/metrics';
import { createHuman } from './user.js';

import { Config } from './config.js';

export default async function Setup(accessToken) {
    const org = await createOrg(accessToken);
    const user = await createHuman('gigi', org, accessToken);
    return {org, user};
}

const createOrgTrend = new Trend('setup_create_org_duration', true);
export function createOrg(accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST', 
            url('/v2beta/organizations'), 
            JSON.stringify({
                name: `load-test-${new Date(Date.now()).toISOString()}`
            }), 
            {
                headers: {
                    authorization: `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                    'x-zitadel-orgid': Config.orgId
                }
            }
        )

        response.then((res) => {
            check(res, {
                'org created': (r) => r.status === 201  || reject(`unable to create org status: ${res.status} || body: ${res.body}`)
            });
        
            createOrgTrend.add(res.timings.duration);
        
            resolve(res.json());
        });
    })
}

const addProjectTrend = new Trend('setup_add_project_duration', true);
export function createProject(name, org, accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST',
            url('/management/v1/projects'),
            JSON.stringify({
                name: name
            }),
            {
                headers: {
                    authorization: `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                    'x-zitadel-orgid': org.organizationId
                }
            }
        );
        response.then((res) => {
            check(res, {
                "add project status ok": (r) => r.status === 200
            }) || reject(`unable to add project status: ${res.status} body: ${res.body}`);
            resolve(res.json());

            addProjectTrend.add(res.timings.duration);
            resolve(res.json());
        });
    });
}

const addAPITrend = new Trend('setup_add_app_duration', true);
export function createAPI(name, projectId, org, accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST',
            url(`/management/v1/projects/${projectId}/apps/api`),
            JSON.stringify({
                name: name,
                authMethodType: "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
            }),
            {
                headers: {
                    authorization: `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                    'x-zitadel-orgid': org.organizationId
                }
            }
        );
        response.then((res) => {
            check(res, {
                "add api status ok": (r) => r.status === 200
            }) || reject(`unable to add api project: ${projectId} status: ${res.status} body: ${res.body}`);
            resolve(res.json());
            
            addAPITrend.add(res.timings.duration);
            resolve(res.json());
        });
    });
}

const addAppKeyTrend = new Trend('setup_add_app_key_duration', true);
export function createAppKey(appId, projectId, org, accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST',
            url(`/management/v1/projects/${projectId}/apps/${appId}/keys`),
            JSON.stringify({
                type: "KEY_TYPE_JSON"
            }),
            {
                headers: {
                    authorization: `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                    'x-zitadel-orgid': org.organizationId
                }
            }
        );
        response.then((res) => {
            check(res, {
                "add app key status ok": (r) => r.status === 200
            }) || reject(`unable to add app key project: ${projectId} app: ${appId} status: ${res.status} body: ${res.body}`);
            resolve(res.json());
            
            
            addAPITrend.add(res.timings.duration);
            resolve(res.json());
        });
    });
}